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
