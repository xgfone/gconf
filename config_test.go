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
	"fmt"
	"os"
	"testing"
)

func ExampleConfig_Observe() {
	conf := New()
	conf.RegisterCliOpt(Int("watchval", 123, "test watch int value"))
	conf.NewGroup("test").RegisterOpt(Str("watchval", "abc", "test watch str value"))
	conf.Observe(func(group, opt string, old, new interface{}) {
		fmt.Printf("group=%s, option=%s, old=%v, new=%v\n", group, opt, old, new)
	})

	// Start the config
	conf.Parse()

	// Set the option vlaue during the program is running.
	conf.UpdateOptValue("", "watchval", 456)
	conf.UpdateOptValue("test", "watchval", "123")

	// Output:
	// group=DEFAULT, option=watchval, old=<nil>, new=123
	// group=test, option=watchval, old=<nil>, new=abc
	// group=DEFAULT, option=watchval, old=123, new=456
	// group=test, option=watchval, old=abc, new=123
}

func ExampleNewEnvVarParser() {
	// Simulate the environment variable.
	os.Setenv("TEST_VAR1", "abc")
	os.Setenv("TEST_GROUP1_GROUP2_VAR2", "123")

	conf := New().AddParser(NewEnvVarParser(10, "test"))

	opt1 := Str("var1", "", "the environment var 1")
	opt2 := Int("var2", 0, "the environment var 2")

	conf.RegisterOpt(opt1)
	conf.NewGroup("group1").NewGroup("group2").RegisterOpt(opt2)

	if err := conf.Parse(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("var1=%s\n", conf.String("var1"))
	fmt.Printf("var2=%d\n", conf.Group("group1.group2").Int("var2"))
	fmt.Printf("var2=%d\n", conf.Group("group1").Group("group2").Int("var2"))

	// Output:
	// var1=abc
	// var2=123
	// var2=123
}

func ExampleConfig() {
	cliOpts1 := []Opt{
		StrOpt("", "required", "", "required").SetValidators(NewStrLenValidator(1, 10)),
		BoolOpt("", "yes", true, "test bool option"),
	}

	cliOpts2 := []Opt{
		BoolOpt("", "no", false, "test bool option"),
		StrOpt("", "optional", "optional", "optional"),
	}

	opts := []Opt{
		Str("opt", "", "test opt"),
	}

	conf := New().AddParser(NewFlagCliParser(nil, true))
	conf.RegisterCliOpts(cliOpts1)
	conf.NewGroup("cli").RegisterCliOpts(cliOpts2)
	conf.NewGroup("group1").NewGroup("group2").RegisterCliOpts(opts)

	// For Test
	cliArgs := []string{
		"--cli.no=0",
		"--required", "required",
		"--group1.group2.opt", "option",
	}

	if err := conf.Parse(cliArgs...); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("yes=%v\n", conf.Bool("yes"))
	fmt.Printf("required=%v\n", conf.StringD("required", "abc"))

	fmt.Println()

	fmt.Printf("cli.no=%v\n", conf.Group("cli").Bool("no"))
	fmt.Printf("cli.optional=%v\n", conf.Group("cli").String("optional"))

	fmt.Println()

	fmt.Printf("group1.group2.opt=%v\n", conf.Group("group1.group2").StringD("opt", "opt"))
	fmt.Printf("group1.group2.opt=%v\n", conf.Group("group1").Group("group2").StringD("opt", "opt"))

	// Output:
	// yes=true
	// required=required
	//
	// cli.no=false
	// cli.optional=optional
	//
	// group1.group2.opt=option
	// group1.group2.opt=option
}

func ExampleConfig_RegisterStruct() {
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

	// Set the debug to output the process that handles the configuration.
	// Conf.SetDebug(true)

	var config Config
	Conf.AddParser(NewEnvVarParser(10, "")).RegisterCliStruct(&config)

	// Only for test
	os.Setenv("DB1_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db1")
	os.Setenv("DB02_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db2")
	os.Setenv("DB03_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db3")
	os.Setenv("DB004_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db4")
	os.Setenv("DB_DB1_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db5-1")
	os.Setenv("DB_DB222_MYSQL_CONN", "user:pass@tcp(localhost:3306)/db5-2")
	cliArgs := []string{"--addr", "0.0.0.0:80", "--log.file", "/var/log/test.log"}

	if err := Conf.Parse(cliArgs...); err != nil {
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
	fmt.Printf("Addr: %s\n", Conf.String("addr"))
	fmt.Printf("File: %s\n", Conf.Group("log").String("file"))
	fmt.Printf("Level: %s\n", Conf.Group("log").String("level"))
	fmt.Printf("Ignore: %v\n", Conf.BoolD("ignore", true))
	fmt.Printf("DB1.MySQL.Conn: %s\n", Conf.Group("db1").Group("mysql").String("conn"))
	fmt.Printf("DB1.MySQL.MaxConnNum: %d\n", Conf.Group("db1.mysql").Int("maxconn"))
	fmt.Printf("DB2.MySQL.Conn: %s\n", Conf.Group("db02.mysql").String("conn"))
	fmt.Printf("DB2.MySQL.MaxConnNum: %d\n", Conf.Group("db02").Group("mysql").Int("maxconn"))
	fmt.Printf("DB3.MySQL.Conn: %s\n", Conf.Group("db03.mysql").String("conn"))
	fmt.Printf("DB3.MySQL.MaxConnNum: %d\n", Conf.Group("db03").Group("mysql").Int("maxconn"))
	fmt.Printf("DB4.MySQL.Conn: %s\n", Conf.Group("db004").Group("mysql").String("conn"))
	fmt.Printf("DB4.MySQL.MaxConnNum: %d\n", Conf.Group("db004.mysql").Int("maxconn"))
	fmt.Printf("DB5.DB1.MySQL.Conn: %s\n", Conf.Group("db").Group("db1").Group("mysql").String("conn"))
	fmt.Printf("DB5.DB1.MySQL.MaxConnNum: %d\n", Conf.Group("db.db1").Group("mysql").Int("maxconn"))
	fmt.Printf("DB5.DB2.MySQL.Conn: %s\n", Conf.Group("db").Group("db222.mysql").String("conn"))
	fmt.Printf("DB5.DB2.MySQL.MaxConnNum: %d\n", Conf.Group("db.db222.mysql").Int("maxconn"))

	// Print the group tree for debug.
	fmt.Printf("\n------ Debug ------\n")
	Conf.PrintTree(os.Stdout)

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

func ExampleCommand() {
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

	var c Config
	conf := NewDefault(nil)
	conf.RegisterCliStruct(&c)

	// conf.Parse() // We only test.

	conf.PrintTree(os.Stdout)

	// Output:
	// [DEFAULT]
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

func ExampleNewCliParser() {
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

	conf := NewConfig("", "the cli command").AddParser(NewDefaultCliParser(true))

	// conf.SetDebug(true)
	conf.SetAction(func() error {
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
	}).RegisterCliStruct(new(Config)).Parse(
		"--opt1", "123", // Global Options
		"cmd2", "--opt7=xyz", // The command "cmd2" and its options
		"cmd3", "-opt8", "456", // The sub-command "cmd3" and its options
		"arg1", "arg2", "arg3", // The rest arguments
	)

	// Output:
	// opt1=123
	// opt2=hij
	// cmd2.opt7=xyz
	// cmd2.cmd3.opt8=456
	// args=[arg1 arg2 arg3]
}

func TestNewSimpleIniParser(t *testing.T) {
	conf := NewDefault(nil)
	conf.RegisterCliOpt(Str("opt1", "", ""))
	conf.RegisterCliOpt(Str("opt2", "abc", ""))
	conf.RegisterCliOpt(Str("opt3", "opq", ""))

	// Write the test config data into the file.
	filename := "test_new_simple_ini_parser.ini"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
		return
	}
	file.WriteString("[DEFAULT]\nopt1=123\nopt2=456\nopt3=789")
	file.Close()
	defer os.Remove(filename)

	if err = conf.Parse("--config-file", filename, "--opt3=xyz"); err != nil {
		t.Error(err)
		return
	}

	if opt1 := conf.String("opt1"); opt1 != "123" {
		t.Error(opt1)
	}
	if opt2 := conf.String("opt2"); opt2 != "456" {
		t.Error(opt2)
	}
	if opt3 := conf.String("opt3"); opt3 != "xyz" {
		t.Error(opt3)
	}
}
