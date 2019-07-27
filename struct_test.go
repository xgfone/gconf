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
	"flag"
	"fmt"
	"os"
	"time"
)

func ExampleOptField() {
	type AppConfig struct {
		Bool      BoolOptField
		BoolT     BoolTOptField
		Int       IntOptField
		Int32     Int32OptField
		Int64     Int64OptField
		Uint      UintOptField
		Uint32    Uint32OptField
		Uint64    Uint64OptField
		Float64   Float64OptField
		String    StringOptField
		Duration  DurationOptField
		Time      TimeOptField
		Ints      IntSliceOptField
		Uints     UintSliceOptField
		Float64s  Float64SliceOptField
		Strings   StringSliceOptField
		Durations DurationSliceOptField

		// Pointer Example
		IntP   *IntOptField `default:"123"`
		Ignore *StringOptField
	}

	// Notice: for the pointer to the option field, it must be initialized.
	// Or it will be ignored.
	config := AppConfig{IntP: &IntOptField{}}
	conf := New()
	conf.RegisterStruct(&config)

	fmt.Println("--- Registered Options ---")
	for _, opt := range conf.AllOpts() {
		fmt.Println(opt.Name)
	}

	fmt.Println("--- Before Updating ---")
	fmt.Printf("bool=%v\n", config.Bool.Get())
	fmt.Printf("boolt=%v\n", config.BoolT.Get())
	fmt.Printf("int=%v\n", config.Int.Get())
	fmt.Printf("int32=%v\n", config.Int32.Get())
	fmt.Printf("int64=%v\n", config.Int64.Get())
	fmt.Printf("uint=%v\n", config.Uint.Get())
	fmt.Printf("uint32=%v\n", config.Uint32.Get())
	fmt.Printf("uint64=%v\n", config.Uint64.Get())
	fmt.Printf("float64=%v\n", config.Float64.Get())
	fmt.Printf("string=%v\n", config.String.Get())
	fmt.Printf("duration=%v\n", config.Duration.Get())
	fmt.Printf("time=%v\n", config.Time.Get().Format(time.RFC3339))
	fmt.Printf("ints=%v\n", config.Ints.Get())
	fmt.Printf("uints=%v\n", config.Uints.Get())
	fmt.Printf("float64s=%v\n", config.Float64s.Get())
	fmt.Printf("strings=%v\n", config.Strings.Get())
	fmt.Printf("durations=%v\n", config.Durations.Get())
	fmt.Printf("intp=%v\n", config.IntP.Get())

	conf.Set("bool", true)
	conf.Set("boolt", false)
	conf.Set("int", 123)
	conf.Set("int32", 123)
	conf.Set("int64", 123)
	conf.Set("uint", 123)
	conf.Set("uint32", 123)
	conf.Set("uint64", 123)
	conf.Set("float64", 123)
	conf.Set("string", "abc")
	conf.Set("duration", "10s")
	conf.Set("time", "2019-07-27 15:39:34")
	conf.Set("ints", []int{1, 2, 3})
	conf.Set("uints", []uint{4, 5, 6})
	conf.Set("float64s", []float64{1, 2, 3})
	conf.Set("strings", []string{"a", "b", "c"})
	conf.Set("durations", []time.Duration{time.Second, time.Second * 2, time.Second * 3})
	conf.Set("intp", 456)

	fmt.Println("--- After Updating ---")
	fmt.Printf("bool=%v\n", config.Bool.Get())
	fmt.Printf("boolt=%v\n", config.BoolT.Get())
	fmt.Printf("int=%v\n", config.Int.Get())
	fmt.Printf("int32=%v\n", config.Int32.Get())
	fmt.Printf("int64=%v\n", config.Int64.Get())
	fmt.Printf("uint=%v\n", config.Uint.Get())
	fmt.Printf("uint32=%v\n", config.Uint32.Get())
	fmt.Printf("uint64=%v\n", config.Uint64.Get())
	fmt.Printf("float64=%v\n", config.Float64.Get())
	fmt.Printf("string=%v\n", config.String.Get())
	fmt.Printf("duration=%v\n", config.Duration.Get())
	fmt.Printf("time=%v\n", config.Time.Get().Format(time.RFC3339))
	fmt.Printf("ints=%v\n", config.Ints.Get())
	fmt.Printf("uints=%v\n", config.Uints.Get())
	fmt.Printf("float64s=%v\n", config.Float64s.Get())
	fmt.Printf("strings=%v\n", config.Strings.Get())
	fmt.Printf("durations=%v\n", config.Durations.Get())
	fmt.Printf("intp=%v\n", config.IntP.Get())

	// Output:
	// --- Registered Options ---
	// bool
	// boolt
	// duration
	// durations
	// float64
	// float64s
	// int
	// int32
	// int64
	// intp
	// ints
	// string
	// strings
	// time
	// uint
	// uint32
	// uint64
	// uints
	// --- Before Updating ---
	// bool=false
	// boolt=true
	// int=0
	// int32=0
	// int64=0
	// uint=0
	// uint32=0
	// uint64=0
	// float64=0
	// string=
	// duration=0s
	// time=0001-01-01T00:00:00Z
	// ints=[]
	// uints=[]
	// float64s=[]
	// strings=[]
	// durations=[]
	// intp=123
	// --- After Updating ---
	// bool=true
	// boolt=false
	// int=123
	// int32=123
	// int64=123
	// uint=123
	// uint32=123
	// uint64=123
	// float64=123
	// string=abc
	// duration=10s
	// time=2019-07-27T15:39:34Z
	// ints=[1 2 3]
	// uints=[4 5 6]
	// float64s=[1 2 3]
	// strings=[a b c]
	// durations=[1s 2s 3s]
	// intp=456
}

