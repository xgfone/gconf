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
	"fmt"
	"testing"
)

func ExampleOpt_F() {
	fix := func(v interface{}) (interface{}, error) { return v.(int) + 1, nil }
	opt1 := IntOpt("opt1", "test fix with default").D(10).F(fix, true)
	opt2 := IntOpt("opt2", "test fix without default").D(20).F(fix)

	conf := New()
	conf.RegisterOpts(opt1, opt2)

	fmt.Printf("opt1=%s\n", conf.MustString("opt1"))
	fmt.Printf("opt2=%s\n", conf.MustString("opt2"))

	conf.UpdateValue("opt1", 30)
	conf.UpdateValue("opt2", 40)

	fmt.Printf("opt1=%s\n", conf.MustString("opt1"))
	fmt.Printf("opt2=%s\n", conf.MustString("opt2"))

	// Output:
	// opt1=11
	// opt2=20
	// opt1=31
	// opt2=41
}

func TestOptObserver(t *testing.T) {
	var value string
	opt := StrOpt("opt", "").D("abc").O(func(v interface{}) { value = v.(string) })

	conf := New()
	conf.RegisterOpts(opt)
	conf.UpdateOptValue("", "opt", "xyz")

	if value != "xyz" {
		t.Error(value, conf.MustString("opt"))
	}
}

func ExampleConfig_Traverse() {
	conf := New()
	conf.RegisterOpts(StrOpt("opt1", "").D("abc"))
	conf.NewGroup("group1").RegisterOpts(IntOpt("opt2", "").D(123))
	conf.NewGroup("group1").NewGroup("group2").RegisterOpts(IntOpt("opt3", "").D(456))

	conf.Traverse(func(group, opt string, value interface{}) {
		fmt.Printf("group=%s, opt=%s, value=%v\n", group, opt, value)
	})

	// Output:
	// group=, opt=opt1, value=abc
	// group=group1, opt=opt2, value=123
	// group=group1.group2, opt=opt3, value=456
}

func ExampleOptGroup_Migrate() {
	conf := New()
	conf.RegisterOpts(StrOpt("opt1", "").D("abc"))
	conf.RegisterOpts(StrOpt("opt2", "").D("efg"))
	conf.NewGroup("group1").RegisterOpts(StrOpt("opt3", "").D("opq"))
	conf.NewGroup("group1").RegisterOpts(StrOpt("opt4", "").D("rst"))

	conf.Migrate("opt1", "group1.opt3")
	conf.Group("group1").Migrate("opt4", "opt2")

	fmt.Printf("--- Before Updating ---\n")
	fmt.Printf("opt1=%v\n", conf.MustString("opt1"))
	fmt.Printf("opt2=%v\n", conf.MustString("opt2"))
	fmt.Printf("group1.opt3=%v\n", conf.G("group1").MustString("opt3"))
	fmt.Printf("group1.opt4=%v\n", conf.G("group1").MustString("opt4"))

	conf.UpdateValue("opt1", "uvw")
	conf.UpdateValue("group1.opt4", "xyz")

	fmt.Printf("--- After Updating ---\n")
	fmt.Printf("opt1=%v\n", conf.MustString("opt1"))
	fmt.Printf("opt2=%v\n", conf.MustString("opt2"))
	fmt.Printf("group1.opt3=%v\n", conf.G("group1").MustString("opt3"))
	fmt.Printf("group1.opt4=%v\n", conf.G("group1").MustString("opt4"))

	// Output:
	// --- Before Updating ---
	// opt1=abc
	// opt2=efg
	// group1.opt3=opq
	// group1.opt4=rst
	// --- After Updating ---
	// opt1=uvw
	// opt2=xyz
	// group1.opt3=uvw
	// group1.opt4=xyz
}

func ExampleConfig_Observe() {
	conf := New()
	conf.RegisterOpts(StrOpt("opt1", "").D("abc"))
	conf.NewGroup("group").RegisterOpts(IntOpt("opt2", "").D(123))
	conf.Group("group").NewGroup("subgroup").RegisterOpts(IntOpt("opt3", ""))
	conf.Observe(func(group, opt string, old, new interface{}) {
		fmt.Printf("Setting: group=%s, opt=%s, old=%v, new=%v\n", group, opt, old, new)
	})

	conf.Set("opt1", "xyz")
	conf.Group("group").Set("opt2", 789)
	conf.Group("group").Group("subgroup").Set("opt3", 456)

	// Output:
	// Setting: group=, opt=opt1, old=abc, new=xyz
	// Setting: group=group, opt=opt2, old=123, new=789
	// Setting: group=group.subgroup, opt=opt3, old=0, new=456
}

