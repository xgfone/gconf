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
	"os"
	"strings"
	"unicode"
)

type propertyParser struct {
	parsed int32

	sep  string
	prio int

	init    func(*Config) error
	getData func(*Config) ([]byte, error)
}

// NewSimplePropertyParser returns a property parser with the priority 100,
// which registers the CLI option, cliOptName, into the default group and reads
// the data from the property file appointed by cliOptName.
func NewSimplePropertyParser(cliOptName string) Parser {
	return NewPropertyParser(100, func(c *Config) error {
		c.RegisterCliOpt(Str(cliOptName, "", "The path of the property config file."))
		return nil
	}, func(c *Config) ([]byte, error) {
		// Read the content of the config file.
		if filename := c.StringD(cliOptName, ""); filename == "" {
			return nil, nil
		} else if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
			c.Debugf("[property] Warning: the file named '%s' does not exist", filename)
			return nil, nil
		} else if data, err := ioutil.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return data, nil
		}
	})
}

// NewPropertyParser returns a new property parser based on the file.
//
// The first argument sets the Init function.
//
// The second argument sets the Init function to initialize the parser, such as
// registering the CLI option.
//
// The third argument is used to read the data to be parsed, which will
// be called at the start of calling the method Parse().
//
// The property parser supports the line comments starting with "#", "//" or ";".
// The key and the value is separated by an equal sign, that's =. The key must
// be in one of ., :, _, -, number and letter.
//
// If the value ends with "\", it will continue the next line. The lines will
// be joined by "\n" together.
func NewPropertyParser(priority int, init func(*Config) error, getData func(*Config) ([]byte, error)) Parser {
	return &propertyParser{
		sep:  "=",
		prio: priority,

		init:    init,
		getData: getData,
	}
}

func (p *propertyParser) Name() string {
	return "property"
}

func (p *propertyParser) Priority() int {
	return p.prio
}

func (p *propertyParser) Pre(c *Config) error {
	if p.init != nil {
		return p.init(c)
	}
	return nil
}

func (p *propertyParser) Post(c *Config) error {
	return nil
}

func (p *propertyParser) Parse(c *Config) error {
	data, err := p.getData(c)
	if err != nil {
		return err
	} else if len(data) == 0 {
		return nil
	}

	// Parse the config file.
	opts := make([][3]interface{}, 0)
	lines := strings.Split(string(data), "\n")
	for index, maxIndex := 0, len(lines); index < maxIndex; {
		line := strings.TrimSpace(lines[index])
		index++

		c.Debugf("[%s] Parsing %dth line: '%s'", p.Name(), index, line)

		// Ignore the empty line.
		if len(line) == 0 {
			continue
		}

		// Ignore the line comments starting with "#", ";" or "//".
		if (line[0] == '#') || (line[0] == ';') ||
			(len(line) > 1 && line[0] == '/' && line[1] == '/') {
			continue
		}

		ss := strings.SplitN(line, p.sep, 2)
		if len(ss) != 2 {
			return fmt.Errorf("the %dth line misses the separator '%s'", index, p.sep)
		}

		key := strings.TrimSpace(ss[0])
		if len(key) == 0 {
			return fmt.Errorf("empty identifier key")
		}
		for _, r := range key {
			if unicode.IsSpace(r) || !unicode.IsPrint(r) {
				return fmt.Errorf("invalid identifier key '%s'", key)
			}
		}

		value := strings.TrimSpace(ss[1])
		if value != "" {
			for index < maxIndex && value[len(value)-1] == '\\' {
				value = strings.TrimRight(value, "\\") + strings.TrimSpace(lines[index])
				index++
				c.Debugf("[%s] Parsing %dth line: '%s'", p.Name(), index, lines[index])
			}
		}

		var gname, optname string
		ss = strings.Split(key, c.GetGroupSeparator())
		switch _len := len(ss) - 1; _len {
		case 0:
			optname = key
		default:
			gname = strings.Join(ss[:_len], c.GetGroupSeparator())
			optname = ss[_len]
		}

		if group := c.Group(gname); group == nil {
			continue
		} else if opt := group.Opt(optname); opt == nil {
			continue
		} else if v, err := opt.Parse(value); err != nil {
			return err
		} else {
			opts = append(opts, [3]interface{}{group, optname, v})
		}
	}

	for _, opt := range opts {
		opt[0].(*OptGroup).UpdateOptValue(opt[1].(string), opt[2])
	}

	return nil
}
