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
	parsed int32

	sep  string
	prio int

	init    func(*Config) error
	getData func(*Config) ([]byte, error)
}

// NewSimpleIniParser returns a INI parser with the priority 100,
// which registers the CLI option, cliOptName, into the default group and reads
// the data from the INI file appointed by cliOptName.
func NewSimpleIniParser(cliOptName string) Parser {
	return NewIniParser(100, func(c *Config) error {
		c.RegisterCliOpt(Str(cliOptName, "", "The path of the INI config file."))
		return nil
	}, func(c *Config) ([]byte, error) {
		// Read the content of the config file.
		if filename := c.StringD(cliOptName, ""); filename == "" {
			return nil, nil
		} else if data, err := ioutil.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return data, nil
		}
	})
}

// NewIniParser returns a new ini parser based on the file.
//
// The first argument sets the Init function.
//
// The second argument sets the Init function to initialize the parser, such as
// registering the CLI option.
//
// The third argument is used to read the data to be parsed, which will
// be called at the start of calling the method Parse().
//
// The ini parser supports the line comments starting with "#", "//" or ";".
// The key and the value is separated by an equal sign, that's =. The key must
// be in one of ., :, _, -, number and letter.
//
// If the value ends with "\", it will continue the next line. The lines will
// be joined by "\n" together.
func NewIniParser(priority int, init func(*Config) error, getData func(*Config) ([]byte, error)) Parser {
	return &iniParser{
		sep:  "=",
		prio: priority,

		init:    init,
		getData: getData,
	}
}

func (p *iniParser) Name() string {
	return "ini"
}

func (p *iniParser) Priority() int {
	return p.prio
}

func (p *iniParser) Pre(c *Config) error {
	if p.init != nil {
		return p.init(c)
	}
	return nil
}

func (p *iniParser) Post(c *Config) error {
	return nil
}

func (p *iniParser) Parse(c *Config) error {
	data, err := p.getData(c)
	if err != nil {
		return err
	} else if len(data) == 0 {
		return nil
	}

	// Parse the config file.
	opts := make([][3]interface{}, 0)
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
		if len(key) == 0 {
			return fmt.Errorf("empty identifier key")
		}
		for _, r := range key {
			if unicode.IsSpace(r) || !unicode.IsPrint(r) {
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

		if group := c.Group(gname); group == nil {
			continue
		} else if opt := group.Opt(key); opt == nil {
			continue
		} else if v, err := opt.Parse(value); err != nil {
			return err
		} else {
			opts = append(opts, [3]interface{}{group, key, v})
		}
	}

	for _, opt := range opts {
		opt[0].(*OptGroup).UpdateOptValue(opt[1].(string), opt[2])
	}

	return nil
}
