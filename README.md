# gconf [![Build Status](https://travis-ci.org/xgfone/gconf.svg?branch=master)](https://travis-ci.org/xgfone/gconf) [![GoDoc](https://godoc.org/github.com/xgfone/gconf?status.svg)](http://godoc.org/github.com/xgfone/gconf) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/gconf/master/LICENSE)
An extensible and powerful go configuration manager, which is inspired by [oslo.config](https://github.com/openstack/oslo.config), [github.com/micro/go-micro/config](https://github.com/micro/go-micro/tree/master/config) and [viper](https://github.com/spf13/viper).

The current version is `v2`. See [DOC](https://godoc.org/github.com/xgfone/gconf).

The supported Go version: `1.10+`.


## Goal

1. A atomic key-value configuration center with the multi-group and the option.
2. Support the multi-parser to parse the configurations from many sources with the different format.
3. Change the configuration dynamically during running and watch it.
4. Observe the change of the configuration.


## Source

Source is used to read the configuration data. You can load lots of sources to read the configuration data from many storage locations. The default has implemented some sources, such as `flag`, `cli`, `env`, `file`, `url`. But you can also implement other sources, such as `ZooKeeper`, `ETCD`, etc.

**Notice:** If the source supports the watcher, it will add it to watch the changed of the source data automatically.

## Decoder

The source reads the original data, that's `[]byte`, and it must be decoded. The default has implemented the `json` and `INI` decoders.


## Read and Update the option value

You can get the group or sub-group by calling `Group(name)`, then get the option value again by calling `GetBool("optionName")`, `GetString("optionName")`, `GetInt("optionName")`, etc. However, if you don't known whether the option has a value, you can call `Get("optionName")`, which returns `nil` if no option or value, etc.

Beside, you can update the value of the option dynamically by calling `UpdateOptValue(groupFullName, optName, newOptValue)` during the program is running. For the default group, `groupFullName` is a empty string(`""`).

**Notce:**
1. Both of Reading and Updating are goroutine-safe.
2. For the modifiable type, such as slice or map, in order to update them, you should clone them firstly, then update the cloned option vlaue and call `UpdateOptValue` with it.


## Observe the changed configuration

You can use the method `Observe(callback func(groupName, optName string, oldOptValue, newOptValue interface{}))` to monitor what the configuration is updated to: when a certain configuration is updated, the callback function will be called asynchronizedly.


## Usage
```go
package main

import (
	"flag"
	"fmt"

	"github.com/xgfone/gconf"
)

func main() {
	// Register options
	conf := gconf.New()
	conf.RegisterOpt(gconf.StrOpt("ip", "the ip address").D("0.0.0.0").V(gconf.NewIPValidator()))
	conf.RegisterOpt(gconf.IntOpt("port", "the port").D(80))
	conf.NewGroup("redis").RegisterOpt(gconf.StrOpt("conn", "the redis connection url"))

	// Set the CLI version and exit when giving the CLI option version.
	conf.SetVersion(gconf.VersionOpt.D("1.0.0"))
	gconf.AddAndParseOptFlag(conf)

	// Load the sources
	conf.LoadSource(gconf.NewFlagSource())

	// Read and print the option
	fmt.Println(conf.GetString("ip"))
	fmt.Println(conf.GetInt("port"))
	fmt.Println(conf.Group("redis").GetString("conn"))
	fmt.Println(flag.Args())

	// Execute:
	//     PROGRAM --ip 1.2.3.4 --redis.conn=redis://127.0.0.1:6379/0 aa bb cc
	//
	// Output:
	//     1.2.3.4
	//     80
	//     redis://127.0.0.1:6379/0
	//     [aa bb cc]
}
```

The package has created a global default `Config`, that's, `Conf`. You can use it, like the global variable `CONF` in `oslo.config`. For example,
```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

var opts = []gconf.Opt{
	gconf.StrOpt("ip", "the ip address").D("0.0.0.0").V(gconf.NewIPValidator()),
	gconf.IntOpt("port", "the port").D(80).V(gconf.NewPortValidator()),
}

func main() {
	// Register options
	gconf.Conf.RegisterOpts(opts)

	// Add the options to flag.CommandLine and parse the CLI
	gconf.AddAndParseOptFlag(gconf.Conf)

	// Load the sources
	gconf.Conf.LoadSource(gconf.NewFlagSource())

	// Read and print the option
	fmt.Println(gconf.Conf.GetString("ip"))
	fmt.Println(gconf.Conf.GetInt("port"))

	// Execute:
	//     PROGRAM --ip 1.2.3.4
	//
	// Output:
	//     1.2.3.4
	//     80
}
```

You can watch the change of the configuration option.
```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/gconf"
)

func main() {
	// Register the options
	gconf.Conf.RegisterOpt(gconf.StrOpt("opt1", "").D("abc"))
	gconf.Conf.NewGroup("group").RegisterOpt(gconf.IntOpt("opt2", "").D(123))

	// Add the observer
	gconf.Conf.Observe(func(group, opt string, old, new interface{}) {
		fmt.Printf("[Observer] Setting: group=%s, opt=%s, old=%v, new=%v\n", group, opt, old, new)
	})

	// Update the value of the option.
	gconf.Conf.UpdateOptValue("", "opt1", "xyz") // The first way
	gconf.Conf.Group("group").Set("opt2", 789)   // The second way

	// Sleep a while, because observer is asynchronized.
	time.Sleep(time.Millisecond * 10)

	// Output:
	// [Observer] Setting: group=, opt=opt1, old=abc, new=xyz
	// [Observer] Setting: group=group, opt=opt2, old=123, new=789
}
```

### The `cli` Command

The `flag` does not support the command, so you can use `github.com/urfave/cli`.

```go
package main

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/xgfone/gconf"
)

func main() {
	// Register options into the group
	gconf.Conf.RegisterOpt(gconf.StrOpt("opt1", "").D("abc"))
	gconf.Conf.NewGroup("cmd1").RegisterOpt(gconf.IntOpt("opt2", ""))
	gconf.Conf.NewGroup("cmd1").NewGroup("cmd2").RegisterOpt(gconf.IntOpt("opt3", ""))

	// Create and run cli app.
	app := cli.NewApp()
	app.Flags = []cli.Flag{cli.StringFlag{Name: "opt1"}}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "cmd1",
			Flags: []cli.Flag{cli.IntFlag{Name: "opt2"}},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "cmd2",
					Flags: []cli.Flag{cli.IntFlag{Name: "opt3"}},
					Action: func(ctx *cli.Context) error {
						// Load the sources
						gconf.Conf.LoadSource(gconf.NewCliSource(ctx, "cmd1.cmd2"))          // cmd2
						gconf.Conf.LoadSource(gconf.NewCliSource(ctx.Parent(), "cmd1"))      // cmd1
						gconf.Conf.LoadSource(gconf.NewCliSource(ctx.Parent().Parent(), "")) // global

						// Read and print the option
						fmt.Println(gconf.Conf.GetString("opt1"))
						fmt.Println(gconf.Conf.Group("cmd1").GetInt("opt2"))
						fmt.Println(gconf.Conf.Group("cmd1.cmd2").GetInt("opt3"))

						return nil
					},
				},
			},
		},
	}
	app.RunAndExitOnError()

	// Execute:
	//     PROGRAM --opt1=xyz cmd1 --opt2=123 cmd2 --opt3=456
	//
	// Output:
	//     xyz
	//     123
	//     456
}
```

### Use the config file

The default file source supports watcher, it will watch the change of the given filename then reload it.

```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/gconf"
)

var opts = []gconf.Opt{
	gconf.StrOpt("ip", "the ip address").D("0.0.0.0").V(gconf.NewIPValidator()),
	gconf.IntOpt("port", "the port").D(80).V(gconf.NewPortValidator()),
}

func main() {
	// Register options
	//
	// Notice: the default global Conf has registered gconf.O.
	gconf.Conf.RegisterOpts(opts)

	// Add the options to flag.CommandLine and parse the CLI
	gconf.AddAndParseOptFlag(gconf.Conf)

	// Load the flag & file sources
	gconf.Conf.LoadSource(gconf.NewFlagSource())
	gconf.Conf.LoadSource(gconf.NewFileSource(gconf.Conf.GetString(gconf.ConfigFileOpt.Name)))

	// Read and print the option
	for {
		time.Sleep(time.Second * 10)
		fmt.Printf("%s:%d\n", gconf.Conf.GetString("ip"), gconf.Conf.GetInt("port"))
	}

	// $ PROGRAM --config-file /path/to/file.json &
	// 0.0.0.0:80
	//
	// $ echo '{"ip": "1.2.3.4", "port":8000}' >config.json
	// 1.2.3.4:8000
	//
	// $ echo '{"ip": "5.6.7.8", "port":9000}' >config.json
	// 5.6.7.8:9000
}
```

**Notice:** Because there are two kinds of default decoders, `json` and `INI`, the file is only the format of json or ini. But you can register other decoders to support more format files.

### Register a struct as the group and the option
You also register a struct then use it.
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xgfone/gconf"
)

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

func main() {
	// Register the option from struct
	var data DataConfig
	conf := gconf.New()
	conf.RegisterStruct(&data)

	// Add options to flag, and parse them from flag.
	gconf.AddAndParseOptFlag(conf)
	conf.LoadSource(gconf.NewFlagSource())

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
	fmt.Println("--- Group ---")
	conf.PrintGroup(os.Stdout)

	// RUN:
	//     PROGRAM --bool=true --time=2019-06-11T20:00:00Z --group1.bool=1
	//
	// Output:
	//     --- Struct ---
	//     bool: true
	//     int: 789
	//     int32: 789
	//     int64: 789
	//     uint: 789
	//     uint32: 789
	//     uint64: 789
	//     float64: 789
	//     string: xyz
	//     duration: 13m9s
	//     time: 2019-06-11 20:00:00 +0000 UTC
	//     ints: [7 8 9]
	//     uints: [7 8 9]
	//     float64s: [7 8 9]
	//     strings: [x y z]
	//     durations: [7s 8s 9s]
	//     --- Config ---
	//     bool: true
	//     int: 789
	//     int32: 789
	//     int64: 789
	//     uint: 789
	//     uint32: 789
	//     uint64: 789
	//     float64: 789
	//     string: xyz
	//     duration: 13m9s
	//     time: 2019-06-11 20:00:00 +0000 UTC
	//     ints: [7 8 9]
	//     uints: [7 8 9]
	//     float64s: [7 8 9]
	//     strings: [x y z]
	//     durations: [7s 8s 9s]
	//     --- Group ---
	//     [DEFAULT]
	//         bool
	//         duration
	//         durations
	//         float64
	//         float64s
	//         int
	//         int32
	//         int64
	//         ints
	//         string
	//         strings
	//         time
	//         uint
	//         uint32
	//         uint64
	//         uints
	//     [group1]
	//         bool
	//         duration
	//         durations
	//         float64
	//         float64s
	//         int
	//         int32
	//         int64
	//         ints
	//         string
	//         strings
	//         time
	//         uint
	//         uint32
	//         uint64
	//         uints
	//     [group2]
	//         bool
	//         duration
	//         durations
	//         float64
	//         float64s
	//         int
	//         int32
	//         int64
	//         ints
	//         string
	//         strings
	//         time
	//         uint
	//         uint32
	//         uint64
	//         uints
	//     [group2.group3]
	//         bool
	//         duration
	//         durations
	//         float64
	//         float64s
	//         int
	//         int32
	//         int64
	//         ints
	//         string
	//         strings
	//         time
	//         uint
	//         uint32
	//         uint64
	//         uints
}
```
