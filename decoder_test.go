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
	"testing"
)

var yamlData = `
opt1: xyz
group:
    opt2: 456
`

func TestNewYamlDecoder(t *testing.T) {
	conf := New()
	conf.RegisterOpt(StrOpt("opt1", "").D("abc"))
	conf.NewGroup("group").RegisterOpt(IntOpt("opt2", "").D(123))

	ms := make(map[string]interface{})
	if err := NewYamlDecoder().Decode([]byte(yamlData), ms); err != nil {
		t.Error(err)
	} else {
		conf.LoadMap(ms)
		if v := conf.GetString("opt1"); v != "xyz" {
			t.Error(v)
		} else if v := conf.Group("group").GetInt("opt2"); v != 456 {
			t.Error(v)
		}
	}
}

var tomlData = `
opt1 = "xyz"
[group]
opt2 = 456
`

func TestNewTomlDecoder(t *testing.T) {
	conf := New()
	conf.RegisterOpt(StrOpt("opt1", "").D("abc"))
	conf.NewGroup("group").RegisterOpt(IntOpt("opt2", "").D(123))

	ms := make(map[string]interface{})
	if err := NewTomlDecoder().Decode([]byte(tomlData), ms); err != nil {
		t.Error(err)
	} else {
		conf.LoadMap(ms)
		if v := conf.GetString("opt1"); v != "xyz" {
			t.Error(v)
		} else if v := conf.Group("group").GetInt("opt2"); v != 456 {
			t.Error(v)
		}
	}
}
