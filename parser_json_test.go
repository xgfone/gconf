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

import "testing"

func TestNewJSONParser(t *testing.T) {
	type Command1OptGroup struct {
		Opt1 string
		Opt2 int
	}

	type Command1 struct {
		Opt1 int
		Opt2 string

		Group Command1OptGroup
	}

	type Command2OptGroup struct {
		Opt1 int
		Opt2 string
	}

	type Command2 struct {
		Opt1 string
		Opt2 int

		Group Command2OptGroup
	}

	type OptGroup struct {
		Opt1 string
		Opt2 int
	}

	type config struct {
		Opt1 int
		Opt2 string

		Group OptGroup `cli:"false"`

		Cmd1 Command1 `cmd:"cmd1"`
		Cmd2 Command2 `cmd:"cmd2"`
	}

	parser := NewJSONParser(100, nil, func(c *Config) ([]byte, error) {
		data := []byte(`{
			"opt1": 123,
			"opt2": "abc",
			"opt3": 456,
			"opt4": "xyz",
			"group1": {"key": "value"},
			"group": {"opt1": "efg", "opt2": 789},
			"cmd1": {"opt1": 234, "opt2": "hij", "group": {"opt1": "lmn", "opt2": 345}},
			"cmd2": {"opt1": "opq", "opt2": 567, "group": {"opt1": 890, "opt2": "rst"}}
		}`)
		return data, nil
	})

	var c config
	conf := New().AddParser(parser)
	conf.RegisterStruct(&c)
	if err := conf.Parse(); err != nil {
		t.Error(err)
	} else if c.Opt1 != 123 || c.Opt2 != "abc" {
		t.Error(c)
	} else if c.Group.Opt1 != "efg" || c.Group.Opt2 != 789 {
		t.Error(c)
	} else if c.Cmd1.Opt1 != 234 || c.Cmd1.Opt2 != "hij" || c.Cmd1.Group.Opt1 != "lmn" || c.Cmd1.Group.Opt2 != 345 {
		t.Error(c)
	} else if c.Cmd2.Opt1 != "opq" || c.Cmd2.Opt2 != 567 || c.Cmd2.Group.Opt1 != 890 || c.Cmd2.Group.Opt2 != "rst" {
		t.Error(c)
	}
}
