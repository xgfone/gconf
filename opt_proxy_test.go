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
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestOptProxy(t *testing.T) {
	Conf.reset()

	boolopt := NewBool("bool", false, "bool")
	intopt := NewInt("int", 1, "int")
	int32opt := NewInt32("int32", 2, "int32")
	int64opt := NewInt64("int64", 3, "int64")
	uintopt := NewUint("uint", 4, "uint")
	uint32opt := NewUint32("uint32", 5, "uint32")
	uint64opt := NewUint64("uint64", 6, "uint64")
	float64opt := NewFloat64("float64", 7, "float64")
	stringopt := NewString("string", "a", "string")
	durationopt := NewDuration("duration", time.Second, "duration")
	timeopt := NewTime("time", time.Unix(1629877851, 0), "time")
	intsopt := NewIntSlice("ints", []int{1}, "ints")
	uintsopt := NewUintSlice("uints", []uint{4}, "uints")
	float64sopt := NewFloat64Slice("float64s", []float64{7}, "float64s")
	stringsopt := NewStringSlice("strings", []string{"a"}, "strings")
	durationsopt := NewDurationSlice("durations", []time.Duration{time.Second}, "durations")

	if v := boolopt.Get(); v != false {
		t.Errorf("option value expect '%v', but got '%v'", false, v)
	}
	if v := intopt.Get(); v != 1 {
		t.Errorf("option value expect '%v', but got '%v'", 1, v)
	}
	if v := int32opt.Get(); v != 2 {
		t.Errorf("option value expect '%v', but got '%v'", 2, v)
	}
	if v := int64opt.Get(); v != 3 {
		t.Errorf("option value expect '%v', but got '%v'", 3, v)
	}
	if v := uintopt.Get(); v != 4 {
		t.Errorf("option value expect '%v', but got '%v'", 4, v)
	}
	if v := uint32opt.Get(); v != 5 {
		t.Errorf("option value expect '%v', but got '%v'", 5, v)
	}
	if v := uint64opt.Get(); v != 6 {
		t.Errorf("option value expect '%v', but got '%v'", 6, v)
	}
	if v := float64opt.Get(); v != 7 {
		t.Errorf("option value expect '%v', but got '%v'", 7, v)
	}
	if v := stringopt.Get(); v != "a" {
		t.Errorf("option value expect '%v', but got '%v'", "a", v)
	}
	if v := durationopt.Get(); v != time.Second {
		t.Errorf("option value expect '%v', but got '%v'", time.Second, v)
	}
	if v := timeopt.Get(); v != time.Unix(1629877851, 0) {
		t.Errorf("option value expect '%v', but got '%v'", time.Unix(1629877851, 0), v)
	}
	if v := intsopt.Get(); !reflect.DeepEqual(v, []int{1}) {
		t.Errorf("option value expect '%v', but got '%v'", []int{1}, v)
	}
	if v := uintsopt.Get(); !reflect.DeepEqual(v, []uint{4}) {
		t.Errorf("option value expect '%v', but got '%v'", []uint{4}, v)
	}
	if v := float64sopt.Get(); !reflect.DeepEqual(v, []float64{7}) {
		t.Errorf("option value expect '%v', but got '%v'", []float64{7}, v)
	}
	if v := stringsopt.Get(); !reflect.DeepEqual(v, []string{"a"}) {
		t.Errorf("option value expect '%v', but got '%v'", []string{"a"}, v)
	}
	if v := durationsopt.Get(); !reflect.DeepEqual(v, []time.Duration{time.Second}) {
		t.Errorf("option value expect '%v', but got '%v'", []time.Duration{time.Second}, v)
	}

	boolopt.Set(true)
	intopt.Set(11)
	int32opt.Set(22)
	int64opt.Set(33)
	uintopt.Set(44)
	uint32opt.Set(55)
	uint64opt.Set(66)
	float64opt.Set(77)
	stringopt.Set("abc")
	durationopt.Set(time.Hour)
	timeopt.Set(time.Unix(1629878851, 0))
	intsopt.Set("1,2,3")
	uintsopt.Set("4,5,6")
	float64sopt.Set("7,8,9")
	stringsopt.Set("x,y,z")
	durationsopt.Set("1s,1h")

	if v := boolopt.Get(); v != true {
		t.Errorf("option value expect '%v', but got '%v'", true, v)
	}
	if v := intopt.Get(); v != 11 {
		t.Errorf("option value expect '%v', but got '%v'", 11, v)
	}
	if v := int32opt.Get(); v != 22 {
		t.Errorf("option value expect '%v', but got '%v'", 22, v)
	}
	if v := int64opt.Get(); v != 33 {
		t.Errorf("option value expect '%v', but got '%v'", 33, v)
	}
	if v := uintopt.Get(); v != 44 {
		t.Errorf("option value expect '%v', but got '%v'", 44, v)
	}
	if v := uint32opt.Get(); v != 55 {
		t.Errorf("option value expect '%v', but got '%v'", 55, v)
	}
	if v := uint64opt.Get(); v != 66 {
		t.Errorf("option value expect '%v', but got '%v'", 66, v)
	}
	if v := float64opt.Get(); v != 77 {
		t.Errorf("option value expect '%v', but got '%v'", 77, v)
	}
	if v := stringopt.Get(); v != "abc" {
		t.Errorf("option value expect '%v', but got '%v'", "abc", v)
	}
	if v := durationopt.Get(); v != time.Hour {
		t.Errorf("option value expect '%v', but got '%v'", time.Hour, v)
	}
	if v := timeopt.Get(); v != time.Unix(1629878851, 0) {
		t.Errorf("option value expect '%v', but got '%v'", time.Unix(1629878851, 0), v)
	}
	if v := intsopt.Get(); !reflect.DeepEqual(v, []int{1, 2, 3}) {
		t.Errorf("option value expect '%v', but got '%v'", []int{1, 2, 3}, v)
	}
	if v := uintsopt.Get(); !reflect.DeepEqual(v, []uint{4, 5, 6}) {
		t.Errorf("option value expect '%v', but got '%v'", []uint{4, 5, 6}, v)
	}
	if v := float64sopt.Get(); !reflect.DeepEqual(v, []float64{7, 8, 9}) {
		t.Errorf("option value expect '%v', but got '%v'", []float64{7, 8, 9}, v)
	}
	if v := stringsopt.Get(); !reflect.DeepEqual(v, []string{"x", "y", "z"}) {
		t.Errorf("option value expect '%v', but got '%v'", []string{"x", "y", "z"}, v)
	}
	if v := durationsopt.Get(); !reflect.DeepEqual(v, []time.Duration{time.Second, time.Hour}) {
		t.Errorf("option value expect '%v', but got '%v'", []time.Duration{time.Second, time.Hour}, v)
	}
}

