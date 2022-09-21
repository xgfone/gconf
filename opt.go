// Copyright 2021 xgfone
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

type optsT []Opt

func (a optsT) Len() int           { return len(a) }
func (a optsT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a optsT) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Parser is used to parse the option value intput.
type Parser func(input interface{}) (output interface{}, err error)

// Opt is used to represent a option vlaue.
type Opt struct {
	// Name is the long name of the option, which should be lower case.
	//
	// Required!
	Name string

	// Default is the default value of the option, which will be used to
	// indicate the type of the option.
	//
	// Required!
	Default interface{}

	// Parser is used to parse the input option value to a specific type.
	//
	// Required!
	Parser Parser

	// Short is the short name of the option, which should be a single-
	// character string, such as "v" for "version".
	//
	// Optional?
	Short string

	// Help is the help or usage information.
	//
	// Optional?
	Help string

	// IsCli indicates whether the option can be used for the CLI flag.
	//
	// Optional?
	IsCli bool

	// The list of the aliases of the option.
	//
	// Optional?
	Aliases []string

	// Validators is used to validate whether the option value is valid
	// after parsing it and before updating it.
	//
	// Optional?
	Validators []Validator

	// OnUpdate is called when the option value is updated.
	OnUpdate func(oldValue, newValue interface{})
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

// As returns a new Opt with the new aliases based on the current option,
// which will append them.
func (o Opt) As(aliases ...string) Opt {
	o.Aliases = append(o.Aliases, aliases...)
	return o
}

// Cli returns a new Opt with the cli flag based on the current option.
func (o Opt) Cli(cli bool) Opt {
	o.IsCli = cli
	return o
}

// N returns a new Opt with the given name based on the current option.
func (o Opt) N(name string) Opt {
	if name == "" {
		panic("the option name must not be empty")
	}
	o.Name = name
	return o
}

// H returns a new Opt with the given help based on the current option.
func (o Opt) H(help string) Opt {
	o.Help = help
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

// V returns a new Opt with the given validators based on the current option,
// which will append them.
func (o Opt) V(validators ...Validator) Opt {
	o.Validators = append(o.Validators, validators...)
	return o
}

// D returns a new Opt with the given default value based on the current option.
func (o Opt) D(_default interface{}) Opt {
	if _default == nil {
		panic("the default value of the option must not be nil")
	}

	if o.Parser == nil {
		o.Default = _default
	} else if value, err := o.Parser(_default); err != nil {
		panic(fmt.Errorf("fail to parse '%v' for the option named '%s': %s",
			_default, o.Name, err))
	} else {
		o.Default = value
	}
	return o
}

// P returns a new Opt with the given parser based on the current option.
func (o Opt) P(parser Parser) Opt {
	if parser == nil {
		panic("the parser of the option must not be nil")
	}
	o.Parser = parser
	if o.Default != nil {
		if value, err := o.Parser(o.Default); err != nil {
			panic(fmt.Errorf("fail to parse '%v' for the option named '%s': %s",
				o.Default, o.Name, err))
		} else {
			o.Default = value
		}
	}
	return o
}

// U returns a new Opt with the update callback function on the current option.
func (o Opt) U(callback func(oldValue, newValue interface{})) Opt {
	o.OnUpdate = callback
	return o
}

// NewOpt returns a new Opt that IsCli is true.
func NewOpt(name, help string, _default interface{}, parser Parser) Opt {
	return Opt{IsCli: true, Name: name, Help: help}.D(_default).P(parser)
}

// BoolOpt is the same NewOpt, but uses ToBool to parse the value as bool.
func BoolOpt(name string, help string) Opt {
	return NewOpt(name, help, false, func(v interface{}) (interface{}, error) {
		return ToBool(v)
	})
}

// IntOpt is the same NewOpt, but uses ToInt to parse the value as int.
func IntOpt(name string, help string) Opt {
	return NewOpt(name, help, 0, func(v interface{}) (interface{}, error) {
		return ToInt(v)
	})
}

// Int16Opt is the same NewOpt, but uses ToInt16 to parse the value as int16.
func Int16Opt(name string, help string) Opt {
	return NewOpt(name, help, int16(0),
		func(v interface{}) (interface{}, error) {
			return ToInt16(v)
		})
}

// Int32Opt is the same NewOpt, but uses ToInt32 to parse the value as int32.
func Int32Opt(name string, help string) Opt {
	return NewOpt(name, help, int32(0),
		func(v interface{}) (interface{}, error) {
			return ToInt32(v)
		})
}

// Int64Opt is the same NewOpt, but uses ToInt64 to parse the value as int64.
func Int64Opt(name string, help string) Opt {
	return NewOpt(name, help, int64(0),
		func(v interface{}) (interface{}, error) {
			return ToInt64(v)
		})
}

// UintOpt is the same NewOpt, but uses ToUint to parse the value as uint.
func UintOpt(name string, help string) Opt {
	return NewOpt(name, help, uint(0),
		func(v interface{}) (interface{}, error) {
			return ToUint(v)
		})
}

// Uint16Opt is the same NewOpt, but uses ToUint16 to parse the value as uint16.
func Uint16Opt(name string, help string) Opt {
	return NewOpt(name, help, uint16(0),
		func(v interface{}) (interface{}, error) {
			return ToUint16(v)
		})
}

// Uint32Opt is the same NewOpt, but uses ToUint32 to parse the value as uint32.
func Uint32Opt(name string, help string) Opt {
	return NewOpt(name, help, uint32(0),
		func(v interface{}) (interface{}, error) {
			return ToUint32(v)
		})
}

// Uint64Opt is the same NewOpt, but uses ToUint64 to parse the value as uint64.
func Uint64Opt(name string, help string) Opt {
	return NewOpt(name, help, uint64(0),
		func(v interface{}) (interface{}, error) {
			return ToUint64(v)
		})
}

// Float64Opt is the same NewOpt, but uses ToFloat64
// to parse the value as float64.
func Float64Opt(name string, help string) Opt {
	return NewOpt(name, help, 0.0, func(v interface{}) (interface{}, error) {
		return ToFloat64(v)
	})
}

// StrOpt is the same NewOpt, but uses ToString to parse the value as string.
func StrOpt(name string, help string) Opt {
	return NewOpt(name, help, "", func(v interface{}) (interface{}, error) {
		return ToString(v)
	})
}

// DurationOpt is the same NewOpt, but uses ToDuration
// to parse the value as time.Duration.
func DurationOpt(name string, help string) Opt {
	return NewOpt(name, help, time.Duration(0),
		func(v interface{}) (interface{}, error) {
			return ToDuration(v)
		})
}

// TimeOpt is the same NewOpt, but uses ToTime to parse the value as time.Time.
func TimeOpt(name string, help string) Opt {
	return NewOpt(name, help, time.Time{},
		func(v interface{}) (interface{}, error) {
			return ToTime(v)
		})
}

// StrSliceOpt is the same NewOpt, but uses ToStringSlice
// to parse the value as []string.
func StrSliceOpt(name string, help string) Opt {
	return NewOpt(name, help, []string{},
		func(v interface{}) (interface{}, error) {
			return ToStringSlice(v)
		})
}

// IntSliceOpt is the same NewOpt, but uses ToIntSlice
// to parse the value as []int.
func IntSliceOpt(name string, help string) Opt {
	return NewOpt(name, help, []int{},
		func(v interface{}) (interface{}, error) {
			return ToIntSlice(v)
		})
}

// UintSliceOpt is the same NewOpt, but uses ToUintSlice
// to parse the value as []uint.
func UintSliceOpt(name string, help string) Opt {
	return NewOpt(name, help, []uint{},
		func(v interface{}) (interface{}, error) {
			return ToUintSlice(v)
		})
}

// Float64SliceOpt is the same NewOpt, but uses ToFloat64Slice
// to parse the value as []float64.
func Float64SliceOpt(name string, help string) Opt {
	return NewOpt(name, help, []float64{},
		func(v interface{}) (interface{}, error) {
			return ToFloat64Slice(v)
		})
}

// DurationSliceOpt is the same NewOpt, but uses ToDurationSlice
// to parse the value as []time.Duration.
func DurationSliceOpt(name string, help string) Opt {
	return NewOpt(name, help, []time.Duration{},
		func(v interface{}) (interface{}, error) {
			return ToDurationSlice(v)
		})
}
