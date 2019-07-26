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
	"time"
)

// VersionOpt reprensents a version option.
var VersionOpt = StrOpt("version", "Print the version and exit.").S("v").D("1.0.0").V(NewStrNotEmptyValidator())

// Opt is used to represent a option vlaue.
type Opt struct {
	// Name is the long name of the option.
	// It's necessary and must not be empty.
	Name string

	// Short is the short name of the option, which is optional.
	//
	// It should be a single-character string, such as "v" for "version".
	Short string

	// Help is the help or usage information, which is optional.
	Help string

	// Default is the default value of the option, which is necessary
	// and must not be nil. It will be used to indicate the type of the option.
	Default interface{}

	// Parser is used to parse the input option value to a specific type,
	// Which is necessary.
	//
	// Notice: it must not panic.
	Parser func(input interface{}) (output interface{}, err error)

	// Validators is the validators of the option, which is optional.
	//
	// When updating the option value, the validators will validate it.
	// If there is a validator returns an error, the updating fails and
	// returns the error. That's, these validators are the AND relations.
	Validators []Validator
}

func (o Opt) check() {
	if o.Name == "" {
		panic("the option name must not be empty")
	} else if o.Default == nil {
		panic(fmt.Errorf("the default value of the option '%s' must not be nil", o.Name))
	} else if o.Parser == nil {
		panic(fmt.Errorf("the parser of the option '%s' must not be nil", o.Name))
	} else if len(o.Short) > 1 {
		panic(fmt.Errorf("the short name of the option '%s' is more than one character", o.Name))
	}
}

func (o Opt) validate(value interface{}) (err error) {
	for _, validator := range o.Validators {
		if err = validator(value); err != nil {
			return err
		}
	}
	return
}

// N returns a new Opt with the given name based on the current option.
func (o Opt) N(name string) Opt {
	if name == "" {
		panic("the option name must not be empty")
	}
	o.Name = name
	return o
}

// S returns a new Opt with the given short name based on the current option.
func (o Opt) S(shortName string) Opt {
	if len(shortName) > 1 {
		panic("the short name of the option is not a single character")
	}
	o.Short = shortName
	return o
}

// H returns a new Opt with the given help based on the current option.
func (o Opt) H(help string) Opt {
	o.Help = help
	return o
}

// D returns a new Opt with the given default value based on the current option.
func (o Opt) D(_default interface{}) Opt {
	if _default == nil {
		panic("the default value of the option must not be nil")
	}
	o.Default = _default
	return o
}

// P returns a new Opt with the given parser based on the current option.
func (o Opt) P(parser func(interface{}) (interface{}, error)) Opt {
	if parser == nil {
		panic("the parser of the option must not be nil")
	}
	o.Parser = parser
	return o
}

// V returns a new Opt with the given validators based on the current option,
// which will append them.
func (o Opt) V(validators ...Validator) Opt {
	o.Validators = append(o.Validators, validators...)
	return o
}

// NewOpt returns a new Opt.
func NewOpt(name string, _default interface{}, parser func(interface{}) (interface{}, error)) Opt {
	return Opt{}.N(name).D(_default).P(parser)
}

// BoolOpt returns a bool Opt, which is equal to
//   NewOpt(name, false, ToBool).H(help)
func BoolOpt(name string, help string) Opt {
	return NewOpt(name, false, func(v interface{}) (interface{}, error) { return ToBool(v) }).H(help)
}

// IntOpt returns a int Opt, which is equal to
//   NewOpt(name, 0, ToInt).H(help)
func IntOpt(name string, help string) Opt {
	return NewOpt(name, 0, func(v interface{}) (interface{}, error) { return ToInt(v) }).H(help)
}

// Int32Opt returns a int32 Opt, which is equal to
//   NewOpt(name, int32(0), ToInt32).H(help)
func Int32Opt(name string, help string) Opt {
	return NewOpt(name, int32(0), func(v interface{}) (interface{}, error) { return ToInt32(v) }).H(help)
}