func ExampleOptGroup_RegisterStruct() {
	type Group struct {
		Bool     bool          `help:"test bool"`
		Int      int           `default:"123" help:"test int"`
		Int32    int32         `default:"123" help:"test int32"`
		Int64    int64         `default:"123" help:"test int64"`
		Uint     uint          `default:"123" help:"test uint"`
		Uint32   uint32        `default:"123" help:"test uint32"`
		Uint64   uint64        `default:"123" help:"test uint64"`
		Float64  float64       `default:"123" help:"test float64"`
		String   string        `default:"abc" help:"test string"`
		Duration time.Duration `default:"123s" help:"test time.Duration"`
		Time     time.Time     `help:"test time.Time"`

		Ints      []int           `default:"1,2,3" help:"test []int"`
		Uints     []uint          `default:"1,2,3" help:"test []uint"`
		Float64s  []float64       `default:"1,2,3" help:"test []float64"`
		Strings   []string        `default:"a,b,c" help:"test []string"`
		Durations []time.Duration `default:"1s,2s,3s" help:"test []time.Duration"`
	}

	type WrapGroup struct {
		Bool     bool          `help:"test bool"`
		Int      int           `default:"456" help:"test int"`
		Int32    int32         `default:"456" help:"test int32"`
		Int64    int64         `default:"456" help:"test int64"`
		Uint     uint          `default:"456" help:"test uint"`
		Uint32   uint32        `default:"456" help:"test uint32"`
		Uint64   uint64        `default:"456" help:"test uint64"`
		Float64  float64       `default:"456" help:"test float64"`
		String   string        `default:"efg" help:"test string"`
		Duration time.Duration `default:"456s" help:"test time.Duration"`
		Time     time.Time     `help:"test time.Time"`

		Ints      []int           `default:"4,5,6" help:"test []int"`
		Uints     []uint          `default:"4,5,6" help:"test []uint"`
		Float64s  []float64       `default:"4,5,6" help:"test []float64"`
		Strings   []string        `default:"e,f,g" help:"test []string"`
		Durations []time.Duration `default:"4s,5s,6s" help:"test []time.Duration"`

		Group Group `group:"group3" name:"group33"`
	}

	type DataConfig struct {
		Bool     bool          `help:"test bool"`
		Int      int           `default:"789" help:"test int"`
		Int32    int32         `default:"789" help:"test int32"`
		Int64    int64         `default:"789" help:"test int64"`
		Uint     uint          `default:"789" help:"test uint"`
		Uint32   uint32        `default:"789" help:"test uint32"`
		Uint64   uint64        `default:"789" help:"test uint64"`
		Float64  float64       `default:"789" help:"test float64"`
		String   string        `default:"xyz" help:"test string"`
		Duration time.Duration `default:"789s" help:"test time.Duration"`
		Time     time.Time     `help:"test time.Time"`

		Ints      []int           `default:"7,8,9" help:"test []int"`
		Uints     []uint          `default:"7,8,9" help:"test []uint"`
		Float64s  []float64       `default:"7,8,9" help:"test []float64"`
		Strings   []string        `default:"x,y,z" help:"test []string"`
		Durations []time.Duration `default:"7s,8s,9s" help:"test []time.Duration"`

		Group1 Group     `group:"group1"`
		Group2 WrapGroup `name:"group2"`
	}

	// Register the option from struct
	var data DataConfig
	conf := New()
	conf.RegisterStruct(&data)

	// Add options to flag, and parse them from flag.
	flagSet := flag.NewFlagSet("test_struct", flag.ExitOnError)
	AddOptFlag(conf, flagSet)
	flagSet.Parse([]string{
		"--bool=true",
		"--time=2019-06-11T20:00:00Z",
		"--group1.bool=1",
	})
	conf.LoadSource(NewFlagSource(flagSet))

	fmt.Println("--- Struct ---")
	fmt.Printf("bool: %t\n", data.Bool)
	fmt.Printf("int: %d\n", data.Int)
	fmt.Printf("int32: %d\n", data.Int32)
	fmt.Printf("int64: %d\n", data.Int64)
	fmt.Printf("uint: %d\n", data.Uint)
	fmt.Printf("uint32: %d\n", data.Uint32)
	fmt.Printf("uint64: %d\n", data.Uint64)
	fmt.Printf("float64: %v\n", data.Float64)
	fmt.Printf("string: %s\n", data.String)
	fmt.Printf("duration: %s\n", data.Duration)
	fmt.Printf("time: %s\n", data.Time)
	fmt.Printf("ints: %v\n", data.Ints)
	fmt.Printf("uints: %v\n", data.Uints)
	fmt.Printf("float64s: %v\n", data.Float64s)
	fmt.Printf("strings: %v\n", data.Strings)
	fmt.Printf("durations: %v\n", data.Durations)
	// ...
	fmt.Println("--- Config ---")
	fmt.Printf("bool: %t\n", conf.GetBool("bool"))
	fmt.Printf("int: %d\n", conf.GetInt("int"))
	fmt.Printf("int32: %d\n", conf.GetInt32("int32"))
	fmt.Printf("int64: %d\n", conf.GetInt64("int64"))
	fmt.Printf("uint: %d\n", conf.GetUint("uint"))
	fmt.Printf("uint32: %d\n", conf.GetUint32("uint32"))
	fmt.Printf("uint64: %d\n", conf.GetUint64("uint64"))
	fmt.Printf("float64: %v\n", conf.GetFloat64("float64"))
	fmt.Printf("string: %s\n", conf.GetString("string"))
	fmt.Printf("duration: %s\n", conf.GetDuration("duration"))
	fmt.Printf("time: %s\n", conf.GetTime("time"))
	fmt.Printf("ints: %v\n", conf.GetIntSlice("ints"))
	fmt.Printf("uints: %v\n", conf.GetUintSlice("uints"))
	fmt.Printf("float64s: %v\n", conf.GetFloat64Slice("float64s"))
	fmt.Printf("strings: %v\n", conf.GetStringSlice("strings"))
	fmt.Printf("durations: %v\n", conf.GetDurationSlice("durations"))
	// ...
	conf.PrintGroup(os.Stdout)

	// Output:
	// --- Struct ---
	// bool: true
	// int: 789
	// int32: 789
	// int64: 789
	// uint: 789
	// uint32: 789
	// uint64: 789
	// float64: 789
	// string: xyz
	// duration: 13m9s
	// time: 2019-06-11 20:00:00 +0000 UTC
	// ints: [7 8 9]
	// uints: [7 8 9]
	// float64s: [7 8 9]
	// strings: [x y z]
	// durations: [7s 8s 9s]
	// --- Config ---
	// bool: true
	// int: 789
	// int32: 789
	// int64: 789
	// uint: 789
	// uint32: 789
	// uint64: 789
	// float64: 789
	// string: xyz
	// duration: 13m9s
	// time: 2019-06-11 20:00:00 +0000 UTC
	// ints: [7 8 9]
	// uints: [7 8 9]
	// float64s: [7 8 9]
	// strings: [x y z]
	// durations: [7s 8s 9s]
	// [DEFAULT]
	//     bool
	//     duration
	//     durations
	//     float64
	//     float64s
	//     int
	//     int32
	//     int64
	//     ints
	//     string
	//     strings
	//     time
	//     uint
	//     uint32
	//     uint64
	//     uints
	// [group1]
	//     bool
	//     duration
	//     durations
	//     float64
	//     float64s
	//     int
	//     int32
	//     int64
	//     ints
	//     string
	//     strings
	//     time
	//     uint
	//     uint32
	//     uint64
	//     uints
	// [group2]
	//     bool
	//     duration
	//     durations
	//     float64
	//     float64s
	//     int
	//     int32
	//     int64
	//     ints
	//     string
	//     strings
	//     time
	//     uint
	//     uint32
	//     uint64
	//     uints
	// [group2.group3]
	//     bool
	//     duration
	//     durations
	//     float64
	//     float64s
	//     int
	//     int32
	//     int64
	//     ints
	//     string
	//     strings
	//     time
	//     uint
	//     uint32
	//     uint64
	//     uints
}
