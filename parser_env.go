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
	"os"
	"strings"
)

type envVarParser struct {
	prefix   string
	priority int
}

// NewEnvVarParser returns a new environment variable parser.
//
// For the environment variable name, it's the format "PREFIX_GROUP_OPTION".
// If the prefix is empty, it's "GROUP_OPTION". For the default group, it's
// "PREFIX_OPTION". When the prefix is empty and the group is the default,
// it's "OPTION". "GROUP" is the group name, and "OPTION" is the option name.
//
// Notice: the prefix, the group name and the option name will be converted to
// the upper, and the group separator will be converted to "_".
func NewEnvVarParser(priority int, prefix string) Parser {
	return envVarParser{prefix: prefix, priority: priority}
}

func (e envVarParser) Name() string {
	return "env"
}

func (e envVarParser) Priority() int {
	return e.priority
}

func (e envVarParser) Pre(c *Config) error {
	return nil
}

func (e envVarParser) Post(c *Config) error {
	return nil
}

func (e envVarParser) Parse(c *Config) (err error) {
	// Initialize the prefix
	prefix := e.prefix
	if prefix != "" {
		prefix += "_"
	}

	// Convert the option to the variable name
	env2opts := make(map[string][]string, len(c.AllGroups())*8)
	for _, group := range c.AllGroups() {
		var gname string
		if group.FullName() != c.GetDefaultGroupName() {
			gname = strings.Replace(group.FullName(), c.GetGroupSeparator(), "_", -1) + "_"
		}
		for _, opt := range group.AllOpts() {
			e := fmt.Sprintf("%s%s%s", prefix, gname, opt.Name())
			env2opts[strings.ToUpper(e)] = []string{group.FullName(), opt.Name()}
		}
	}

	// Get the option value from the environment variable.
	envs := os.Environ()
	for _, env := range envs {
		c.Printf("[%s] Parsing Env '%s'", e.Name(), env)
		items := strings.SplitN(env, "=", 2)
		if len(items) == 2 {
			if info, ok := env2opts[items[0]]; ok {
				if err = c.SetOptValue(e.priority, info[0], info[1], items[1]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
