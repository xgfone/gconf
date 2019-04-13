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
	"reflect"
	"time"
)

// Opt stands for an option value.
type Opt interface {
	// Name returns the long name of the option.
	// It's necessary and must not be empty.
	Name() string

	// Short returns the short name of the option.
	// It's optional. If having no short name, it should return "".
	Short() string

	// Help returns the help or usage information.
	// If having no help doc, it should return "".
	Help() string

	// Default returns the default value.
	// If having no default value, it should return nil.
	Default() interface{}

	// Zero returns the zero value of this type.
	//
	// For the slice, it should use the empty slice instead of nil.
	Zero() interface{}

	// Parse parses the argument to the type of this option.
	// If failed to parse, it should return an error to explain the reason.
	Parse(interface{}) (interface{}, error)
}

type optType int

func (ot optType) String() string {
	return optTypeMap[ot]
}

const (
	noneType optType = iota
	boolType
	stringType
	intType
	int8Type
	int16Type
	int32Type
	int64Type
	uintType
	uint8Type
	uint16Type
	uint32Type
	uint64Type
	float32Type
	float64Type
	durationType
	timeType

	stringsType
	intsType
	int64sType
	uintsType
	uint64sType
	float64sType
	durationsType
	timesType
)

var optTypeMap = map[optType]string{
	noneType:     "none",
	boolType:     "bool",
	stringType:   "string",
	intType:      "int",
	int8Type:     "int8",
	int16Type:    "int16",
	int32Type:    "int32",
	int64Type:    "int64",
	uintType:     "uint",
	uint8Type:    "uint8",
	uint16Type:   "uint16",
	uint32Type:   "uint32",
	uint64Type:   "uint64",
	float32Type:  "float32",
	float64Type:  "float64",
	durationType: "time.Duration",
	timeType:     "time.Time",

	stringsType:   "[]string",
	intsType:      "[]int",
	int64sType:    "[]int64",
	uintsType:     "[]uint",
	uint64sType:   "[]uint64",
	float64sType:  "[]float64",
	durationsType: "[]time.Duration",
	timesType:     "[]time.Time",
}

var kind2optType = map[reflect.Kind]optType{
	reflect.Bool:    boolType,
	reflect.String:  stringType,
	reflect.Int:     intType,
	reflect.Int8:    int8Type,
	reflect.Int16:   int16Type,
	reflect.Int32:   int32Type,
	reflect.Int64:   int64Type,
	reflect.Uint:    uintType,
	reflect.Uint8:   uint8Type,
	reflect.Uint16:  uint16Type,
	reflect.Uint32:  uint32Type,
	reflect.Uint64:  uint64Type,
	reflect.Float32: float32Type,
	reflect.Float64: float64Type,
}

func getOptType(v reflect.Value) optType {
	if t, ok := kind2optType[v.Kind()]; ok {
		return t
	}

	switch v.Interface().(type) {
	case time.Duration:
		return durationType
	case time.Time:
		return timeType
	case []string:
		return stringsType
	case []int:
		return intsType
	case []int64:
		return int64sType
	case []uint:
		return uintsType
	case []uint64:
		return uint64sType
	case []float64:
		return float64sType
	case []time.Duration:
		return durationsType
	case []time.Time:
		return timesType
	default:
		panic(fmt.Errorf("doesn't support the type %s", v.Type().Name()))
	}
}

type baseOpt struct {
	name     string
	help     string
	short    string
	_default interface{}

	_type      optType
	validators []Validator
}

var _ ValidatorChainOpt = baseOpt{}

func newBaseOpt(short, name string, _default interface{}, help string,
	optType optType) baseOpt {
	o := baseOpt{
		short:    short,
		name:     name,
		help:     help,
		_default: _default,
		_type:    optType,
	}
	o.Default()
	return o
}

// SetValidators resets the validator chain.
func (o baseOpt) SetValidators(vs ...Validator) ValidatorChainOpt {
	o.validators = vs
	return o
}

