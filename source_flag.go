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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PrintFlagUsage prints the flag usage instead of the default.
func PrintFlagUsage(flagSet *flag.FlagSet) {
	flagSet.VisitAll(func(f *flag.Flag) {
		// Two spaces before -; see next two comments.
		prefix := "  -"
		if len(f.Name) > 1 {
			prefix += "-"
		}

		s := fmt.Sprintf(prefix+"%s", f.Name)
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		} else {
			vf := reflect.ValueOf(f.Value)
			if vf.Kind() == reflect.Ptr {
				vf = vf.Elem()
			}
			if vf.Kind() == reflect.Bool {
				s += " bool"
			}
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)
		s += fmt.Sprintf(" (default: %q)", f.DefValue)
		fmt.Fprint(flagSet.Output(), s, "\n")
	})
}

// AddOptFlag adds the option to the flagSet, which is flag.CommandLine
// by default.
//
// Notice: for the slice option, it maybe occur many times, and they are
// combined with the comma as the string representation of slice. For example,
//
//   $APP --slice-opt v1  --slice-opt v2  --slice-opt v3
//   $APP --slice-opt v1,v2  --slice-opt v3
//   $APP --slice-opt v1,v2,v3
//
// They are equivalent.
func AddOptFlag(c *Config, flagSet ...*flag.FlagSet) {
	addAndParseOptFlag(false, c, flagSet...)
}

// AddAndParseOptFlag is the same as AddOptFlag, but parses the CLI arguments.
//
// Notice: if there is the version flag and it is true, it will print the version
// and exit.
func AddAndParseOptFlag(c *Config, flagSet ...*flag.FlagSet) error {
	return addAndParseOptFlag(true, c, flagSet...)
}

func addAndParseOptFlag(parse bool, c *Config, flagSet ...*flag.FlagSet) error {
	flagset := flag.CommandLine
	if len(flagSet) > 0 && flagSet[0] != nil {
		flagset = flagSet[0]
	}

	var vName, value string
	if opt := c.GetVersion(); opt.Name != "" && opt.Default != nil {
		flagset.Bool(opt.Name, false, opt.Help)
		vName = opt.Name
		value = opt.Default.(string)
	}

	flagset.Usage = func() { PrintFlagUsage(flagset) }
	for _, group := range c.AllGroups() {
		for _, opt := range group.AllOpts() {
			if !opt.Cli {
				continue
			}

			name := opt.Name
			if gname := group.Name(); gname != "" {
				name = fmt.Sprintf("%s.%s", gname, opt.Name)
			}
			name = strings.Replace(name, "_", "-", -1)

			switch v := opt.Default.(type) {
			case nil:
				flagset.String(name, "", opt.Help)
			case bool:
				flagset.Bool(name, v, opt.Help)
			case int, int8, int16, int32, int64:
				flagset.Int64(name, reflect.ValueOf(v).Int(), opt.Help)
			case uint, uint8, uint16, uint32, uint64:
				flagset.Uint64(name, reflect.ValueOf(v).Uint(), opt.Help)
			case float32, float64:
				flagset.Float64(name, reflect.ValueOf(v).Float(), opt.Help)
			case time.Duration:
				flagset.Duration(name, v, opt.Help)
			default:
				switch vf := reflect.ValueOf(opt.Default); vf.Kind() {
				case reflect.Slice, reflect.Array:
					sv := &flagSliceValue{values: make([]string, vf.Len())}
					for i, _len := 0, vf.Len(); i < _len; i++ {
						sv.values[i] = fmt.Sprint(vf.Index(i).Interface())
					}
					flagset.Var(sv, name, opt.Help)
				default:
					flagset.String(name, fmt.Sprintf("%v", v), opt.Help)
				}
			}
		}
	}

	if parse {
		if err := flagset.Parse(os.Args[1:]); err != nil {
			return err
		}

		if vName != "" {
			if flag := flagset.Lookup(vName); flag != nil {
				if yes, _ := strconv.ParseBool(flag.Value.String()); yes {
					fmt.Println(value)
					os.Exit(0)
				}
			}
		}
	}

	return nil
}

type flagSliceValue struct {
	values []string
	isset  bool
}

func (v *flagSliceValue) String() string {
	if v == nil {
		return ""
	}
	return strings.Join(v.values, ",")
}

func (v *flagSliceValue) Set(s string) error {
	if v != nil {
		if !v.isset {
			v.isset = true
			v.values = []string{s}
		} else {
			v.values = append(v.values, s)
		}
	}
	return nil
}

// NewFlagSource returns a new source based on flag.FlagSet,
// which is flag.CommandLine by default.
func NewFlagSource(flagSet ...*flag.FlagSet) Source {
	flagset := flag.CommandLine
	if len(flagSet) > 0 && flagSet[0] != nil {
		flagset = flagSet[0]
	}
	return flagSource{flagSet: flagset}
}

type flagSource struct {
	flagSet *flag.FlagSet
}

func (f flagSource) Watch(load func(DataSet, error) bool, exit <-chan struct{}) {}

func (f flagSource) Read() (DataSet, error) {
	if !f.flagSet.Parsed() {
		if err := f.flagSet.Parse(os.Args[1:]); err != nil {
			return DataSet{Source: "flag", Format: "json"}, err
		}
	}

	vs := make(map[string]string, 32)
	f.flagSet.Visit(func(f *flag.Flag) {
		vs[strings.Replace(f.Name, "-", "_", -1)] = f.Value.String()
	})

	data, err := json.Marshal(vs)
	if err != nil {
		return DataSet{Source: "flag", Format: "json"}, err
	}
	ds := DataSet{
		Args:      f.flagSet.Args(),
		Data:      data,
		Format:    "json",
		Source:    "flag",
		Timestamp: time.Now(),
	}
	ds.Checksum = fmt.Sprintf("md5:%s", ds.Md5())
	return ds, nil
}
