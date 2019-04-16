# gconf [![Build Status](https://travis-ci.org/xgfone/gconf.svg?branch=master)](https://travis-ci.org/xgfone/gconf) [![GoDoc](https://godoc.org/github.com/xgfone/gconf?status.svg)](http://godoc.org/github.com/xgfone/gconf) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/gconf/master/LICENSE)
An extensible and powerful go configuration manager, which is inspired by [oslo.config](https://github.com/openstack/oslo.config).

The current version is `v1`. See [DOC](https://godoc.org/github.com/xgfone/gconf).

The supported Go version: `1.8+`.


## Goal

1. A atomic key-value configuration center with the multi-group and the option.
2. Support the multi-parser to parse the configurations from many sources with the different format.
3. Change the configuration dynamically during running and watch it.
4. Observe the change of the configuration.


## Parser

The parser is used to parse the configurations from many sources with a certain format. There have been implemented some parsers, such as `cli`, `flag`, `env`, `INI`, `property`, etc. Of these, `flag` is the CLI parser based on the stdlib `flag`, and `cli` is another CLI parser based on `github.com/urfave/cli`.

You can develop yourself parser, only needing to implement the interface `Parser` as follow.
```go
type Parser interface {
	// Name returns the name of the parser to identify it.
	Name() string

	// Priority reports the priority of the current parser, which should be
	// a natural number.
	//
	// The smaller the number, the higher the priority. And the higher priority
	// parser can cover the option value set by the lower priority parser.
	//
	// For the cli parser, it maybe return 0 to indicate the highest priority.
	Priority() int

	// Pre is called before parsing the configuration, so it may be used to
	// initialize the parser, such as registering the itself options.
	Pre(*Config) error

	// Parse the value of the registered options.
	//
	// The parser can get any information from the argument, config.
	//
	// When the parser parsed out the option value, it should call
	// config.SetOptValue(), which will set the group option.
	// For the default group, the group name may be "" instead,
	//
	// For the CLI parser, it should get the parsed CLI argument by calling
	// config.ParsedCliArgs(), which is a string slice, not nil, but it maybe
	// have no elements. The CLI parser should not use os.Args[1:]
	// as the parsed CLI arguments. After parsing, If there are the rest CLI
	// arguments, which are those that does not start with the prefix "-", "--",
	// the CLI parser should call config.SetCliArgs() to set them.
	//
	// If there is any error, the parser should stop to parse and return it.
	//
	// If a certain option has no value, the parser should not return a default
	// one instead. Also, the parser has no need to convert the value to the
	// corresponding specific type, and just string is ok. Because the Config
	// will convert the value to the specific type automatically. Certainly,
	// it's not harmless for the parser to convert the value to the specific type.
	Parse(*Config) error

	// Pre is called before parsing the configuration, so it may be used to
	// clean the parser.
	Post(*Config) error
}
```

`Config` does not distinguish the CLI parser and the common parser, all of which have the same interface `Parser`. You can add them by calling `AddParser()`. See the example below.

**Notice:** the priority of the CLI parser should be higher than that of other parsers, because the higher parser will be parsed preferentially. And the same priority parsers will be parsed in turn by the added order.


## Read and Update the option value

You can get the group or sub-group by calling `Group(name)`, then get the option value again by calling `Bool("optionName")`, `String("optionName")`, `Int("optionName")`, etc. However, if you don't known whether the option has a value, you can call `Value("optionName")`, which returns `nil` if no value, or call `BoolE()`, `BoolD()`, `StringE()`, `StringD()`, `IntE()`, `IntD()`, etc.

Beside, you can update the value of the option dynamically by calling `SetOptValue(priority int, groupFullName, optName, newOptValue) error`, during the program is running. For the default group, `groupFullName` may be `""`. If the setting fails, it will return an error.

**Notce:**
1. All of Reading and Updating are goroutine-safe.
2. For the modifiable type, such as slice or map, in order to update them, you should clone them firstly, then update the cloned option vlaue and call `SetOptValue` with it.
3. For updating the value of the option, you can use the priority `0`, which is the highest priority and will overwrite any other value.


## Observe the changed configuration

You can use the method `Observe(callback func(groupFullName, optName string, oldOptValue, newOptValue interface{}))` to monitor what the configuration is updated to: when a certain configuration is updated, the callback function will be called.

Notice: the callback should finish as soon as possible because the callback is called synchronously at when the configuration is updated.


## Usage
```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

func main() {
	cliParser := gconf.NewDefaultFlagCliParser(true)
	iniParser := gconf.NewSimpleIniParser(gconf.ConfigFileName)
	conf := gconf.New().AddParser(cliParser, iniParser)

	ipOpt := gconf.StrOpt("", "ip", "", "the ip address").SetValidators(gconf.NewIPValidator())
	conf.RegisterCliOpt(ipOpt)
	conf.RegisterCliOpt(gconf.IntOpt("", "port", 80, "the port"))
	conf.NewGroup("redis").RegisterCliOpt(gconf.StrOpt("", "conn", "redis://127.0.0.1:6379/0", "the redis connection url"))

	// Print the version and exit when giving the CLI option version.
	conf.SetCliVersion("v", "version", "1.0.0", "Print the version and exit")

	if err := conf.Parse(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(conf.String("ip"))
	fmt.Println(conf.Int("port"))
	fmt.Println(conf.Group("redis").String("conn"))
	fmt.Println(conf.CliArgs())

	// Execute:
	//     PROGRAM --ip 0.0.0.0 aa bb cc
	//
	// Output:
	//     0.0.0.0
	//     80
	//     redis://127.0.0.1:6379/0
	//     [aa bb cc]
}
```

You can also create a new `Config` by the `NewDefault()`, which will use `NewDefaultFlagCliParser(true)` as the CLI parser, add the ini parser `NewSimpleIniParser(ConfigFileName)` that will register the CLI option `config-file`. **Notice:** `NewDefault()` does not add the environment variable parser, and you need to add it by hand, such as `NewDefault().AddParser(NewEnvVarParser(""))`.

The package has created a global default `Config` created by `NewDefault()`, that's, `Conf`. You can use it, like the global variable `CONF` in `oslo.config`. For example,
```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

var opts = []gconf.Opt{
	gconf.Str("ip", "", "the ip address").AddValidators(gconf.NewIPValidator()),
	gconf.Int("port", 80, "the port").AddValidators(gconf.NewPortValidator()),
}

func main() {
	gconf.Conf.RegisterCliOpts(opts)
	gconf.Conf.Parse()

	fmt.Println(gconf.Conf.String("ip"))
	fmt.Println(gconf.Conf.Int("port"))

	// Execute:
	//     PROGRAM --ip 0.0.0.0
	//
	// Output:
	//     0.0.0.0
	//     80
}
```

You can watch the change of the configuration option.
```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

func main() {
	conf := gconf.New()

	conf.NewGroup("test").RegisterCliOpt(gconf.Str("watchval", "abc", "test watch value"))
	conf.Observe(func(groupName, optName string, old, new interface{}) {
		fmt.Printf("[Observer] group=%s, option=%s, old=%v, new=%v\n",
			groupName, optName, old, new)
	})

	conf.Parse()

	// Set the option vlaue during the program is running.
	conf.SetOptValue(0, "test", "watchval", "123")

	// Output:
	// [Observer] group=test, option=watchval, old=<nil>, new=abc
	// [Observer] group=test, option=watchval, old=abc, new=123
}
```

### Register a struct as the group and the option
You also register a struct then use it.
```go
package main

import (
	"fmt"
	"os"

	"github.com/xgfone/gconf"
)

type MySQL struct {
	Conn       string `help:"the connection to mysql server"`
	MaxConnNum int    `name:"maxconn" default:"3" help:"the maximum number of connections"`
}

type DB struct {
	MySQL MySQL
}

type DBWrapper struct {
	DB1 DB
	DB2 DB `group:"db222"`
}

type Config struct {
	Addr  string `default:":80" help:"the address to listen to"`
	File  string `default:"" group:"log" help:"the log file path"`
	Level string `default:"debug" group:"log" help:"the log level, such as debug, info, etc"`

	Ignore bool `name:"-" default:"true"`

	DB1 DB `cli:"false"`
	DB2 DB `cli:"off" name:"db02"`
	DB3 DB `cli:"f" group:"db03"`              // equal to `name:"db03"` for this case
	DB4 DB `cli:"0" name:"db04" group:"db004"` // use the tag "group", not "name"

	DB5 DBWrapper `group:"db"`
}

func main() {
	var config Config
	gconf.Conf.AddParser(gconf.NewEnvVarParser(10, "")).RegisterCliStruct(&config)

	// Only for test
	os.Setenv("DB1_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db1")
	os.Setenv("DB02_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db2")
	os.Setenv("DB03_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db3")
	os.Setenv("DB004_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db4")
	os.Setenv("DB_DB1_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db5-1")
	os.Setenv("DB_DB222_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db5-2")
	cliArgs := []string{"--addr", "0.0.0.0:80", "--log.file", "/var/log/test.log"}

	if err := gconf.Conf.Parse(cliArgs...); err != nil {
		fmt.Println(err)
		return
	}

	// Get the configuration by the struct.
	fmt.Printf("------ Struct ------\n")
	fmt.Printf("Addr: %s\n", config.Addr)
	fmt.Printf("File: %s\n", config.File)
	fmt.Printf("Level: %s\n", config.Level)
	fmt.Printf("Ignore: %v\n", config.Ignore)
	fmt.Printf("DB1.MySQL.Conn: %s\n", config.DB1.MySQL.Conn)
	fmt.Printf("DB1.MySQL.MaxConnNum: %d\n", config.DB1.MySQL.MaxConnNum)
	fmt.Printf("DB2.MySQL.Conn: %s\n", config.DB2.MySQL.Conn)
	fmt.Printf("DB2.MySQL.MaxConnNum: %d\n", config.DB2.MySQL.MaxConnNum)
	fmt.Printf("DB3.MySQL.Conn: %s\n", config.DB3.MySQL.Conn)
	fmt.Printf("DB3.MySQL.MaxConnNum: %d\n", config.DB3.MySQL.MaxConnNum)
	fmt.Printf("DB4.MySQL.Conn: %s\n", config.DB4.MySQL.Conn)
	fmt.Printf("DB4.MySQL.MaxConnNum: %d\n", config.DB4.MySQL.MaxConnNum)
	fmt.Printf("DB5.DB1.MySQL.Conn: %s\n", config.DB5.DB1.MySQL.Conn)
	fmt.Printf("DB5.DB1.MySQL.MaxConnNum: %d\n", config.DB5.DB1.MySQL.MaxConnNum)
	fmt.Printf("DB5.DB2.MySQL.Conn: %s\n", config.DB5.DB2.MySQL.Conn)
	fmt.Printf("DB5.DB2.MySQL.MaxConnNum: %d\n", config.DB5.DB2.MySQL.MaxConnNum)

	// Get the configuration by the Config.
	fmt.Printf("\n------ Config ------\n")
	fmt.Printf("Addr: %s\n", gconf.Conf.String("addr"))
	fmt.Printf("File: %s\n", gconf.Conf.Group("log").String("file"))
	fmt.Printf("Level: %s\n", gconf.Conf.Group("log").String("level"))
	fmt.Printf("Ignore: %v\n", gconf.Conf.BoolD("ignore", true))
	fmt.Printf("DB1.MySQL.Conn: %s\n", gconf.Conf.Group("db1").Group("mysql").String("conn"))
	fmt.Printf("DB1.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db1.mysql").Int("maxconn"))
	fmt.Printf("DB2.MySQL.Conn: %s\n", gconf.Conf.Group("db02.mysql").String("conn"))
	fmt.Printf("DB2.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db02").Group("mysql").Int("maxconn"))
	fmt.Printf("DB3.MySQL.Conn: %s\n", gconf.Conf.Group("db03.mysql").String("conn"))
	fmt.Printf("DB3.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db03").Group("mysql").Int("maxconn"))
	fmt.Printf("DB4.MySQL.Conn: %s\n", gconf.Conf.Group("db004").Group("mysql").String("conn"))
	fmt.Printf("DB4.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db004.mysql").Int("maxconn"))
	fmt.Printf("DB5.DB1.MySQL.Conn: %s\n", gconf.Conf.Group("db").Group("db1").Group("mysql").String("conn"))
	fmt.Printf("DB5.DB1.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db.db1").Group("mysql").Int("maxconn"))
	fmt.Printf("DB5.DB2.MySQL.Conn: %s\n", gconf.Conf.Group("db").Group("db222.mysql").String("conn"))
	fmt.Printf("DB5.DB2.MySQL.MaxConnNum: %d\n", gconf.Conf.Group("db.db222.mysql").Int("maxconn"))

	// Print the group tree for debug.
	fmt.Printf("\n------ Debug ------\n")
	gconf.Conf.PrintTree(os.Stdout)

	// Output:
	// ------ Struct ------
	// Addr: 0.0.0.0:80
	// File: /var/log/test.log
	// Level: debug
	// Ignore: false
	// DB1.MySQL.Conn: user:pass@tcp(localhost:3306)/db1
	// DB1.MySQL.MaxConnNum: 3
	// DB2.MySQL.Conn: user:pass@tcp(localhost:3306)/db2
	// DB2.MySQL.MaxConnNum: 3
	// DB3.MySQL.Conn: user:pass@tcp(localhost:3306)/db3
	// DB3.MySQL.MaxConnNum: 3
	// DB4.MySQL.Conn: user:pass@tcp(localhost:3306)/db4
	// DB4.MySQL.MaxConnNum: 3
	// DB5.DB1.MySQL.Conn: user:pass@tcp(localhost:3306)/db5-1
	// DB5.DB1.MySQL.MaxConnNum: 3
	// DB5.DB2.MySQL.Conn: user:pass@tcp(localhost:3306)/db5-2
	// DB5.DB2.MySQL.MaxConnNum: 3
	//
	// ------ Config ------
	// Addr: 0.0.0.0:80
	// File: /var/log/test.log
	// Level: debug
	// Ignore: true
	// DB1.MySQL.Conn: user:pass@tcp(localhost:3306)/db1
	// DB1.MySQL.MaxConnNum: 3
	// DB2.MySQL.Conn: user:pass@tcp(localhost:3306)/db2
	// DB2.MySQL.MaxConnNum: 3
	// DB3.MySQL.Conn: user:pass@tcp(localhost:3306)/db3
	// DB3.MySQL.MaxConnNum: 3
	// DB4.MySQL.Conn: user:pass@tcp(localhost:3306)/db4
	// DB4.MySQL.MaxConnNum: 3
	// DB5.DB1.MySQL.Conn: user:pass@tcp(localhost:3306)/db5-1
	// DB5.DB1.MySQL.MaxConnNum: 3
	// DB5.DB2.MySQL.Conn: user:pass@tcp(localhost:3306)/db5-2
	// DB5.DB2.MySQL.MaxConnNum: 3
	//
	// ------ Debug ------
	// [DEFAULT]
	// |--- addr*
	// |--- config-file*
	// |-->[db]
	// |   |-->[db1]
	// |   |   |-->[mysql]
	// |   |   |   |--- conn*
	// |   |   |   |--- maxconn*
	// |   |-->[db222]
	// |   |   |-->[mysql]
	// |   |   |   |--- conn*
	// |   |   |   |--- maxconn*
	// |-->[db004]
	// |   |-->[mysql]
	// |   |   |--- conn
	// |   |   |--- maxconn
	// |-->[db02]
	// |   |-->[mysql]
	// |   |   |--- conn
	// |   |   |--- maxconn
	// |-->[db03]
	// |   |-->[mysql]
	// |   |   |--- conn
	// |   |   |--- maxconn
	// |-->[db1]
	// |   |-->[mysql]
	// |   |   |--- conn
	// |   |   |--- maxconn
	// |-->[log]
	// |   |--- file*
	// |   |--- level*
}
```

### Register the commands

```go
package main

import (
	"os"

	"github.com/xgfone/gconf"
)

type Command1OptGroup struct {
	Opt1 string
	Opt2 int
}

type Command1 struct {
	Opt1 int
	Opt2 string

	Group Command1OptGroup
}

type Command2OptGroup struct {
	Opt1 int
	Opt2 string
}

type Command2 struct {
	Opt1 string
	Opt2 int

	Group Command2OptGroup
}

type OptGroup struct {
	Opt1 string
	Opt2 int
}

type Config struct {
	Opt1 int
	Opt2 string

	Group OptGroup `cli:"false"`

	Cmd1 Command1 `cmd:"cmd1"`
	Cmd2 Command2 `cmd:"cmd2"`
}

func main() {
	gconf.Conf.RegisterCliStruct(new(Config))
	gconf.Conf.Parse()
	gconf.Conf.PrintTree(os.Stdout)

	// Output:
	// [DEFAULT]
	// |--- config-file*
	// |--- opt1*
	// |--- opt2*
	// |-->[group]
	// |   |--- opt1
	// |   |--- opt2
	// |-->{cmd1}
	// |   |--- opt1*
	// |   |--- opt2*
	// |   |-->[group]
	// |   |   |--- opt1*
	// |   |   |--- opt2*
	// |-->{cmd2}
	// |   |--- opt1*
	// |   |--- opt2*
	// |   |-->[group]
	// |   |   |--- opt1*
	// |   |   |--- opt2*
}
```

### Parse Cli Command

The `flag` parser cannot understand the command and the sub-command, so it will ignore them. But `cli` parser based on `github.com/urfave/cli` can understand them. For example,

```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

type OptGroup struct {
	Opt3 string
	Opt4 int
}

type Command1 struct {
	Opt5 int
	Opt6 string
}

type Command2 struct {
	Opt7 string `default:"abc"`
	Opt8 int    `cmd:"cmd3" help:"test sub-command" action:"cmd3_action"`
}

type Config struct {
	Opt1  int
	Opt2  string   `default:"hij"`
	Group OptGroup `cli:"false"`

	Cmd1 Command1 `cmd:"cmd1" help:"test cmd1" action:"cmd1_action"`
	Cmd2 Command2 `cmd:"cmd2" help:"test cmd2" action:"cmd2_action"`
}

func main() {
	conf := gconf.NewConfig("", "the cli command").AddParser(gconf.NewDefaultCliParser(true))
	// conf.SetDebug(true)

	// In order to let the help/version option work correctly, you should set
	// the actions for the main and the commands.
	conf.SetAction(func() error { // the main, that's, the non-command.
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterAction("cmd1_action", func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd1.opt5=%d\n", conf.Command("cmd1").Int("opt5"))
		fmt.Printf("cmd1.opt6=%s\n", conf.Command("cmd1").String("opt6"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterAction("cmd2_action", func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd2.opt7=%s\n", conf.Command("cmd2").String("opt7"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterAction("cmd3_action", func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd2.opt7=%s\n", conf.Command("cmd2").String("opt7"))
		fmt.Printf("cmd2.cmd3.opt8=%d\n", conf.Command("cmd2").Command("cmd3").Int("opt8"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterCliStruct(new(Config)).Parse()

	// RUN Example 1: the main
	// $ PROGRAM -h
	// NAME:
	//    PROGRAM - the cli command
	//
	// USAGE:
	//    main [global options] command [command options] [arguments...]
	//
	// VERSION:
	//    0.0.0
	//
	// COMMANDS:
	//      cmd1     test cmd1
	//      cmd2     test cmd2
	//      help, h  Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --opt1 value   (default: 0)
	//    --opt2 value   (default: "hij")
	//    --help, -h     show help
	//    --version, -v  print the version

	// RUN Example 2: the command "cmd2"
	// $ PROGRAM cmd2 -h
	// NAME:
	//    main cmd2 - test cmd2
	//
	// USAGE:
	//    main cmd2 command [command options] [arguments...]
	//
	// COMMANDS:
	//      cmd3  test sub-command
	//
	// OPTIONS:
	//    --opt7 value  (default: "abc")
	//    --help, -h    show help

	// RUN Example 3: the sub-command "cmd3"
	// $ PROGRAM cmd1 cmd3 -h
	// NAME:
	//    main cmd2 cmd3 - test sub-command
	//
	// USAGE:
	//    main cmd2 cmd3 [command options] [arguments...]
	//
	// OPTIONS:
	//    --opt8 value  test sub-command (default: 0)
}
```

Beside, you maybe build the commands by hand instead of using the struct. So the program above is equal to that below.

```go
package main

import (
	"fmt"

	"github.com/xgfone/gconf"
)

func main() {
	conf := gconf.NewConfig("", "the cli command").AddParser(gconf.NewDefaultCliParser(true))
	// conf.SetDebug(true)

	// Build the main
	conf.SetAction(func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterOpts([]gconf.Opt{
		gconf.Int("opt1", 0, ""),
		gconf.Str("opt2", "hij", ""),
	}).NewGroup("group").RegisterOpts([]gconf.Opt{
		gconf.Str("opt3", "", ""),
		gconf.Int("opt4", 0, ""),
	})

	// Build the command "cmd1"
	conf.NewCommand("cmd1", "test cmd1").SetAction(func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd1.opt5=%d\n", conf.Command("cmd1").Int("opt5"))
		fmt.Printf("cmd1.opt6=%s\n", conf.Command("cmd1").String("opt6"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterCliOpts([]gconf.Opt{
		gconf.Int("opt5", 0, ""),
		gconf.Str("opt6", "", ""),
	})

	// Build the command "cmd2"
	conf.NewCommand("cmd2", "test cmd2").SetAction(func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd2.opt7=%s\n", conf.Command("cmd2").String("opt7"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterCliOpt(gconf.Str("opt7", "abc", ""))

	// Build the sub-command "cmd3" of the command "cmd2".
	conf.NewCommand("cmd3", "test sub-command").SetAction(func() error {
		fmt.Printf("opt1=%d\n", conf.Int("opt1"))
		fmt.Printf("opt2=%s\n", conf.String("opt2"))
		fmt.Printf("cmd2.opt7=%s\n", conf.Command("cmd2").String("opt7"))
		fmt.Printf("cmd2.cmd3.opt8=%d\n", conf.Command("cmd2").Command("cmd3").Int("opt8"))
		fmt.Printf("args=%v\n", conf.CliArgs())
		return nil
	}).RegisterCliOpt(gconf.Int("opt8", 0, ""))

	// Parse and run the command.
	conf.Parse()
}
```
