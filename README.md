# Go Config [![Build Status](https://github.com/xgfone/gconf/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/gconf/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/gconf)](https://pkg.go.dev/github.com/xgfone/gconf/v6) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/gconf/master/LICENSE)

An extensible and powerful go configuration manager supporting `Go1.18+`, which is inspired by [oslo.config](https://github.com/openstack/oslo.config), [viper](https://github.com/spf13/viper) and [github.com/micro/go-micro/config](https://github.com/micro/go-micro/tree/master/config).


## Install
```shell
$ go get -u github.com/xgfone/gconf/v6
```


## Features

- A atomic key-value configuration center.
- Support kinds of decoders to decode the data from the source.
- Support to get the configuration data from many data sources.
- Support to change of the configuration option thread-safely during running.
- Support to observe the change of the configration options.


## Basic

```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf/v6"
)

// Pre-define a set of options.
var opts = []gconf.Opt{
	gconf.BoolOpt("opt1", "opt1 help doc"),
	gconf.StrOpt("opt2", "opt2 help doc").D("default"),
	gconf.IntOpt("opt3", "opt3 help doc").D(123).S("o"), // For short name
	gconf.Int32Opt("opt4", "opt4 help doc").As("opt5"),  // For alias name
	gconf.UintOpt("opt6", "opt6 help doc").D(1).V(gconf.NewIntegerRangeValidator(1, 100)),
	gconf.Float64Opt("opt7", "opt7 help doc").Cli(false),
}

func main() {
	// Register the options.
	gconf.RegisterOpts(opts...)

	// Print the registered options.
	for _, opt := range gconf.GetAllOpts() {
		fmt.Printf("Option: name=%s, value=%v\n", opt.Name, opt.Default)
	}

	// Add the observer to observe the change of the options.
	gconf.Observe(func(optName string, oldValue, newValue interface{}) {
		fmt.Printf("option=%s: %v -> %v\n", optName, oldValue, newValue)
	})

	// Update the value of the option thread-safely during app is running.
	gconf.Set("opt1", true)
	gconf.Set("opt2", "abc")
	gconf.Set("opt3", "456")
	gconf.Set("opt4", 789)
	gconf.Set("opt6", 100)
	gconf.Set("opt7", 1.2)

	// Get the values of the options thread-safely.
	fmt.Println(gconf.Get("opt1"))
	fmt.Println(gconf.Get("opt2"))
	fmt.Println(gconf.Get("opt3"))
	fmt.Println(gconf.Get("opt4"))
	fmt.Println(gconf.Get("opt6"))
	fmt.Println(gconf.Get("opt7"))

	// Output:
	// Option: name=opt1, value=false
	// Option: name=opt2, value=default
	// Option: name=opt3, value=123
	// Option: name=opt4, value=0
	// Option: name=opt6, value=1
	// Option: name=opt7, value=0
	// option=opt1: false -> true
	// option=opt2: default -> abc
	// option=opt3: 123 -> 456
	// option=opt4: 0 -> 789
	// option=opt6: 1 -> 100
	// option=opt7: 0 -> 1.2
	// true
	// abc
	// 456
	// 789
	// 100
	// 1.2

	// TODO ...
}
```


## Option Proxy

```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/gconf/v6"
)

func main() {
	// New the proxy of the option, which will new an option and register them,
	// then return the proxy of the option. So you can use the option proxy
	// to update and get the value of the option.
	opt1 := gconf.NewInt("opt1", 111, "opt1 help doc")
	opt2 := gconf.NewDuration("opt2", time.Second, "opt2 help doc")

	// Update the value of the option by the proxy.
	opt1.Set("222")
	opt2.Set("1m")

	// Get the value of the option by the proxy.
	fmt.Println(opt1.Get())
	fmt.Println(opt2.Get())

	// Output:
	// 222
	// 1m0s
}
```


## Option Group

```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/gconf/v6"
)

// Pre-define a set of options.
var opts = []gconf.Opt{
	gconf.StrOpt("opt1", "opt1 help doc").D("abc"),
	gconf.IntOpt("opt2", "opt2 help doc").D(123),
}

func main() {
	group1 := gconf.Group("group1")  // New the group "group1"
	group2 := group1.Group("group2") // New the sub-group "group1.group2"

	gconf.RegisterOpts(opts...)  // Register opts
	group1.RegisterOpts(opts...) // Register opts with group1
	group2.RegisterOpts(opts...) // Register opts with group2

	opt3 := group1.NewFloat64("opt3", 1.2, "opt3 help doc")          // For "group1.opt3"
	opt4 := group1.NewDuration("opt4", time.Second, "opt4 help doc") // For "group1.opt4"
	opt5 := group2.NewUint("opt5", 456, "opt5 help doc")             // For "group1.group2.opt5"
	opt6 := group2.NewBool("opt6", false, "opt6 help doc")           // For "group1.group2.opt6"

	/// Update the value of the option thread-safely during app is running.
	//
	// Method 1: Update the value of the option by the full name.
	gconf.Set("opt1", "aaa")
	gconf.Set("opt2", "111")
	gconf.Set("group1.opt1", "bbb")
	gconf.Set("group1.opt2", 222)
	gconf.Set("group1.opt3", 2.4)
	gconf.Set("group1.opt4", "1m")
	gconf.Set("group1.group2.opt1", "ccc")
	gconf.Set("group1.group2.opt2", 333)
	gconf.Set("group1.group2.opt5", 444)
	gconf.Set("group1.group2.opt6", "true")
	//
	// Method 2: Update the value of the option by the group proxy.
	group1.Set("opt1", "bbb")
	group1.Set("opt2", 222)
	group1.Set("opt3", 2.4)
	group1.Set("opt4", "1m")
	group2.Set("opt1", "ccc")
	group2.Set("opt2", 333)
	group2.Set("opt5", 444)
	group2.Set("opt6", "true")
	//
	// Method 3: Update the value of the option by the option proxy.
	opt3.Set(2.4)
	opt4.Set("1m")
	opt5.Set(444)
	opt6.Set("true")

	/// Get the values of the options thread-safely.
	//
	// Method 1: Get the value of the option by the full name.
	gconf.Get("opt1")
	gconf.Get("opt2")
	gconf.Get("group1.opt1")
	gconf.Get("group1.opt2")
	gconf.Get("group1.opt3")
	gconf.Get("group1.opt4")
	gconf.Get("group1.group2.opt1")
	gconf.Get("group1.group2.opt2")
	gconf.Get("group1.group2.opt5")
	gconf.Get("group1.group2.opt6")
	//
	// Method 2: Get the value of the option by the group proxy.
	group1.Get("opt1")
	group1.Get("opt2")
	group1.Get("opt3")
	group1.Get("opt4")
	group2.Get("opt1")
	group2.Get("opt2")
	group2.Get("opt5")
	group2.Get("opt6")
	//
	// Method 3: Get the value of the option by the option proxy.
	opt3.Get()
	opt4.Get()
	opt5.Get()
	opt6.Get()
}
```


## Data Decoder
The data decoder is a function like `func(src []byte, dst map[string]interface{}) error`, which is used to decode the configration data from the data source.

`Config` supports three kinds of decoders by default, such as `ini`, `json`, `yaml`, and `yml` is the alias of `yaml`. You can customize yourself decoder, then add it into `Config`, such as `Config.AddDecoder("type", NewCustomizedDecoder())`.


## Data Source
A source is used to read the configuration from somewhere the data is. And it can also watch the change the data.

```go
type Source interface {
	// String is the description of the source, such as "env", "file:/path/to".
	String() string

	// Read reads the source data once, which should not block.
	Read() (DataSet, error)

	// Watch watches the change of the source, then call the callback load.
	//
	// close is used to notice the underlying watcher to close and clean.
	Watch(close <-chan struct{}, load func(DataSet, error) (success bool))
}
```

You can load lots of sources to update the options. It has implemented the sources based on `flag`, `env`, `file` and `url`. But you can implement other sources, such as `ZooKeeper`, `ETCD`, etc.

```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/gconf/v6"
)

// Pre-define a set of options.
var opts = []gconf.Opt{
	gconf.StrOpt("opt1", "opt1 help doc").D("abc"),
	gconf.IntOpt("opt2", "opt2 help doc").D(123),
}

func main() {
	gconf.RegisterOpts(opts...)

	group := gconf.Group("group")
	group.RegisterOpts(opts...)

	// Convert the options to flag.Flag, and parse the CLI arguments with "flag".
	gconf.AddAndParseOptFlag(gconf.Conf)

	// Load the sources "flag" and "env".
	gconf.LoadSource(gconf.NewFlagSource())
	gconf.LoadSource(gconf.NewEnvSource(""))

	// Load and watch the file source.
	configFile := gconf.GetString(gconf.ConfigFileOpt.Name)
	gconf.LoadAndWatchSource(gconf.NewFileSource(configFile))

	for _, opt := range gconf.GetAllOpts() {
		fmt.Printf("%s: %v\n", opt.Name, gconf.Get(opt.Name))
	}

	gconf.Observe(func(optName string, oldValue, newValue interface{}) {
		fmt.Printf("%s: %v -> %v\n", optName, oldValue, newValue)
	})

	time.Sleep(time.Minute)

	// ## Run:
	// $ GROUP_OPT2=456 go run main.go --config-file conf.json
	// config-file: conf.json
	// group.opt1: abc
	// group.opt2: 456
	// opt1: abc
	// opt2: 123
	//
	// ## echo '{"opt1":"aaa","opt2":111,"group":{"opt1":"bbb","opt2": 222}}' > conf.json
	// opt1: abc -> aaa
	// opt2: 123 -> 111
	// group.opt1: abc -> bbb
	// group.opt2: 456 -> 222
	//
	// ## echo '{"opt1":"ccc","opt2":333,"group":{"opt1":"ddd","opt2":444}}' > conf.json
	// opt1: aaa -> ccc
	// opt2: 111 -> 333
	// group.opt1: bbb -> ddd
	// group.opt2: 222 -> 444
	//
	// ## echo '{"opt1":"eee","opt2":555,"group":{"opt1":"fff","opt2":666}}' > conf.json
	// opt1: ccc -> eee
	// opt2: 333 -> 555
	// group.opt1: ddd -> fff
	// group.opt2: 444 -> 666
}
```


## Snapshot & Backup

```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf/v6"
)

// Pre-define a set of options.
var opts = []gconf.Opt{
	gconf.StrOpt("opt1", "opt1 help doc").D("abc"),
	gconf.IntOpt("opt2", "opt2 help doc").D(123),
}

func main() {
	gconf.RegisterOpts(opts...)

	group := gconf.Group("group")
	group.RegisterOpts(opts...)

	// Convert the options to flag.Flag, and parse the CLI arguments with "flag".
	gconf.AddAndParseOptFlag(gconf.Conf)

	// Load the sources "flag" and "env".
	gconf.LoadSource(gconf.NewFlagSource())
	gconf.LoadSource(gconf.NewEnvSource(""))

	// Load and update the configuration from the backup file which will watch
	// the change of all configuration options and write them into the backup
	// file to wait to be loaded when the program starts up next time.
	gconf.LoadBackupFile("config-file.backup")

	fmt.Println(gconf.Get("opt1"))
	fmt.Println(gconf.Get("opt2"))
	fmt.Println(group.Get("opt1"))
	fmt.Println(group.Get("opt2"))

	/// Get the snapshot of all configuration options at any time.
	// generation, snapshots := gconf.Snapshot()
	// fmt.Println(generation, snapshots)

	// $ go run main.go
	// abc
	// 123
	// abc
	// 123
	//
	// $ echo '{"opt1":"aaa","opt2":111,"group":{"opt1":"bbb","opt2":222}}' > config-file.backup
	// $ go run main.go
	// aaa
	// 111
	// bbb
	// 222
}
```
