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
	"strings"
	"time"
)

// OptField is used to describe a struct field option, which can get or set
// the value of the struct field safely.
type OptField interface {
	// Default returns the default value of the field option.
	Default() interface{}

	// Parse converts the input to output.
	//
	// Notice: the type of output must be identical with the types of the value
	// returned by Default() and the argument of Set().
	Parse(input interface{}) (output interface{}, err error)

	// Set is used to update the value of the field, which must be goroutine-safe.
	Set(interface{})
}

// RegisterStruct retusters the struct fields as the options into the current group.
//
// Supproted types for the struct filed:
//
//   bool
//   int
//   int32
//   int64
//   uint
//   uint32
//   uint64
//   float64
//   string
//   time.Duration
//   time.Time
//   []int
//   []uint
//   []float64
//   []string
//   []time.Duration
//   OptField<interface{}>
//
// Other types will be ignored.
//
// The tag of the field supports "name"(string), "short"(string),
// "help"(string), "default"(string), "group"(string).
//
//   1. "name", "short", "default" and "help" are used to create a option
//      with the name, the short name, the default value and the help doc.
//   2. "group" is used to change the group of the option to "group".
//      For a struct, if no "group", it will use "name".
//
// If "name" or "group" is "-", that's `name:"-"` or `group:"-"`,
// the corresponding field will be ignored.
//
// The bool value will be parsed by `strconv.ParseBool`, so "1", "t", "T",
// "TRUE", "true", "True", "0", "f", "F", "FALSE", "false" and "False"
// are accepted.
//
// For the field that is a struct, it is a new sub-group of the current group,
// and the lower-case of the field name is the name of the new sub-group.
// But you can use the tag "group" or "name" to overwrite it, and "group" is
// preference to "name".
//
// Notice:
//   1. All the tags are optional.
//   2. OptField is goroutine-safe, but the others is not.
//   3. The struct must be a pointer to a struct variable, or panic.
//   4. The struct supports the nested struct, but not the pointer field
//      except OptField.
//
func (g *OptGroup) RegisterStruct(v interface{}) {
	if v == nil {
		panic("the struct value must not be nil")
	}

	sv := reflect.ValueOf(v)
	if !sv.IsValid() {
		panic("the struct is invalid or can't be set")
	}

	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	if sv.Kind() != reflect.Struct {
		panic("the struct is not a struct")
	}

	g.conf.registerStructByValue(sv, sv)
}

func (g *OptGroup) registerStructByValue(sv, orig reflect.Value) {
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	st := sv.Type()

	var optFieldValue OptField
	optFieldT := reflect.ValueOf(&optFieldValue).Elem().Type()

	// Register the field as the option
	num := sv.NumField()
	for i := 0; i < num; i++ {
		field := st.Field(i)
		fieldV := sv.Field(i)
		group := g

		// Check whether the field can be set.
		if !fieldV.CanSet() {
			continue
		}

		// Parse the tag "name": the option name.
		name := strings.ToLower(field.Name)
		tagname := strings.TrimSpace(field.Tag.Get("name"))
		if tagname == "-" {
			continue
		} else if tagname != "" {
			name = tagname
		}

		// Parse the tag "group": rename the group name.
		gname := strings.TrimSpace(field.Tag.Get("group"))
		if gname != "" {
			group = group.NewGroup(gname)
		}

		var opt Opt
		var setter func(interface{})

		fv := fieldV
		if kind := fv.Kind(); kind != reflect.Ptr && kind != reflect.Interface {
			fv = fv.Addr()
		}

		if fv.Type().Implements(optFieldT) { // For OptField
			if fv.IsNil() {
				continue
			}
			optField := fv.Interface().(OptField)
			setter = optField.Set
			opt = NewOpt(name, optField.Default(), optField.Parse)
		} else {
			if t := field.Type.Kind(); t == reflect.Struct {
				if _, ok := fieldV.Interface().(time.Time); !ok { // For struct config
					if gname == "" {
						group = group.NewGroup(name)
					}
					group.registerStructByValue(fieldV, orig)
					continue
				}
			}

			switch v := fieldV.Interface().(type) {
			case bool:
				opt = BoolOpt(name, "").D(v)
			case int:
				opt = IntOpt(name, "").D(v)
			case int32:
				opt = Int32Opt(name, "").D(v)
			case int64:
				opt = Int64Opt(name, "").D(v)
			case uint:
				opt = UintOpt(name, "").D(v)
			case uint32:
				opt = Uint32Opt(name, "").D(v)
			case uint64:
				opt = Uint64Opt(name, "").D(v)
			case float64:
				opt = Float64Opt(name, "").D(v)
			case string:
				opt = StrOpt(name, "").D(v)
			case time.Duration:
				opt = DurationOpt(name, "").D(v)
			case time.Time:
				opt = TimeOpt(name, "").D(v)
			case []int:
				opt = IntSliceOpt(name, "").D(v)
			case []uint:
				opt = UintSliceOpt(name, "").D(v)
			case []float64:
				opt = Float64SliceOpt(name, "").D(v)
			case []string:
				opt = StrSliceOpt(name, "").D(v)
			case []time.Duration:
				opt = DurationSliceOpt(name, "").D(v)
			default:
				continue
			}

			setter = func(value interface{}) { fieldV.Set(reflect.ValueOf(value)) }
		}

		// Parse the tag "help": the help document.
		if help := strings.TrimSpace(field.Tag.Get("help")); help != "" {
			opt.Help = help
		}

		// Parse the tag "short": the short name of the option.
		if short := strings.TrimSpace(field.Tag.Get("short")); short != "" {
			opt.Short = short
		}

		// Parse the tag "default": the default value of the option.
		if v, ok := field.Tag.Lookup("default"); ok {
			if _default, err := opt.Parser(strings.TrimSpace(v)); err != nil {
				panic(fmt.Errorf("can't parse the default tag in the field %s: %s", field.Name, err))
			} else {
				opt.Default = _default
				setter(_default)
			}
		}

		group.registerOpts(opt)
		group.setOptWatch(opt.Name, setter)
	}
}