// Int64Opt returns a int64 Opt, which is equal to
//   NewOpt(name, int64(0), ToInt64).H(help)
func Int64Opt(name string, help string) Opt {
	return NewOpt(name, int64(0), func(v interface{}) (interface{}, error) { return ToInt64(v) }).H(help)
}

// UintOpt returns a uint Opt, which is equal to
//   NewOpt(name, uint(0), ToUint).H(help)
func UintOpt(name string, help string) Opt {
	return NewOpt(name, uint(0), func(v interface{}) (interface{}, error) { return ToUint(v) }).H(help)
}

// Uint32Opt returns a uint32 Opt, which is equal to
//   NewOpt(name, uint32(0), ToUint32).H(help)
func Uint32Opt(name string, help string) Opt {
	return NewOpt(name, uint32(0), func(v interface{}) (interface{}, error) { return ToUint32(v) }).H(help)
}

// Uint64Opt returns a uint64 Opt, which is equal to
//   NewOpt(name, uint64(0), ToUint64).H(help)
func Uint64Opt(name string, help string) Opt {
	return NewOpt(name, uint64(0), func(v interface{}) (interface{}, error) { return ToUint64(v) }).H(help)
}

// Float64Opt returns a float64 Opt, which is equal to
//   NewOpt(name, 0.0, ToFloat64).H(help)
func Float64Opt(name string, help string) Opt {
	return NewOpt(name, 0.0, func(v interface{}) (interface{}, error) { return ToFloat64(v) }).H(help)
}

// StrOpt returns a string Opt, which is equal to
//   NewOpt(name, "", ToString).H(help)
func StrOpt(name string, help string) Opt {
	return NewOpt(name, "", func(v interface{}) (interface{}, error) { return ToString(v) }).H(help)
}

// DurationOpt returns a time.Duration Opt, which is equal to
//   NewOpt(name, time.Duration(0), ToDuration).H(help)
func DurationOpt(name string, help string) Opt {
	return NewOpt(name, time.Duration(0), func(v interface{}) (interface{}, error) { return ToDuration(v) }).H(help)
}

// TimeOpt returns a time.Time Opt, which is equal to
//   NewOpt(name, time.Time{}, ToTime).H(help)
func TimeOpt(name string, help string) Opt {
	return NewOpt(name, time.Time{}, func(v interface{}) (interface{}, error) { return ToTime(v) }).H(help)
}

// StrSliceOpt returns a []string Opt, which is equal to
//   NewOpt(name, []string{}, ToStringSlice).H(help)
func StrSliceOpt(name string, help string) Opt {
	return NewOpt(name, []string{}, func(v interface{}) (interface{}, error) { return ToStringSlice(v) }).H(help)
}

// IntSliceOpt returns a []int Opt, which is equal to
//   NewOpt(name, []int{}, ToIntSlice).H(help)
func IntSliceOpt(name string, help string) Opt {
	return NewOpt(name, []int{}, func(v interface{}) (interface{}, error) { return ToIntSlice(v) }).H(help)
}

// UintSliceOpt returns a []uint Opt, which is equal to
//   NewOpt(name, []uint{}, ToUintSlice).H(help)
func UintSliceOpt(name string, help string) Opt {
	return NewOpt(name, []uint{}, func(v interface{}) (interface{}, error) { return ToUintSlice(v) }).H(help)
}

// Float64SliceOpt returns a []float64 Opt, which is equal to
//   NewOpt(name, []float64{}, ToFloat64Slice).H(help)
func Float64SliceOpt(name string, help string) Opt {
	return NewOpt(name, []float64{}, func(v interface{}) (interface{}, error) { return ToFloat64Slice(v) }).H(help)
}

// DurationSliceOpt returns a []time.Duration Opt, which is equal to
//   NewOpt(name, []time.Duration{}, ToDurationSlice).H(help)
func DurationSliceOpt(name string, help string) Opt {
	return NewOpt(name, []time.Duration{}, func(v interface{}) (interface{}, error) { return ToDurationSlice(v) }).H(help)
}
