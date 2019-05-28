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
	"io/ioutil"
	"os"
)

// NewSimpleJSONParser returns a INI parser based on file with the priority 100,
// which registers the CLI option, cliOptName, into the default group and reads
// the data from the INI file appointed by cliOptName.
func NewSimpleJSONParser(cliOptName string) Parser {
	return NewJSONParser(100, func(c *Config) error {
		c.RegisterCliOpt(Str(cliOptName, "", "The path of the JSON config file."))
		return nil
	}, func(c *Config) ([]byte, error) {
		// Read the content of the config file.
		if filename := c.StringD(cliOptName, ""); filename == "" {
			return nil, nil
		} else if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
			c.Debugf("[json] Warning: the file named '%s' does not exist", filename)
			return nil, nil
		} else if data, err := ioutil.ReadFile(filename); err != nil {
			return nil, err
		} else {
			return data, nil
		}
	})
}

// NewJSONParser returns a new JSON parser.
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
func NewJSONParser(priority int, init func(*Config) error, getData func(*Config) ([]byte, error)) Parser {
	return &jsonParser{
		prio:    priority,
		init:    init,
		getData: getData,
	}
}

type jsonParser struct {
	prio    int
	init    func(*Config) error
	getData func(*Config) ([]byte, error)
}

func (j jsonParser) Name() string {
	return "json"
}

func (j jsonParser) Priority() int {
	return j.prio
}

func (j jsonParser) Pre(c *Config) error {
	if j.init != nil {
		return j.init(c)
	}
	return nil
}

func (j jsonParser) Parse(c *Config) error {
	return nil
}

func (j jsonParser) Post(c *Config) error {
	data, err := j.getData(c)
	if err != nil {
		return err
	} else if len(data) == 0 {
		return nil
	}

	var ms map[string]interface{}
	if err = json.Unmarshal(data, &ms); err != nil {
		return err
	}

	return j.update(c, c.OptGroup, ms)
}

func (j jsonParser) update(c *Config, group *OptGroup, ms map[string]interface{}) error {
	for key, value := range ms {
		if _ms, ok := value.(map[string]interface{}); ok {
			if subGroup := group.Group(key); subGroup != nil {
				if err := j.update(c, subGroup, _ms); err != nil {
					return err
				}
			}
			continue
		}

		switch err := group.UpdateOptValue(key, value); err {
		case nil, ErrNoOpt:
		default:
			return err
		}
	}
	return nil
}