func TestOptGroupProxy(t *testing.T) {
	config := New()
	group := config.Group("group1.group2")
	boolopt := group.NewBool("bool", false, "bool")
	intopt := group.NewInt("int", 1, "int")
	int32opt := group.NewInt32("int32", 2, "int32")
	int64opt := group.NewInt64("int64", 3, "int64")
	uintopt := group.NewUint("uint", 4, "uint")
	uint32opt := group.NewUint32("uint32", 5, "uint32")
	uint64opt := group.NewUint64("uint64", 6, "uint64")
	float64opt := group.NewFloat64("float64", 7, "float64")
	stringopt := group.NewString("string", "a", "string")
	durationopt := group.NewDuration("duration", time.Second, "duration")
	timeopt := group.NewTime("time", time.Unix(1629877851, 0), "time")
	intsopt := group.NewIntSlice("ints", []int{1}, "ints")
	uintsopt := group.NewUintSlice("uints", []uint{4}, "uints")
	float64sopt := group.NewFloat64Slice("float64s", []float64{7}, "float64s")
	stringsopt := group.NewStringSlice("strings", []string{"a"}, "strings")
	durationsopt := group.NewDurationSlice("durations", []time.Duration{time.Second}, "durations")

	boolopt.Set(true)
	intopt.Set(11)
	int32opt.Set(22)
	int64opt.Set(33)
	uintopt.Set(44)
	uint32opt.Set(55)
	uint64opt.Set(66)
	float64opt.Set(77)
	stringopt.Set("abc")
	durationopt.Set(time.Hour)
	timeopt.Set(time.Unix(1629878851, 0))
	intsopt.Set("1,2,3")
	uintsopt.Set("4,5,6")
	float64sopt.Set("7,8,9")
	stringsopt.Set("x,y,z")
	durationsopt.Set("1s,1h")

	if v := config.Get(intopt.Name()); v != 11 {
		t.Errorf("option value expect '%v', but got '%v'", 11, v)
	}

	if v := boolopt.Get(); v != true {
		t.Errorf("option value expect '%v', but got '%v'", true, v)
	}
	if v := intopt.Get(); v != 11 {
		t.Errorf("option value expect '%v', but got '%v'", 11, v)
	}
	if v := int32opt.Get(); v != 22 {
		t.Errorf("option value expect '%v', but got '%v'", 22, v)
	}
	if v := int64opt.Get(); v != 33 {
		t.Errorf("option value expect '%v', but got '%v'", 33, v)
	}
	if v := uintopt.Get(); v != 44 {
		t.Errorf("option value expect '%v', but got '%v'", 44, v)
	}
	if v := uint32opt.Get(); v != 55 {
		t.Errorf("option value expect '%v', but got '%v'", 55, v)
	}
	if v := uint64opt.Get(); v != 66 {
		t.Errorf("option value expect '%v', but got '%v'", 66, v)
	}
	if v := float64opt.Get(); v != 77 {
		t.Errorf("option value expect '%v', but got '%v'", 77, v)
	}
	if v := stringopt.Get(); v != "abc" {
		t.Errorf("option value expect '%v', but got '%v'", "abc", v)
	}
	if v := durationopt.Get(); v != time.Hour {
		t.Errorf("option value expect '%v', but got '%v'", time.Hour, v)
	}
	if v := timeopt.Get(); v != time.Unix(1629878851, 0) {
		t.Errorf("option value expect '%v', but got '%v'", time.Unix(1629878851, 0), v)
	}
	if v := intsopt.Get(); !reflect.DeepEqual(v, []int{1, 2, 3}) {
		t.Errorf("option value expect '%v', but got '%v'", []int{1, 2, 3}, v)
	}
	if v := uintsopt.Get(); !reflect.DeepEqual(v, []uint{4, 5, 6}) {
		t.Errorf("option value expect '%v', but got '%v'", []uint{4, 5, 6}, v)
	}
	if v := float64sopt.Get(); !reflect.DeepEqual(v, []float64{7, 8, 9}) {
		t.Errorf("option value expect '%v', but got '%v'", []float64{7, 8, 9}, v)
	}
	if v := stringsopt.Get(); !reflect.DeepEqual(v, []string{"x", "y", "z"}) {
		t.Errorf("option value expect '%v', but got '%v'", []string{"x", "y", "z"}, v)
	}
	if v := durationsopt.Get(); !reflect.DeepEqual(v, []time.Duration{time.Second, time.Hour}) {
		t.Errorf("option value expect '%v', but got '%v'", []time.Duration{time.Second, time.Hour}, v)
	}

	opts := config.GetAllOpts()
	for _, opt := range opts {
		switch opt.Name {
		case
			"group1.group2.bool",
			"group1.group2.int",
			"group1.group2.int32",
			"group1.group2.int64",
			"group1.group2.uint",
			"group1.group2.uint32",
			"group1.group2.uint64",
			"group1.group2.float64",
			"group1.group2.string",
			"group1.group2.duration",
			"group1.group2.time",
			"group1.group2.ints",
			"group1.group2.uints",
			"group1.group2.float64s",
			"group1.group2.strings",
			"group1.group2.durations":
		default:
			t.Errorf("unexpected option '%s'", opt.Name)
		}
	}
}

func TestOptProxyOnUpdate(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	config := New()
	opt := config.NewString("opt", "abc", "help")
	opt.OnUpdate(func(old, new interface{}) {
		fmt.Fprintf(buf, "%s: %v -> %v", opt.Name(), old, new)
	})

	opt.Set("xyz")
	expect := `opt: abc -> xyz`
	if s := buf.String(); s != expect {
		t.Errorf("expect '%s', but got '%s'", expect, s)
	}
}
