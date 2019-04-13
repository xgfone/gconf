// Copyright 2019 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gconf

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
)

// Validator is an interface to validate whether the value v is valid.
//
// When implementing an Opt, you can supply the method Validate to implement
// the interface Validator, too. The config engine will check and call it.
// So the Opt is the same to implement the interface:
//
//    type ValidatorOpt interface {
//        Opt
//        Validator
//    }
//
// In order to be flexible and customized, the builtin validators use the
// validator chain ValidatorChainOpt to handle more than one validator.
// Notice: they both are the valid Opts with the validator function.
type Validator interface {
	// Validate whether the value v of name in the group is valid.
	//
	// Return nil if the value is ok, or an error instead.
	Validate(groupFullName, optName string, v interface{}) error
}

// ValidatorChainOpt is an Opt interface with more than one validator.
//
// The validators in the chain will be called in turn. The validation is
// considered as failure only if one validator returns an error, that's,
// only all the validators return nil, it's successful.
type ValidatorChainOpt interface {
	Opt

	// ReSet the validator chain.
	//
	// Notice: this method should return the option itself.
	SetValidators(...Validator) ValidatorChainOpt

	// Add some new validators into the validator chain.
	AddValidators(...Validator) ValidatorChainOpt

	// Return the validator chain.
	GetValidators() []Validator
}

var (
	errNil       = fmt.Errorf("the value is nil")
	errStrEmtpy  = fmt.Errorf("the string is empty")
	errStrType   = fmt.Errorf("the value is not string type")
	errIntType   = fmt.Errorf("the value is not an integer type")
	errFloatType = fmt.Errorf("the value is not an float type")
)

// ValidatorError stands for a validator error.
type ValidatorError struct {
	Group string
	Name  string
	Value interface{}
	Err   error
}

// NewValidatorError returns a new ValidatorError.
func NewValidatorError(group, name string, value interface{}, err error) ValidatorError {
	return ValidatorError{Group: group, Name: name, Value: value, Err: err}
}

// NewValidatorErrorf returns a new ValidatorError.
func NewValidatorErrorf(group, name string, value interface{},
	format string, args ...interface{}) ValidatorError {
	return ValidatorError{
		Group: group,
		Name:  name,
		Value: value,
		Err:   fmt.Errorf(format, args...),
	}
}

// Error implements the interface Error.
func (v ValidatorError) Error() string {
	if v.Group == "" {
		return fmt.Sprintf("%s: %v", v.Name, v.Err)
	}
	return fmt.Sprintf("[%s:%s]: %v", v.Group, v.Name, v.Err)
}

func toString(v interface{}) (string, error) {
	if v == nil {
		return "", errNil
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", errStrType
}

func toInt64(v interface{}) (int64, error) {
	if v == nil {
		return 0, errNil
	}

	switch _v := v.(type) {
	case int:
		return int64(_v), nil
	case int8:
		return int64(_v), nil
	case int16:
		return int64(_v), nil
	case int32:
		return int64(_v), nil
	case int64:
		return _v, nil
	case uint:
		return int64(_v), nil
	case uint8:
		return int64(_v), nil
	case uint16:
		return int64(_v), nil
	case uint32:
		return int64(_v), nil
	case uint64:
		return int64(_v), nil
	default:
		return 0, errIntType
	}
}

func toFloat64(v interface{}) (float64, error) {
	if v == nil {
		return 0, errNil
	}

	switch _v := v.(type) {
	case float32:
		return float64(_v), nil
	case float64:
		return _v, nil
	default:
		return 0, errFloatType
	}
}

// ValidatorFunc is a wrapper of a function validator.
type ValidatorFunc func(group, name string, v interface{}) error

// Validate implements the method Validate of the interface Validator.
func (f ValidatorFunc) Validate(group, name string, v interface{}) error {
	return f(group, name, v)
}

// NewStrLenValidator returns a validator to validate that the length of the
// string must be between min and max.
func NewStrLenValidator(min, max int) Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}

		_len := len(s)
		if _len > max || _len < min {
			return NewValidatorErrorf(group, name,
				"the length of '%s' is %d, not between %d and %d",
				s, _len, min, max)
		}
		return nil
	})
}

// NewStrNotEmptyValidator returns a validator to validate that the value must
// not be an empty string.
func NewStrNotEmptyValidator() Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}

		if len(s) == 0 {
			return NewValidatorError(group, name, v, errStrEmtpy)
		}
		return nil
	})
}

// NewStrArrayValidator returns a validator to validate that the value is in
// the array.
func NewStrArrayValidator(array []string) Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		for _, v := range array {
			if s == v {
				return nil
			}
		}
		return NewValidatorErrorf(group, name, "the value '%s' is not in %v", s, array)
	})
}

// NewRegexpValidator returns a validator to validate whether the value match
// the regular expression.
//
// This validator uses regexp.MatchString(pattern, s) to validate it.
func NewRegexpValidator(pattern string) Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}

		if ok, err := regexp.MatchString(pattern, s); err != nil {
			return NewValidatorError(group, name, v, err)
		} else if !ok {
			return NewValidatorErrorf(group, name,
				"'%s' doesn't match the value '%s'", s, pattern)
		}
		return nil
	})
}

// NewURLValidator returns a validator to validate whether a url is valid.
func NewURLValidator() Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if _, err = url.Parse(s); err != nil {
			return NewValidatorError(group, name, v, err)
		}
		return nil
	})
}

// NewIPValidator returns a validator to validate whether an ip is valid.
func NewIPValidator() Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if net.ParseIP(s) == nil {
			return NewValidatorErrorf(group, name, v, "the value '%s' is not a valid ip", s)
		}
		return nil
	})
}

// NewIntegerRangeValidator returns a validator to validate whether the integer
// value is between the min and the max.
//
// This validator can be used to validate the value of the type int, int8,
// int16, int32, int64, uint, uint8, uint16, uint32, uint64.
func NewIntegerRangeValidator(min, max int64) Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		i, err := toInt64(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if min > i || i > max {
			return NewValidatorErrorf(group, name, v,
				"the value '%d' is not between %d and %d", i, min, max)
		}
		return nil
	})
}

// NewFloatRangeValidator returns a validator to validate whether the float
// value is between the min and the max.
//
// This validator can be used to validate the value of the type float32 and
// float64.
func NewFloatRangeValidator(min, max float64) Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		f, err := toFloat64(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if min > f || f > max {
			return NewValidatorErrorf(group, name, v,
				"the value '%f' is not between %f and %f", f, min, max)
		}
		return nil
	})
}

// NewPortValidator returns a validator to validate whether a port is between
// 0 and 65535.
func NewPortValidator() Validator {
	return NewIntegerRangeValidator(0, 65535)
}

// NewEmailValidator returns a validator to validate whether an email is valid.
func NewEmailValidator() Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if _, err = mail.ParseAddress(s); err != nil {
			return NewValidatorError(group, name, v, err)
		}
		return nil
	})
}

// NewAddressValidator returns a validator to validate whether an address is
// like "host:port", "host%zone:port", "[host]:port" or "[host%zone]:port".
//
// This validator uses net.SplitHostPort() to validate it.
func NewAddressValidator() Validator {
	return ValidatorFunc(func(group, name string, v interface{}) error {
		s, err := toString(v)
		if err != nil {
			return NewValidatorError(group, name, v, err)
		}
		if _, _, err = net.SplitHostPort(s); err != nil {
			return NewValidatorError(group, name, v, err)
		}
		return nil
	})
}
