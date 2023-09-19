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
	"sort"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
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
	RegisterOpts(opts...)

	registeredOpts := GetAllOpts()
	if len(registeredOpts) != len(opts) {
		t.Errorf("expect %d opts, but got %d", len(opts), len(registeredOpts))
		return
	}

	sort.SliceStable(opts, func(i, j int) bool { return opts[i].Name < opts[j].Name })
	sort.SliceStable(registeredOpts, func(i, j int) bool {
		return registeredOpts[i].Name < registeredOpts[j].Name
	})
	for i := 0; i < len(registeredOpts); i++ {
		if registeredOpts[i].Name != opts[i].Name {
			t.Errorf("expect the option '%s', but got '%s'", opts[i].Name, registeredOpts[i].Name)
		}
	}

	if v := GetBool("bool"); v {
		t.Errorf("bool option value expects '%v', but got '%v'", false, v)
	}
	if v := GetInt("int"); v != 0 {
		t.Errorf("int option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetInt32("int32"); v != 0 {
		t.Errorf("int32 option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetInt64("int64"); v != 0 {
		t.Errorf("int64 option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetUint("uint"); v != 0 {
		t.Errorf("uint option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetUint32("uint32"); v != 0 {
		t.Errorf("uint32 option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetUint64("uint64"); v != 0 {
		t.Errorf("uint64 option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetFloat64("float64"); v != 0 {
		t.Errorf("float64 option value expects '%v', but got '%v'", 0, v)
	}
	if v := GetString("string"); v != "" {
		t.Errorf("string option value expects '%v', but got '%v'", "", v)
	}
	if v := GetDuration("duration"); v != 0 {
		t.Errorf("duration option value expects '%v', but got '%v'", time.Duration(0), v)
	}
	if v := GetTime("time"); !v.IsZero() {
		t.Errorf("time option value expects '%v', but got '%v'", time.Time{}, v)
	}
	if v := GetIntSlice("ints"); len(v) != 0 {
		t.Errorf("ints option value expects '%v', but got '%v'", []int{}, v)
	}
	if v := GetUintSlice("uints"); len(v) != 0 {
		t.Errorf("uints option value expects '%v', but got '%v'", []uint{}, v)
	}
	if v := GetFloat64Slice("float64s"); len(v) != 0 {
		t.Errorf("float64s option value expects '%v', but got '%v'", []float64{}, v)
	}
	if v := GetStringSlice("strings"); len(v) != 0 {
		t.Errorf("strings option value expects '%v', but got '%v'", []string{}, v)
	}
	if v := GetDurationSlice("durations"); len(v) != 0 {
		t.Errorf("durations option value expects '%v', but got '%v'", []time.Duration{}, v)
	}

	_ = LoadMap(map[string]interface{}{
		"bool":      true,
		"string":    "abc",
		"int":       111,
		"int32":     222,
		"int64":     333,
		"uint":      444,
		"uint32":    555,
		"uint64":    666,
		"float64":   777,
		"duration":  "1s",
		"time":      1629877851,
		"ints":      "1, 2, 3",
		"uints":     "4,5 ,6",
		"float64s":  "7,8, 9",
		"strings":   []string{"a", "b", "c"},
		"durations": "1s, 2s, 3s",
	})

	if v := GetBool("bool"); v != true {
		t.Errorf("bool option value expects '%v', but got '%v'", true, v)
	}
	if v := GetInt("int"); v != 111 {
		t.Errorf("int option value expects '%v', but got '%v'", 111, v)
	}
	if v := GetInt32("int32"); v != 222 {
		t.Errorf("int32 option value expects '%v', but got '%v'", 222, v)
	}
	if v := GetInt64("int64"); v != 333 {
		t.Errorf("int64 option value expects '%v', but got '%v'", 333, v)
	}
	if v := GetUint("uint"); v != 444 {
		t.Errorf("uint option value expects '%v', but got '%v'", 444, v)
	}
	if v := GetUint32("uint32"); v != 555 {
		t.Errorf("uint32 option value expects '%v', but got '%v'", 555, v)
	}
	if v := GetUint64("uint64"); v != 666 {
		t.Errorf("uint64 option value expects '%v', but got '%v'", 666, v)
	}
	if v := GetFloat64("float64"); v != 777 {
		t.Errorf("float64 option value expects '%v', but got '%v'", 777, v)
	}
	if v := GetString("string"); v != "abc" {
		t.Errorf("string option value expects '%v', but got '%v'", "abc", v)
	}
	if v := GetDuration("duration"); v != time.Second {
		t.Errorf("duration option value expects '%v', but got '%v'", time.Second, v)
	}
	if v := GetTime("time"); !v.Equal(time.Unix(1629877851, 0)) {
		t.Errorf("time option value expects '%v', but got '%v'", time.Unix(1629877851, 0), v)
	}
	if v := GetIntSlice("ints"); !reflect.DeepEqual(v, []int{1, 2, 3}) {
		t.Errorf("ints option value expects '%v', but got '%v'", []int{1, 2, 3}, v)
	}
	if v := GetUintSlice("uints"); !reflect.DeepEqual(v, []uint{4, 5, 6}) {
		t.Errorf("uints option value expects '%v', but got '%v'", []uint{4, 5, 6}, v)
	}
	if v := GetFloat64Slice("float64s"); !reflect.DeepEqual(v, []float64{7, 8, 9}) {
		t.Errorf("float64s option value expects '%v', but got '%v'", []float64{7, 8, 9}, v)
	}
	if v := GetStringSlice("strings"); !reflect.DeepEqual(v, []string{"a", "b", "c"}) {
		t.Errorf("strings option value expects '%v', but got '%v'", []string{"a", "b", "c"}, v)
	}

	ds := []time.Duration{time.Second, time.Second * 2, time.Second * 3}
	if v := GetDurationSlice("durations"); !reflect.DeepEqual(v, ds) {
		t.Errorf("durations option value expects '%v', but got '%v'", ds, v)
	}

	names := make([]string, len(opts))
	for i := 0; i < len(opts); i++ {
		names[i] = opts[i].Name
	}
	UnregisterOpts(names...)

	registeredOpts = GetAllOpts()
	if len(registeredOpts) != 0 {
		t.Errorf("unexpected options: %v", registeredOpts)
	}
}

func TestConfig_Observe(t *testing.T) {
	Conf.reset()
	type option struct {
		name string
		old  interface{}
		new  interface{}
	}

	var options []option
	Observe(func(name string, old, new interface{}) {
		options = append(options, option{name: name, old: old, new: new})
	})

	opts := []Opt{
		StrOpt("str", "str"),
		Int32Opt("int", "int"),
	}
	RegisterOpts(opts...)

	if err := Set("str", "abc"); err != nil {
		t.Errorf("unknown error: %s", err)
	}
	if err := Set("int", 123); err != nil {
		t.Errorf("unknown error: %s", err)
	}

	if len(options) != 2 {
		t.Errorf("expect %d changed opts, but got %d", 2, len(options))
	} else if o := options[0]; o.name != "str" || o.old.(string) != "" || o.new.(string) != "abc" {
		t.Errorf("unexpected changed opt: %+v", o)
	} else if o := options[1]; o.name != "int" || o.old.(int32) != 0 || o.new.(int32) != 123 {
		t.Errorf("unexpected changed opt: %+v", o)
	}
}
