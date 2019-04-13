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
)

// ConfigFileName is the name of the option to indicate the path
// of the configuration file.
const ConfigFileName = "config-file"

// Conf is the default global Config.
var Conf = NewDefault()

// NewDefault returns a new default Config.
//
// The default Config only contains the Flag parser and the INI parser.
func NewDefault(fset ...*flag.FlagSet) *Config {
	_fset := flag.CommandLine
	if len(fset) > 0 {
		_fset = fset[0]
	}

	cli := NewFlagCliParser(_fset, true)
	ini := NewSimpleIniParser(ConfigFileName)
	return New().AddParser(cli, ini)
}
