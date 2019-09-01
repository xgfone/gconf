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

var (
	errStrEmtpy    = fmt.Errorf("the string is empty")
	errNotString   = fmt.Errorf("the value is not string")
	errStrNotEmtpy = fmt.Errorf("the string is not empty")

	errNotStringSlice = fmt.Errorf("the value is not []string")
)

// Validator is used to validate whether the value of the option in the group
// is valid.
type Validator func(value interface{}) error

// Or returns a union validator, which returns nil only if a certain validator
// returns nil or the error that the last validator returns.
func Or(validators ...Validator) Validator {
	return func(value interface{}) (err error) {
		for _, v := range validators {
			if err = v(value); err == nil {
				return nil
			}
		}
		return
	}
}

// NewStrLenValidator returns a validator to validate that the length of the
// string must be between min and max.
func NewStrLenValidator(min, max int) Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		_len := len(s)
		if _len > max || _len < min {
			return fmt.Errorf("the length of '%s' is %d, not between %d and %d",
				s, _len, min, max)
		}
		return nil
	}
}

// NewEmptyStrValidator returns a validator to validate that the value must be
// an empty string.
func NewEmptyStrValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if len(s) == 0 {
			return nil
		}
		return errStrNotEmtpy
	}
}

// NewStrNotEmptyValidator returns a validator to validate that the value must
// not be an empty string.
func NewStrNotEmptyValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if len(s) == 0 {
			return errStrEmtpy
		}
		return nil
	}
}

// NewStrArrayValidator returns a validator to validate that the value is in
// the array.
func NewStrArrayValidator(array []string) Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		for _, v := range array {
			if s == v {
				return nil
			}
		}
		return fmt.Errorf("the value '%s' is not in %v", s, array)
	}
}

// NewStrSliceValidator returns a validator to validate whether the string element
// of the []string value satisfies all the given validators.
func NewStrSliceValidator(strValidators ...Validator) Validator {
	return func(value interface{}) (err error) {
		ss, ok := value.([]string)
		if !ok {
			return errNotStringSlice
		}

		for _, s := range ss {
			for _, validator := range strValidators {
				if err = validator(s); err != nil {
					return
				}
			}
		}

		return nil
	}
}

// NewRegexpValidator returns a validator to validate whether the value match
// the regular expression.
//
// This validator uses regexp.MatchString(pattern, s) to validate it.
func NewRegexpValidator(pattern string) Validator {
	re := regexp.MustCompile(pattern)
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if re.MatchString(s) {
			return nil
		}
		return fmt.Errorf("'%s' doesn't match the value '%s'", s, pattern)
	}
}

// NewURLValidator returns a validator to validate whether a url is valid.
func NewURLValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if _, err := url.Parse(s); err != nil {
			return err
		}
		return nil
	}
}

// NewURLSliceValidator returns a validator to validate whether the string element
// of the []string value is a valid URL.
func NewURLSliceValidator() Validator {
	return NewStrSliceValidator(NewURLValidator())
}

// NewIPValidator returns a validator to validate whether an ip is valid.
func NewIPValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if net.ParseIP(s) == nil {
			return fmt.Errorf("the value '%s' is not a valid ip", s)
		}
		return nil
	}
}

// NewIPSliceValidator returns a validator to validate whether the string element
// of the []string value is a valid IP.
func NewIPSliceValidator() Validator {
	return NewStrSliceValidator(NewIPValidator())
}

// NewEmailValidator returns a validator to validate whether an email is valid.
func NewEmailValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if _, err := mail.ParseAddress(s); err != nil {
			return err
		}
		return nil
	}
}

// NewEmailSliceValidator returns a validator to validate whether the string element
// of the []string value is a valid email.
func NewEmailSliceValidator() Validator {
	return NewStrSliceValidator(NewEmailValidator())
}

// NewAddressValidator returns a validator to validate whether an address is
// like "host:port", "host%zone:port", "[host]:port" or "[host%zone]:port".
//
// This validator uses net.SplitHostPort() to validate it.
func NewAddressValidator() Validator {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errNotString
		}

		if _, _, err := net.SplitHostPort(s); err != nil {
			return err
		}
		return nil
	}
}

// NewAddressSliceValidator returns a validator to validate whether the string element
// of the []string value is a valid address.
func NewAddressSliceValidator() Validator {
	return NewStrSliceValidator(NewAddressValidator())
}

// NewIntegerRangeValidator returns a validator to validate whether the integer
// value is between the min and the max.
//
// This validator can be used to validate the value of the type int, int8,
// int16, int32, int64, uint, uint8, uint16, uint32, uint64.
func NewIntegerRangeValidator(min, max int64) Validator {
	return func(value interface{}) error {
		v, err := ToInt64(value)
		if err != nil {
			return err
		}
		if min > v || v > max {
			return fmt.Errorf("the value '%d' is not between %d and %d", v, min, max)
		}
		return nil
	}
}

// NewFloatRangeValidator returns a validator to validate whether the float
// value is between the min and the max.
//
// This validator can be used to validate the value of the type float32 and
// float64.
func NewFloatRangeValidator(min, max float64) Validator {
	return func(value interface{}) error {
		f, err := ToFloat64(value)
		if err != nil {
			return err
		}

		if min > f || f > max {
			return fmt.Errorf("the value '%f' is not between %f and %f", f, min, max)
		}
		return nil
	}
}

// NewPortValidator returns a validator to validate whether a port is between
// 0 and 65535.
func NewPortValidator() Validator {
	return NewIntegerRangeValidator(0, 65535)
}