// AddValidators adds some new validators into the validator chain.
func (o baseOpt) AddValidators(vs ...Validator) ValidatorChainOpt {
	if len(o.validators) == 0 {
		o.validators = vs
	} else {
		for _, v := range vs {
			o.validators = append(o.validators, v)
		}
	}

	return o
}

// GetValidators returns the validator chain
func (o baseOpt) GetValidators() []Validator {
	return o.validators
}

// GetName returns the name of the option.
func (o baseOpt) Name() string {
	return o.name
}

// GetShort returns the shorthand name of the option.
func (o baseOpt) Short() string {
	return o.short
}

// GetHelp returns the help doc of the option.
func (o baseOpt) Help() string {
	return o.help
}

// GetDefault returns the default value of the option.
func (o baseOpt) Default() interface{} {
	if o._default == nil {
		return nil
	}

	switch o._type {
	case boolType:
		return o._default.(bool)
	case stringType:
		return o._default.(string)
	case intType:
		return o._default.(int)
	case int8Type:
		return o._default.(int8)
	case int16Type:
		return o._default.(int16)
	case int32Type:
		return o._default.(int32)
	case int64Type:
		return o._default.(int64)
	case uintType:
		return o._default.(uint)
	case uint8Type:
		return o._default.(uint8)
	case uint16Type:
		return o._default.(uint16)
	case uint32Type:
		return o._default.(uint32)
	case uint64Type:
		return o._default.(uint64)
	case float32Type:
		return o._default.(float32)
	case float64Type:
		return o._default.(float64)
	case durationType:
		return o._default.(time.Duration)
	case timeType:
		return o._default.(time.Time)
	case durationsType:
		return o._default.([]time.Duration)
	case timesType:
		return o._default.([]time.Time)
	case stringsType:
		return o._default.([]string)
	case intsType:
		return o._default.([]int)
	case int64sType:
		return o._default.([]int64)
	case uintsType:
		return o._default.([]uint)
	case uint64sType:
		return o._default.([]uint64)
	case float64sType:
		return o._default.([]float64)
	default:
		panic(fmt.Errorf("don't support the type %s", o._type))
	}
}

// Zero returns the zero value of this type.
func (o baseOpt) Zero() interface{} {
	switch o._type {
	case boolType:
		return false
	case stringType:
		return ""
	case intType:
		return int(0)
	case int8Type:
		return int8(0)
	case int16Type:
		return int16(0)
	case int32Type:
		return int32(0)
	case int64Type:
		return int64(0)
	case uintType:
		return uint(0)
	case uint8Type:
		return uint8(0)
	case uint16Type:
		return uint16(0)
	case uint32Type:
		return uint32(0)
	case uint64Type:
		return uint64(0)
	case float32Type:
		return float32(0)
	case float64Type:
		return float64(0)
	case durationType:
		return time.Duration(0)
	case timeType:
		return time.Time{}
	case stringsType:
		return []string{}
	case intsType:
		return []int{}
	case int64sType:
		return []int64{}
	case uintsType:
		return []uint{}
	case uint64sType:
		return []uint64{}
	case float64sType:
		return []float64{}
	case durationsType:
		return []time.Duration{}
	case timesType:
		return []time.Time{}
	default:
		panic(fmt.Errorf("don't support the type %s", o._type))
	}
}

// Parse parses the value of the option to a certain type.
func (o baseOpt) Parse(data interface{}) (v interface{}, err error) {
	return parseOpt(data, o._type)
}

