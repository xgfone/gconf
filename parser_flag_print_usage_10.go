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

// +build go1.10

package gconf

import (
	"flag"
	"fmt"
)

func printDefaultFlagUsage(fset *flag.FlagSet) {
	if fset.Name() == "" {
		fmt.Fprintf(fset.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(fset.Output(), "Usage of %s:\n", fset.Name())
	}
	PrintFlagUsage(fset.Output(), fset, false)
}
