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

import "testing"

func TestConfig_Snapshot(t *testing.T) {
	config := New()
	config.RegisterOpts(
		StrOpt("opt1", ""),
		IntOpt("opt2", ""),
	)

	config.Set("opt1", "a")
	config.Set("opt2", 1)
	gen, snaps := config.Snapshot()
	if gen != 2 {
		t.Errorf("expect %d generation, but got %d", 2, gen)
	} else if len(snaps) != 2 {
		t.Errorf("expect %d snapshot elements, but got %d", 2, len(snaps))
	} else {
		for name, value := range snaps {
			switch name {
			case "opt1":
				if value.(string) != "a" {
					t.Errorf("expect the value '%s', but got '%v'", "a", value)
				}
			case "opt2":
				if value.(int) != 1 {
					t.Errorf("expect the value '%d', but got '%v'", 1, value)
				}
			default:
				t.Errorf("unexpected the option '%s'", name)
			}
		}
	}
	config.Set("opt1", "b")
	config.Set("opt2", 2)
	gen, snaps = config.Snapshot()
	if gen != 4 {
		t.Errorf("expect %d generation, but got %d", 4, gen)
	} else if len(snaps) != 2 {
		t.Errorf("expect %d snapshot elements, but got %d", 2, len(snaps))
	} else {
		for name, value := range snaps {
			switch name {
			case "opt1":
				if value.(string) != "b" {
					t.Errorf("expect the value '%s', but got '%v'", "b", value)
				}
			case "opt2":
				if value.(int) != 2 {
					t.Errorf("expect the value '%d', but got '%v'", 2, value)
				}
			default:
				t.Errorf("unexpected the option '%s'", name)
			}
		}
	}

}
