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

// +build !go1.10

package gconf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func printDefaultFlagUsage() {
	if name := filepath.Base(os.Args[0]); name == "" {
		fmt.Fprintf(os.Stderr, "Usage:\n")
	} else {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
	}
	PrintFlagUsage(flag.CommandLine, false)
}

// PrintFlagUsage prints the usage of flag.FlagSet, which is almost equal to
// flag.FlagSet.PrintDefaults(), but print the double prefixes "--"
// for the long name of the option.
func PrintFlagUsage(fset *flag.FlagSet, exceptDefault bool) {
	fset.VisitAll(func(_flag *flag.Flag) {
		// Two spaces before -; see next two comments.
		prefix := "  -"
		if len(_flag.Name) > 1 {
			prefix += "-"
		}

		s := fmt.Sprintf(prefix+"%s", _flag.Name)
		name, usage := flag.UnquoteUsage(_flag)
		if len(name) > 0 {
			s += " " + name
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

		if !exceptDefault || !isZeroValue(_flag, _flag.DefValue) {
			vf := reflect.ValueOf(_flag.Value)
			if vf.Kind() == reflect.Ptr {
				vf = vf.Elem()
			}
			if vf.Kind() == reflect.String {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", _flag.DefValue)
			} else {
				s += fmt.Sprintf(" (default %s)", _flag.DefValue)
			}
		}
		fmt.Fprint(os.Stderr, s, "\n")
	})
}