func ExampleOptGroup_FreezeOpt() {
	conf := New()
	conf.NewGroup("group1").RegisterOpts(StrOpt("opt1", "").D("a"), StrOpt("opt2", "").D("b"))
	conf.NewGroup("group2").RegisterOpts(StrOpt("opt3", "").D("c"), StrOpt("opt4", "").D("d"))
	conf.Group("group1").FreezeOpt("opt2")
	conf.Group("group2").FreezeGroup()

	conf.UpdateValue("group1.opt1", "o")
	conf.UpdateOptValue("group1", "opt2", "p")
	conf.UpdateOptValue("group2", "opt3", "q")
	conf.UpdateOptValue("group2", "opt4", "r")

	fmt.Println(conf.Group("group1").GetString("opt1"))
	fmt.Println(conf.Group("group1").GetString("opt2"))
	fmt.Println(conf.Group("group2").GetString("opt3"))
	fmt.Println(conf.Group("group2").GetString("opt4"))

	// Output:
	// o
	// b
	// c
	// d
}

func ExampleConfig_Snapshot() {
	conf := New()
	conf.RegisterOpts(StrOpt("opt1", ""))
	conf.NewGroup("group1").RegisterOpts(IntOpt("opt2", ""))
	conf.NewGroup("group1").NewGroup("group2").RegisterOpts(IntOpt("opt3", ""))

	// For test
	print := func(snap map[string]interface{}) {
		data, _ := json.Marshal(conf.Snapshot())
		fmt.Println(string(data))
	}

	print(conf.Snapshot())

	conf.Set("opt1", "abc")
	print(conf.Snapshot())

	conf.Group("group1").Set("opt2", 123)
	print(conf.Snapshot())

	conf.Group("group1.group2").Set("opt3", 456)
	print(conf.Snapshot())

	// Output:
	// {}
	// {"opt1":"abc"}
	// {"group1.opt2":123,"opt1":"abc"}
	// {"group1.group2.opt3":456,"group1.opt2":123,"opt1":"abc"}
}

