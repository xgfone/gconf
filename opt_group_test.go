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
	"reflect"
	"testing"
	"time"
)

func TestOptGroup(t *testing.T) {
	Conf.reset()

	opts := []Opt{
		BoolOpt("bool", "bool"),
		StrOpt("string", "string"),
		IntOpt("int", "int"),
		Int32Opt("int32", "int32"),
		Int64Opt("int64", "int64"),
		UintOpt("uint", "uint"),
		Uint32Opt("uint32", "uint32"),
		Uint64Opt("uint64", "uint64"),
		Float64Opt("float64", "float64"),
		DurationOpt("duration", "duration"),
		TimeOpt("time", "time"),
		IntSliceOpt("ints", "ints"),
		UintSliceOpt("uints", "uints"),
		Float64SliceOpt("float64s", "float64s"),
		StrSliceOpt("strings", "strings"),
		DurationSliceOpt("durations", "durations"),
	}

	g1 := Group("group1")
	g1.RegisterOpts(opts...)

	g2 := g1.Group("group2")
	g2.RegisterOpts(StrOpt("string", "string"))

	registeredOpts := GetAllOpts()
	if expect := len(opts) + 1; expect != len(registeredOpts) {
		t.Errorf("expect %d opts, but got %d", expect, len(registeredOpts))
		return
	}

	for _, opt := range registeredOpts {
		switch opt.Name {
		case "group1.bool",
			"group1.string",
			"group1.int",
			"group1.int32",
			"group1.int64",
			"group1.uint",
			"group1.uint32",
			"group1.uint64",
			"group1.float64",
			"group1.duration",
			"group1.time",
			"group1.ints",
			"group1.uints",
			"group1.float64s",
			"group1.strings",
			"group1.durations",
			"group1.group2.string":
		default:
			t.Errorf("unexpected option name '%s'", opt.Name)
		}
	}

	g1.Set("bool", "true")
	g1.Set("string", "abc")
	g1.Set("int", 111)
	g1.Set("int32", 222)
	g1.Set("int64", 333)
	g1.Set("uint", "444")
	g1.Set("uint32", "555")
	g1.Set("uint64", "666")
	g1.Set("float64", 777)
	g1.Set("duration", "1s")
	g1.Set("time", "2021-08-25T21:50:51+08:00")
	g1.Set("ints", "1,2,3")
	g1.Set("uints", "4,5,6")
	g1.Set("float64s", []float64{7, 8, 9})
	g1.Set("strings", []string{"a", "b", "c"})
	g1.Set("durations", []time.Duration{time.Second})
	g2.Set("string", "xyz")

	if v := g1.GetBool("bool"); v != true {
		t.Errorf("bool option value expect '%v', but got '%v'", true, v)
	}
	if v := g1.GetString("string"); v != "abc" {
		t.Errorf("string option value expect '%v', but got '%v'", "abc", v)
	}
	if v := g1.GetInt("int"); v != 111 {
		t.Errorf("int option value expect '%v', but got '%v'", 111, v)
	}
	if v := g1.GetInt32("int32"); v != 222 {
		t.Errorf("int32 option value expect '%v', but got '%v'", 222, v)
	}
	if v := g1.GetInt64("int64"); v != 333 {
		t.Errorf("int64 option value expect '%v', but got '%v'", 333, v)
	}
	if v := g1.GetUint("uint"); v != 444 {
		t.Errorf("uint option value expect '%v', but got '%v'", 444, v)
	}
	if v := g1.GetUint32("uint32"); v != 555 {
		t.Errorf("uint32 option value expect '%v', but got '%v'", 555, v)
	}
	if v := g1.GetUint64("uint64"); v != 666 {
		t.Errorf("uint64 option value expect '%v', but got '%v'", 666, v)
	}
	if v := g1.GetFloat64("float64"); v != 777 {
		t.Errorf("float64 option value expect '%v', but got '%v'", 777, v)
	}
	if v := g1.GetDuration("duration"); v != time.Second {
		t.Errorf("duration option value expect '%v', but got '%v'", time.Second, v)
	}
	_time, _ := time.Parse(time.RFC3339, "2021-08-25T21:50:51+08:00")
	if v := g1.GetTime("time"); !v.Equal(_time) {
		t.Errorf("time option value expect '%v', but got '%v'", _time, v)
	}
	if v := g1.GetIntSlice("ints"); !reflect.DeepEqual(v, []int{1, 2, 3}) {
		t.Errorf("ints option value expect '%v', but got '%v'", []int{1, 2, 3}, v)
	}
	if v := g1.Get("uints"); !reflect.DeepEqual(v, []uint{4, 5, 6}) {
		t.Errorf("uints option value expect '%v', but got '%v'", []uint{4, 5, 6}, v)
	}
	if v := g1.Get("float64s"); !reflect.DeepEqual(v, []float64{7, 8, 9}) {
		t.Errorf("float64s option value expect '%v', but got '%v'", []float64{7, 8, 9}, v)
	}
	if v := g1.Get("strings"); !reflect.DeepEqual(v, []string{"a", "b", "c"}) {
		t.Errorf("strings option value expect '%v', but got '%v'", []string{"a", "b", "c"}, v)
	}
	if v := g1.Get("durations"); !reflect.DeepEqual(v, []time.Duration{time.Second}) {
		t.Errorf("durations option value expect '%v', but got '%v'", []time.Duration{time.Second}, v)
	}
	if v := g2.Get("string"); v != "xyz" {
		t.Errorf("string option value expect '%v', but got '%v'", "xyz", v)
	}

	g2.UnregisterOpts("string")
	for _, opt := range opts {
		g1.UnregisterOpts(opt.Name)
	}

	registeredOpts = GetAllOpts()
	if len(registeredOpts) != 0 {
		t.Errorf("unexpected options: %v", registeredOpts)
	}
}

func TestOptGroupEmptyName(t *testing.T) {
	config := New()
	group1 := config.Group("")
	group1.RegisterOpts(StrOpt("opt1", "help"))

	group2 := group1.Group("")
	group2.RegisterOpts(StrOpt("opt2", "help"))

	group3 := group1.Group("group1")
	group3.RegisterOpts(StrOpt("opt3", "help"))
	group3.Self("help")

	group4 := group3.Group("")
	group4.RegisterOpts(StrOpt("opt4", "help"))

	group5 := group3.Group("group2")
	group5.RegisterOpts(StrOpt("opt5", "help"))

	opts := config.GetAllOpts()
	if len(opts) != 6 {
		t.Error(opts)
	}

	for _, opt := range opts {
		switch opt.Name {
		case
			"opt1",
			"opt2",
			"group1",
			"group1.opt3",
			"group1.opt4",
			"group1.group2.opt5":
		default:
			t.Errorf("unexpected opt '%s'", opt.Name)
		}
	}
}
