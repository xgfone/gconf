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

	// Cli indicates whether the option can be used for the CLI flag.
	Cli bool

	// Tags is the key-value metadata of the Opt.
	Tags map[string]string

	// The list of the aliases of the option, which will be registered to
	// the group that the option is registered when it is being registered.
	Aliases []string

	// Default is the default value of the option, which is necessary
	// and must not be nil. It will be used to indicate the type of the option.
	Default interface{}

	// Parser is used to parse the input option value to a specific type,
	// Which is necessary.
	//
	// Notice: it must not panic.
	Parser func(input interface{}) (output interface{}, err error)

	// Fix is used to fix the parsed value.
	//
	// The different between Parser and Fix:
	//   1. Parser only parses the value from the arbitrary type to a specific.
	//   2. Fix only changes the value, not the type, that's, input and output
	//      should be the same type. For example, input is the NIC name,
	//      and Fix can get the ip by the NIC name then return it as output.
	//      So it ensures that input may be NIC or IP, and that the value
	//      of the option is always a IP.
	Fix func(input interface{}) (output interface{}, err error)

	// Observers are called after the value of the option is updated.
	Observers []func(newValue interface{})

	// Validators is the validators of the option, which is optional.
	//
	// When updating the option value, the validators will validate it.
	// If there is a validator returns an error, the updating fails and
	// returns the error. That's, these validators are the AND relations.
	//
	// Notice: they must not panic.
	Validators []Validator

	fixDefault bool
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

// C returns a new Opt with the cli flag based on the current option.
func (o Opt) C(cli bool) Opt {
	o.Cli = cli
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
		panic(NewOptError("", o.Name, err, _default))
	} else {
		o.Default = value
	}
	return o
}

func (o *Opt) fix() error {
	if o.fixDefault && o.Fix != nil && o.Default != nil {
		_default, err := o.Fix(o.Default)
		if err != nil {
			return err
		}
		o.Default = _default
	}
	return nil
}

// F returns a new Opt with the given fix function based on the current option.
//
// If fixDefault is true, it will fix the default value when registering
// the option.
func (o Opt) F(fix func(interface{}) (interface{}, error), fixDefault ...bool) Opt {
	o.Fix = fix
	if len(fixDefault) > 0 {
		o.fixDefault = fixDefault[0]
	}

	return o
}

// H returns a new Opt with the given help based on the current option.
func (o Opt) H(help string) Opt {
	o.Help = help
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

// O returns a new Opt with the given observer based on the current option,
// which will append them.
func (o Opt) O(observers ...func(interface{})) Opt {
	o.Observers = append(o.Observers, observers...)
	return o
}

// P returns a new Opt with the given parser based on the current option.
func (o Opt) P(parser func(interface{}) (interface{}, error)) Opt {
	if parser == nil {
		panic("the parser of the option must not be nil")
	}
	o.Parser = parser
	if o.Default != nil {
		if value, err := o.Parser(o.Default); err != nil {
			panic(NewOptError("", o.Name, err, o.Default))
		} else {
			o.Default = value
		}
	}
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

// T returns a new Opt with the key-value tag based on the current option,
// which will clone the tags from the current option to the new.
//
// Notice: the key must not be empty, but value may be.
func (o Opt) T(key, value string) Opt {
	if key == "" {
		panic("the tag key must not be empty")
	}

	tags := make(map[string]string, len(o.Tags)*2)
	for k, v := range o.Tags {
		tags[k] = v
	}
	tags[key] = value
	o.Tags = tags
	return o
}

// V returns a new Opt with the given validators based on the current option,
// which will append them.
func (o Opt) V(validators ...Validator) Opt {
	o.Validators = append(o.Validators, validators...)
	return o
}

// NewOpt returns a new Opt.
//
// Notice: Cli is true by default.
func NewOpt(name string, _default interface{}, parser func(interface{}) (interface{}, error)) Opt {
	return Opt{Cli: true}.N(name).D(_default).P(parser)
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