func ExampleConfig() {
	opts := []Opt{
		BoolOpt("bool", "test bool opt"),
		StrOpt("string", "test string opt"),
		IntOpt("int", "test int opt"),
		Int32Opt("int32", "test int32 opt"),
		Int64Opt("int64", "test int64 opt"),
		UintOpt("uint", "test uint opt"),
		Uint32Opt("uint32", "test uint32 opt"),
		Uint64Opt("uint64", "test uint64 opt"),
		Float64Opt("float64", "test float64 opt"),
		DurationOpt("duration", "test time.Duration opt"),
		TimeOpt("time", "test time.Time opt"),

		// Slice
		IntSliceOpt("ints", "test []int opt"),
		UintSliceOpt("uints", "test []uint opt"),
		Float64SliceOpt("float64s", "test []float64 opt"),
		StrSliceOpt("strings", "test []string opt"),
		DurationSliceOpt("durations", "test []time.Duration opt"),
	}

	conf := New()
	conf.RegisterOpts(opts...)

	group1 := conf.NewGroup("group1")
	group1.RegisterOpts(opts...)

	group2 := group1.NewGroup("group2") // Or conf.NewGroup("group1.group2")
	group2.RegisterOpts(opts...)

	conf.Set("bool", "1")
	conf.Set("string", "abc")
	conf.Set("int", "123")
	conf.Set("int32", "123")
	conf.Set("int64", "123")
	conf.Set("uint", "123")
	conf.Set("uint32", "123")
	conf.Set("uint64", "123")
	conf.Set("float64", "123")
	conf.Set("duration", "123s")
	conf.Set("time", "2019-06-10T18:00:00Z")
	conf.Set("ints", "1,2,3")
	conf.Set("uints", "1,2,3")
	conf.Set("float64s", "1,2,3")
	conf.Set("strings", "a,b,c")
	conf.Set("durations", "1s,2s,3s")

	group1.Set("bool", "1")
	group1.Set("string", "efg")
	group1.Set("int", "456")
	group1.Set("int32", "456")
	group1.Set("int64", "456")
	group1.Set("uint", "456")
	group1.Set("uint32", "456")
	group1.Set("uint64", "456")
	group1.Set("float64", "456")
	group1.Set("duration", "456s")
	group1.Set("time", "2019-06-10T19:00:00Z")
	group1.Set("ints", "4,5,6")
	group1.Set("uints", "4,5,6")
	group1.Set("float64s", "4,5,6")
	group1.Set("strings", "e,f,g")
	group1.Set("durations", "4s,5s,6s")

	group2.Set("bool", "1")
	group2.Set("string", "xyz")
	group2.Set("int", "789")
	group2.Set("int32", "789")
	group2.Set("int64", "789")
	group2.Set("uint", "789")
	group2.Set("uint32", "789")
	group2.Set("uint64", "789")
	group2.Set("float64", "789")
	group2.Set("duration", "789s")
	group2.Set("time", "2019-06-10T20:00:00Z")
	group2.Set("ints", "7,8,9")
	group2.Set("uints", "7,8,9")
	group2.Set("float64s", "7,8,9")
	group2.Set("strings", "x,y,z")
	group2.Set("durations", "7s,8s,9s")

	////// Output

	fmt.Println("[DEFAULT]")
	fmt.Println(conf.GetBool("bool"))
	fmt.Println(conf.GetInt("int"))
	fmt.Println(conf.GetInt32("int32"))
	fmt.Println(conf.GetInt64("int64"))
	fmt.Println(conf.GetUint("uint"))
	fmt.Println(conf.GetUint32("uint32"))
	fmt.Println(conf.GetUint64("uint64"))
	fmt.Println(conf.GetFloat64("float64"))
	fmt.Println(conf.GetString("string"))
	fmt.Println(conf.GetDuration("duration"))
	fmt.Println(conf.GetTime("time").UTC())
	fmt.Println(conf.GetIntSlice("ints"))
	fmt.Println(conf.GetUintSlice("uints"))
	fmt.Println(conf.GetFloat64Slice("float64s"))
	fmt.Println(conf.GetStringSlice("strings"))
	fmt.Println(conf.GetDurationSlice("durations"))

	fmt.Printf("\n[%s]\n", group1.Name())
	fmt.Println(group1.GetBool("bool"))
	fmt.Println(group1.GetInt("int"))
	fmt.Println(group1.GetInt32("int32"))
	fmt.Println(group1.GetInt64("int64"))
	fmt.Println(group1.GetUint("uint"))
	fmt.Println(group1.GetUint32("uint32"))
	fmt.Println(group1.GetUint64("uint64"))
	fmt.Println(group1.GetFloat64("float64"))
	fmt.Println(group1.GetString("string"))
	fmt.Println(group1.GetDuration("duration"))
	fmt.Println(group1.GetTime("time").UTC())
	fmt.Println(group1.GetIntSlice("ints"))
	fmt.Println(group1.GetUintSlice("uints"))
	fmt.Println(group1.GetFloat64Slice("float64s"))
	fmt.Println(group1.GetStringSlice("strings"))
	fmt.Println(group1.GetDurationSlice("durations"))

	fmt.Printf("\n[%s]\n", group2.Name())
	fmt.Println(group2.GetBool("bool"))
	fmt.Println(group2.GetInt("int"))
	fmt.Println(group2.GetInt32("int32"))
	fmt.Println(group2.GetInt64("int64"))
	fmt.Println(group2.GetUint("uint"))
	fmt.Println(group2.GetUint32("uint32"))
	fmt.Println(group2.GetUint64("uint64"))
	fmt.Println(group2.GetFloat64("float64"))
	fmt.Println(group2.GetString("string"))
	fmt.Println(group2.GetDuration("duration"))
	fmt.Println(group2.GetTime("time").UTC())
	fmt.Println(group2.GetIntSlice("ints"))
	fmt.Println(group2.GetUintSlice("uints"))
	fmt.Println(group2.GetFloat64Slice("float64s"))
	fmt.Println(group2.GetStringSlice("strings"))
	fmt.Println(group2.GetDurationSlice("durations"))

	// Output:
	// [DEFAULT]
	// true
	// 123
	// 123
	// 123
	// 123
	// 123
	// 123
	// 123
	// abc
	// 2m3s
	// 2019-06-10 18:00:00 +0000 UTC
	// [1 2 3]
	// [1 2 3]
	// [1 2 3]
	// [a b c]
	// [1s 2s 3s]
	//
	// [group1]
	// true
	// 456
	// 456
	// 456
	// 456
	// 456
	// 456
	// 456
	// efg
	// 7m36s
	// 2019-06-10 19:00:00 +0000 UTC
	// [4 5 6]
	// [4 5 6]
	// [4 5 6]
	// [e f g]
	// [4s 5s 6s]
	//
	// [group1.group2]
	// true
	// 789
	// 789
	// 789
	// 789
	// 789
	// 789
	// 789
	// xyz
	// 13m9s
	// 2019-06-10 20:00:00 +0000 UTC
	// [7 8 9]
	// [7 8 9]
	// [7 8 9]
	// [x y z]
	// [7s 8s 9s]
}