func parseOpt(data interface{}, _type optType) (v interface{}, err error) {
	switch _type {
	case boolType:
		return ToBool(data)
	case stringType:
		return ToString(data)
	case intType, int8Type, int16Type, int32Type, int64Type:
		v, err = ToInt64(data)
	case uintType, uint8Type, uint16Type, uint32Type, uint64Type:
		v, err = ToUint64(data)
	case float32Type, float64Type:
		v, err = ToFloat64(data)
	case durationType:
		switch arg := data.(type) {
		case time.Duration:
			return arg, nil
		case int:
			return time.Duration(arg), nil
		case int8:
			return time.Duration(arg), nil
		case int16:
			return time.Duration(arg), nil
		case int32:
			return time.Duration(arg), nil
		case int64:
			return time.Duration(arg), nil
		case uint:
			return time.Duration(arg), nil
		case uint8:
			return time.Duration(arg), nil
		case uint16:
			return time.Duration(arg), nil
		case uint32:
			return time.Duration(arg), nil
		case uint64:
			return time.Duration(arg), nil
		case string:
			return time.ParseDuration(arg)
		default:
			return nil, fmt.Errorf("don't support the type '%s' for time.Duration", _type)
		}
	case timeType:
		switch arg := data.(type) {
		case time.Time:
			return arg, nil
		case string, []byte:
			return ToTime(arg)
		default:
			return nil, fmt.Errorf("don't support the type '%s' for time.Time", _type)
		}
	case stringsType:
		return ToStringSlice(data)
	case intsType:
		return ToIntSlice(data)
	case int64sType:
		return ToInt64Slice(data)
	case uintsType:
		return ToUintSlice(data)
	case uint64sType:
		return ToUint64Slice(data)
	case float64sType:
		return ToFloat64Slice(data)
	case durationsType:
		return ToDurations(data)
	case timesType:
		return ToTimes(time.RFC3339Nano, data)
	default:
		err = fmt.Errorf("don't support the type '%s'", _type)
	}

	if err != nil {
		return
	}

	switch _type {
	// case uint64Type:
	// case int64Type:
	// case float64Type:
	case intType:
		v = int(v.(int64))
	case int8Type:
		v = int8(v.(int64))
	case int16Type:
		v = int16(v.(int64))
	case int32Type:
		v = int32(v.(int64))
	case uintType:
		v = uint(v.(uint64))
	case uint8Type:
		v = uint8(v.(uint64))
	case uint16Type:
		v = uint16(v.(uint64))
	case uint32Type:
		v = uint32(v.(uint64))
	case float32Type:
		v = float32(v.(float64))
	}
	return
}

// BoolOpt return a new bool option.
func BoolOpt(short, name string, _default bool, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, boolType)
}

// StrOpt return a new string option.
func StrOpt(short, name string, _default string, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, stringType)
}

// IntOpt return a new int option.
func IntOpt(short, name string, _default int, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, intType)
}

// Int8Opt return a new int8 option.
func Int8Opt(short, name string, _default int8, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, int8Type)
}

// Int16Opt return a new int16 option.
func Int16Opt(short, name string, _default int16, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, int16Type)
}

// Int32Opt return a new int32 option.
func Int32Opt(short, name string, _default int32, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, int32Type)
}

// Int64Opt return a new int64 option.
func Int64Opt(short, name string, _default int64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, int64Type)
}

// UintOpt return a new uint option.
func UintOpt(short, name string, _default uint, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uintType)
}

// Uint8Opt return a new uint8 option.
func Uint8Opt(short, name string, _default uint8, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uint8Type)
}

// Uint16Opt return a new uint16 option.
func Uint16Opt(short, name string, _default uint16, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uint16Type)
}

// Uint32Opt return a new uint32 option.
func Uint32Opt(short, name string, _default uint32, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uint32Type)
}

// Uint64Opt return a new uint64 option.
func Uint64Opt(short, name string, _default uint64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uint64Type)
}

// Float32Opt return a new float32 option.
func Float32Opt(short, name string, _default float32, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, float32Type)
}

// Float64Opt return a new float64 option.
func Float64Opt(short, name string, _default float64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, float64Type)
}

// DurationOpt return a new time.Duration option.
//
// For the string value, it will use time.ParseDuration to parse it.
func DurationOpt(short, name string, _default time.Duration, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, durationType)
}

// TimeOpt return a new time.Time option.
//
// For the string value, it will be parsed by the layout time.RFC3339Nano.
func TimeOpt(short, name string, _default time.Time, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, timeType)
}

// DurationsOpt return a new []time.Duration option.
//
// For the string value, it will use time.ParseDuration to parse it.
func DurationsOpt(short, name string, _default []time.Duration, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, durationsType)
}

