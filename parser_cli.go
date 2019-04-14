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
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type flagParser struct {
	utoh bool
	fset *flag.FlagSet
}

// NewDefaultFlagCliParser returns a new CLI parser based on flag,
// which is equal to NewFlagCliParser("", 0, underlineToHyphen, flag.CommandLine).
func NewDefaultFlagCliParser(underlineToHyphen ...bool) Parser {
	var u2h bool
	if len(underlineToHyphen) > 0 {
		u2h = underlineToHyphen[0]
	}
	return NewFlagCliParser(flag.CommandLine, u2h)
}

// NewFlagCliParser returns a new CLI parser based on flag.FlagSet.
//
// If flagSet is nil, it will create a default flag.FlagSet, which is equal to
//
//    flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
//
// If underlineToHyphen is true, it will convert the underline to the hyphen.
//
// Notice:
//   1. The flag parser does not support the commands, so will ignore them.
//   2. when other libraries use the default global flag.FlagSet, that's
//      flag.CommandLine, such as github.com/golang/glog, please use
//      flag.CommandLine as flag.FlagSet.
func NewFlagCliParser(flagSet *flag.FlagSet, underlineToHyphen bool) Parser {
	return &flagParser{
		fset: flagSet,
		utoh: underlineToHyphen,
	}
}

func (f *flagParser) Name() string {
	return "flag"
}

func (f *flagParser) Priority() int {
	return 0
}

func (f *flagParser) Pre(c *Config) error {
	if f.fset == nil {
		f.fset = flag.NewFlagSet(c.Name(), flag.ExitOnError)
	}

	if f.fset.Usage == nil {
		f.fset.Usage = func() { fmt.Println(c.Description()) }
	}
	return nil
}

func (f *flagParser) Post(c *Config) error {
	return nil
}

func (f *flagParser) Parse(c *Config) (err error) {
	// Convert the option name.
	name2group := make(map[string]string, 8)
	name2opt := make(map[string]string, 8)
	for _, group := range c.AllNotCommandGroups() {
		gname := group.FullName()
		for _, opt := range group.CliOpts() {
			name := opt.Name()
			if gname != c.GetDefaultGroupName() {
				name = fmt.Sprintf("%s%s%s", gname, c.GetGroupSeparator(), name)
			}

			if f.utoh {
				name = strings.Replace(name, "_", "-", -1)
			}

			name2group[name] = gname
			name2opt[name] = opt.Name()

			switch opt.Zero().(type) {
			case bool:
				var _default bool
				if v := opt.Default(); v != nil {
					_default = v.(bool)
				}
				f.fset.Bool(name, _default, opt.Help())
				c.Printf("[%s] Add the bool flag '%s'", f.Name(), name)
			case int, int8, int16, int32, int64:
				var _default int64
				if v := opt.Default(); v != nil {
					_default, _ = ToInt64(v)
				}
				f.fset.Int64(name, _default, opt.Help())
				c.Printf("[%s] Add the int flag '%s'", f.Name(), name)
			case uint, uint8, uint16, uint32, uint64:
				var _default uint64
				if v := opt.Default(); v != nil {
					_default, _ = ToUint64(v)
				}
				f.fset.Uint64(name, _default, opt.Help())
			case float32, float64:
				var _default float64
				if v := opt.Default(); v != nil {
					_default, _ = ToFloat64(v)
				}
				f.fset.Float64(name, _default, opt.Help())
				c.Printf("[%s] Add the float flag '%s'", f.Name(), name)
			case time.Duration:
				var _default time.Duration
				if v := opt.Default(); v != nil {
					_default = v.(time.Duration)
				}
				f.fset.Duration(name, _default, opt.Help())
				c.Printf("[%s] Add the time.Duration flag '%s'", f.Name(), name)
			default:
				var _default string
				if v := opt.Default(); v != nil {
					_default = fmt.Sprintf("%v", v)
				}
				f.fset.String(name, _default, opt.Help())
				c.Printf("[%s] Add the string flag '%s'", f.Name(), name)
			}
		}
	}

	// Register the version option.
	var _version *bool
	_, vname, version, vhelp := c.GetCliVersion()
	if version != "" {
		_version = f.fset.Bool(vname, false, vhelp)
		c.Printf("[%s] Add the version flag '%s'", f.Name(), vname)
	}

	// Parse the CLI arguments.
	if err = f.fset.Parse(c.ParsedCliArgs()); err != nil {
		return
	}

	if _version != nil && *_version {
		fmt.Println(version)
		os.Exit(0)
	}

	// Acquire the result.
	c.SetCliArgs(f.fset.Args())
	f.fset.Visit(func(fg *flag.Flag) {
		c.Printf("[%s] Parsing flag '%s'", f.Name(), fg.Name)
		gname := name2group[fg.Name]
		optname := name2opt[fg.Name]
		if gname != "" && optname != "" && fg.Name != vname {
			c.SetOptValue(0, gname, optname, fg.Value.String())
		}
	})

	return
}
