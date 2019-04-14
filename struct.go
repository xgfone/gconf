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

// StructValidator is used to validate the struct value.
type StructValidator interface {
	Validate() error
}

// RegisterStruct retusters the struct fields as the options.
//
// The tag of the field supports "name"(string), "short"(string), "help"(string),
// "default"(string), "group"(string), "cmd"(string) and "cli"(bool).
//
//   1. "name", "short", "default" and "help" are used to create a option
//      with the name, the short name, the default value and the help doc.
//   2. "cli" is used to indicate whether the option is the CLI option or not.
//   3. "group" is used to change the group of the option to "group".
//   4. "cmd" is used to indicate that the option belongs the command named
//      by "cmd".
//
// If "name" is "-", that's `name:"-"`, the corresponding field will be ignored.
//
// The bool value may be one of "t", "T", "1", "on", "On", "ON", "true", "True",
// "TRUE", "yes", "Yes", "YES" for true and "f", "F", "0", "off", "Off", "OFF",
// "false", "False", "FALSE", "no", "No", "NO", "" for false.
//
// For the field that is a struct, it is a new sub-group of the current group,
// and the lower-case of the field name is the name of the new sub-group.
// But you can use the tag "group" or "name" to overwrite it, and "group" is
// preference to "name".
//
// Notice: All the tags are optional.
//
// If the struct has implemented the interface StructValidator, this validator
// will be called automatically after having parsed. But the individual field
// does not support the validator by the tag. You maybe choose others, such as
// github.com/asaskevich/govalidator. Beside, you can get the option from the
// corresponding group, then add the validators for it by hand. All of the
// builtin options support the validator chain.
//
// Notice:
//   1. The struct must be a pointer to a struct variable, or panic.
//   2. It must be called before being parsed, or panic.
//   3. The struct supports the nested struct, but not the pointer field.
//   4. For the struct option, you shouldn't call SetOptValue() because of concurrence.
//
func (c *Config) RegisterStruct(s interface{}) *Config {
	return c.registerStruct(false, s)
}

func (c *Config) registerStruct(cli bool, s interface{}) *Config {
	sv := reflect.ValueOf(s)
	if sv.IsNil() || !sv.IsValid() {
		panic("the struct is invalid or can't be set")
	}

	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	if sv.Kind() != reflect.Struct {
		panic("the struct is not a struct")
	}

	c.registerStructByValue(nil, c.OptGroup, sv, cli)

	if v, ok := s.(StructValidator); ok {
		c.validators = append(c.validators, v.Validate)
	}
	return c
}

func (c *Config) registerStructByValue(command *Command, optGroup *OptGroup, sv reflect.Value, cli bool) {
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	st := sv.Type()

	// Register the field as the option
	var err error
	num := sv.NumField()
	for i := 0; i < num; i++ {
		field := st.Field(i)
		fieldV := sv.Field(i)
		group := optGroup

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

		// Parse the tag "cli": the option is CLI or not.
		isCli := cli
		if _cli := strings.TrimSpace(field.Tag.Get("cli")); _cli != "" {
			if isCli, err = ToBool(_cli); err != nil {
				panic(fmt.Errorf("invalid value '%s' for cli", field.Tag.Get("cli")))
			}
		}

		// Parse the tag "help": the help document.
		help := strings.TrimSpace(field.Tag.Get("help"))

		// Parse the tag "cmd": the command.
		var cmd *Command
		if _cmd := strings.TrimSpace(field.Tag.Get("cmd")); _cmd != "" {
			if command == nil {
				cmd = c.NewCommand(_cmd, help)
			} else {
				cmd = command.NewCommand(_cmd, help)
			}

			isCli = true
			group = cmd.OptGroup
		}

		// Parse the tag "group": rename the group name.
		gname := strings.TrimSpace(field.Tag.Get("group"))
		if gname != "" {
			group = group.NewGroup(gname)
		}

		// Check whether the field is the struct.
		if t := field.Type.Kind(); t == reflect.Struct {
			if _, ok := fieldV.Interface().(time.Time); !ok { // For struct config
				if cmd == nil && gname == "" {
					group = group.NewGroup(name)
				}
				if cmd == nil && command != nil {
					cmd = command
				}
				c.registerStructByValue(cmd, group, fieldV, isCli)
				continue
			}
		}

		_type := getOptType(fieldV)
		if _type == int64Type {
			if _, ok := fieldV.Interface().(time.Duration); ok {
				_type = durationType
			}
		}

		// Parse the tag "default": the default value of the option.
		var _default interface{}
		if v, ok := field.Tag.Lookup("default"); ok {
			if _default, err = parseOpt(strings.TrimSpace(v), _type); err != nil {
				panic(fmt.Errorf("can't parse the default tag in the field %s: %s",
					field.Name, err))
			}
		}

		// Parse the tag "short": the short name of the option.
		short := strings.TrimSpace(field.Tag.Get("short"))

		group.registerOpt(isCli, newBaseOpt(short, name, _default, help, _type))
		group.fields[name] = fieldV
	}
}