// TimesOpt return a new []time.Time option.
//
// For the string value, it will be parsed by the layout time.RFC3339Nano.
func TimesOpt(short, name string, _default []time.Time, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, timesType)
}

// StringsOpt return a new []string option.
func StringsOpt(short, name string, _default []string, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, float64sType)
}

// IntsOpt return a new []int option.
func IntsOpt(short, name string, _default []int, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, intsType)
}

// Int64sOpt return a new []int64 option.
func Int64sOpt(short, name string, _default []int64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, int64sType)
}

// UintsOpt return a new []uint option.
func UintsOpt(short, name string, _default []uint, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uintsType)
}

// Uint64sOpt return a new []uint64 option.
func Uint64sOpt(short, name string, _default []uint64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, uint64sType)
}

// Float64sOpt return a new []float64 option.
func Float64sOpt(short, name string, _default []float64, help string) ValidatorChainOpt {
	return newBaseOpt(short, name, _default, help, float64sType)
}

///////////////////////////////////////////////////////////////////////////////

// Bool is equal to BoolOpt("", name, _default, help).
func Bool(name string, _default bool, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, boolType)
}

// Str is equal to StrOpt("", name, _default, help).
func Str(name string, _default string, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, stringType)
}

// Int is equal to IntOpt("", name, _default, help).
func Int(name string, _default int, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, intType)
}

// Int8 is equal to Int8Opt("", name, _default, help).
func Int8(name string, _default int8, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, int8Type)
}

// Int16 is equal to Int16Opt("", name, _default, help).
func Int16(name string, _default int16, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, int16Type)
}

// Int32 is equal to Int32Opt("", name, _default, help).
func Int32(name string, _default int32, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, int32Type)
}

// Int64 is equal to Int64Opt("", name, _default, help).
func Int64(name string, _default int64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, int64Type)
}

// Uint is equal to UintOpt("", name, _default, help).
func Uint(name string, _default uint, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uintType)
}

// Uint8 is equal to Uint8Opt("", name, _default, help).
func Uint8(name string, _default uint8, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uint8Type)
}

// Uint16 is equal to Uint16Opt("", name, _default, help).
func Uint16(name string, _default uint16, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uint16Type)
}

// Uint32 is equal to Uint32Opt("", name, _default, help).
func Uint32(name string, _default uint32, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uint32Type)
}

// Uint64 is equal to Uint64Opt("", name, _default, help).
func Uint64(name string, _default uint64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uint64Type)
}

// Float32 is equal to Float32Opt("", name, _default, help).
func Float32(name string, _default float32, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, float32Type)
}

// Float64 is equal to Float64Opt("", name, _default, help).
func Float64(name string, _default float64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, float64Type)
}

// Duration is equal to DurationOpt("", name, _default, help).
func Duration(name string, _default time.Duration, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, durationType)
}

// Time is equal to TimeOpt("", name, _default, help).
func Time(short, name string, _default time.Time, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, timeType)
}

// Durations is equal to DurationsOpt("", name, _default, help).
func Durations(short, name string, _default []time.Duration, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, durationsType)
}

// Times is equal to TimesOpt("", name, _default, help).
func Times(short, name string, _default []time.Time, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, timesType)
}

// Strings is equal to StringsOpt("", name, _default, help).
func Strings(name string, _default []string, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, stringsType)
}

// Ints is equal to IntsOpt("", name, _default, help).
func Ints(name string, _default []int, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, intsType)
}

// Int64s is equal to Int64sOpt("", name, _default, help).
func Int64s(name string, _default []int64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, int64sType)
}

// Uints is equal to UintsOpt("", name, _default, help).
func Uints(name string, _default []uint, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uintsType)
}

// Uint64s is equal to Uint64sOpt("", name, _default, help).
func Uint64s(name string, _default []uint64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, uint64sType)
}

// Float64s is equal to Float64sOpt("", name, _default, help).
func Float64s(name string, _default []float64, help string) ValidatorChainOpt {
	return newBaseOpt("", name, _default, help, float64sType)
}
