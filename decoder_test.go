// Copyright 2021 xgfone
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
	"testing"
)

func ExampleNewJSONDecoder() {
	data := []byte(`{
		// user name
		"name": "Aaron",
		"age": 123,

		// the other information
		"other": {
			// address
			"home": "http://www.example.com"
		}
	}`)

	ms := make(map[string]interface{})
	err := NewJSONDecoder()(data, ms)

	fmt.Println(err)
	fmt.Println(len(ms))
	fmt.Println(ms["name"])
	fmt.Println(ms["age"])
	fmt.Println(ms["other"])

	// Output:
	// <nil>
	// 3
	// Aaron
	// 123
	// map[home:http://www.example.com]
}

var yamlData = `
opt1: xyz
group1:
    opt2: 456
    group2:
        opt3: 789
`

func TestNewYamlDecoder(t *testing.T) {
	conf := New()
	conf.RegisterOpts(StrOpt("opt1", "").D("abc"))
	conf.Group("group1").RegisterOpts(IntOpt("opt2", "").D(123))
	conf.Group("group1.group2").RegisterOpts(IntOpt("opt3", ""))

	ms := make(map[string]interface{})
	if err := NewYamlDecoder()([]byte(yamlData), ms); err != nil {
		t.Error(err)
	} else if err = conf.LoadMap(ms); err != nil {
		t.Error(err)
	} else {
		if v := conf.GetString("opt1"); v != "xyz" {
			t.Error(v)
		} else if v := conf.Group("group1").GetInt("opt2"); v != 456 {
			t.Error(v)
		} else if v := conf.GetInt("group1.group2.opt3"); v != 789 {
			t.Error(v)
		}
	}
}
