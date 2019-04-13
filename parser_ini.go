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
	"io/ioutil"
	"strings"
	"unicode"
)

type iniParser struct {
	sep  string
	opt  string
	prio int
	init func(*Config) error
}

// NewSimpleIniParser returns a INI parser with the priority 100, which registers
// the option, optName, before parsing the option.
func NewSimpleIniParser(optName string) Parser {
	return NewIniParser(100, optName, func(c *Config) error {
		c.RegisterCliOpt(Str(optName, "", "The path of the INI config file."))
		return nil
	})
}

// NewIniParser returns a new ini parser based on the file.
//
// The first argument is used to customized the priority.
//
// The second argument is the option name which the parser needs. It will be
// registered, and parsed before this parser runs.
//
// The third argument sets the Init function.
//
// The ini parser supports the line comments starting with "#", "//" or ";".
// The key and the value is separated by an equal sign, that's =. The key must
// be in one of ., :, _, -, number and letter. If giving fmtKey, it can convert
// the key in the ini file to the new one.
//
// If the value ends with "\", it will continue the next line. The lines will
// be joined by "\n" together.
//
// Notice: the options that have not been assigned to a certain group will be
// divided into the default group.
func NewIniParser(priority int, optName string, init func(*Config) error) Parser {
	return iniParser{prio: priority, opt: optName, sep: "=", init: init}
}

func (p iniParser) Name() string {
	return "ini"
}

func (p iniParser) Priority() int {
	return p.prio
}

func (p iniParser) Pre(c *Config) error {
	if p.init != nil {
		return p.init(c)
	}
	return nil
}

func (p iniParser) Post(c *Config) error {
	return nil
}

func (p iniParser) Parse(c *Config) error {
	// Read the content of the config file.
	filename := c.StringD(p.opt, "")
	if filename == "" {
		return nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse the config file.
	gname := c.GetDefaultGroupName()
	lines := strings.Split(string(data), "\n")
	for index, maxIndex := 0, len(lines); index < maxIndex; {
		line := strings.TrimSpace(lines[index])
		index++

		c.Printf("[%s] Parsing %dth line: '%s'", p.Name(), index, line)

		// Ignore the empty line.
		if len(line) == 0 {
			continue
		}

		// Ignore the line comments starting with "#", ";" or "//".
		if (line[0] == '#') || (line[0] == ';') ||
			(len(line) > 1 && line[0] == '/' && line[1] == '/') {
			continue
		}

		// Start a new group
		if line[0] == '[' && line[len(line)-1] == ']' {
			gname = strings.TrimSpace(line[1 : len(line)-1])
			if gname == "" {
				return fmt.Errorf("the group is empty")
			}
			continue
		}

		n := strings.Index(line, p.sep)
		if n == -1 {
			return fmt.Errorf("the %dth line misses the separator '%s'", index, p.sep)
		}

		key := strings.TrimSpace(line[0:n])
		for _, r := range key {
			if r != '_' && r != '-' && !unicode.IsNumber(r) && !unicode.IsLetter(r) {
				return fmt.Errorf("invalid identifier key '%s'", key)
			}
		}
		value := strings.TrimSpace(line[n+len(p.sep) : len(line)])

		// The continuation line
		if value != "" && value[len(value)-1] == '\\' {
			vs := []string{strings.TrimSpace(strings.TrimRight(value, "\\"))}
			for index < maxIndex {
				value = strings.TrimSpace(lines[index])
				vs = append(vs, strings.TrimSpace(strings.TrimRight(value, "\\")))
				index++
				c.Printf("[%s] Parsing %dth line: '%s'", p.Name(), index, value)
				if value == "" || value[len(value)-1] != '\\' {
					break
				}
			}
			value = strings.TrimSpace(strings.Join(vs, "\n"))
		}

		if err = c.SetOptValue(p.prio, gname, key, value); err != nil {
			return err
		}
	}

	return nil
}