func ExampleOptGroup_SetOptAlias() {
	conf := New()
	conf.RegisterOpts(IntOpt("newopt", "test alias").D(123))
	conf.SetOptAlias("oldopt", "newopt")

	fmt.Printf("newopt=%d, oldopt=%d\n", conf.GetInt("newopt"), conf.GetInt("oldopt"))
	conf.Set("oldopt", 456)
	fmt.Printf("newopt=%d, oldopt=%d\n", conf.GetInt("newopt"), conf.GetInt("oldopt"))

	// Output:
	// newopt=123, oldopt=123
	// newopt=456, oldopt=456
}

func TestOptGroupAlias(t *testing.T) {
	conf := New()
	conf.RegisterOpts(IntOpt("int", "test alias"))
	conf.SetOptAlias("opt", "int")

	if opt, exist := conf.Opt("opt"); !exist || opt.Name != "int" {
		t.Fail()
	} else if !conf.HasOpt("opt") {
		t.Fail()
	} else if conf.OptIsSet("opt") {
		t.Fail()
	} else if !conf.HasOptAndIsNotSet("opt") {
		t.Fail()
	}

	conf.Set("opt", 123)
	if !conf.OptIsSet("opt") {
		t.Fail()
	} else if conf.HasOptAndIsNotSet("opt") {
		t.Fail()
	} else if v := conf.GetInt("opt"); v != 123 {
		t.Error(v)
	}

	if conf.OptIsFrozen("opt") {
		t.Fail()
	}

	conf.FreezeOpt("opt")
	if !conf.OptIsFrozen("opt") {
		t.Fail()
	}

	conf.UnfreezeOpt("opt")
	if conf.OptIsFrozen("opt") {
		t.Fail()
	}
}

func TestOrValidator(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	conf := New()
	conf.RegisterOpts(StrOpt("ip1", "").V(Or(NewIPValidator(), NewEmptyStrValidator())))
	conf.RegisterOpts(StrOpt("ip2", "").D("0.0.0.0").V(Or(NewIPValidator(), NewEmptyStrValidator())))

	conf.Set("ip1", "127.0.0.1")
	conf.Set("ip2", "")

	if v := conf.GetString("ip1"); v != "127.0.0.1" {
		t.Error(v)
	} else if v = conf.GetString("ip2"); v != "" {
		t.Error(v)
	}
}

func TestOptGroup_SetOptAlias(t *testing.T) {
	conf := New()
	conf.RegisterOpts(StrOpt("opt", "").D("abc"))
	conf.SetOptAlias("opt1", "opt")
	conf.SetOptAlias("opt2", "opt")

	if conf.GetString("opt1") != "abc" {
		t.Fail()
	} else if conf.GetString("opt2") != "abc" {
		t.Fail()
	} else if aliases := conf.MustOpt("opt").Aliases; len(aliases) != 2 {
		t.Error(aliases)
	} else if aliases[0] != "opt1" {
		t.Error(aliases[0])
	} else if aliases[1] != "opt2" {
		t.Error(aliases[1])
	}
}
