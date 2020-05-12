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
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
)

func ExampleNewCliSource() {
	conf := New()
	conf.RegisterOpt(StrOpt("opt1", "").D("abc"))
	conf.NewGroup("cmd1").RegisterOpt(IntOpt("opt2", ""))
	conf.NewGroup("cmd1").NewGroup("cmd2").RegisterOpt(IntOpt("opt3", ""))

	app := cli.NewApp()
	app.Flags = []cli.Flag{&cli.StringFlag{Name: "opt1"}}
	app.Commands = []*cli.Command{
		{
			Name:  "cmd1",
			Flags: []cli.Flag{&cli.IntFlag{Name: "opt2"}},
			Subcommands: []*cli.Command{
				{
					Name:  "cmd2",
					Flags: []cli.Flag{&cli.IntFlag{Name: "opt3"}},
					Action: func(ctx *cli.Context) error {
						ctxs := ctx.Lineage()
						conf.LoadSource(NewCliSource(ctxs[0], "cmd1", "cmd2"))
						conf.LoadSource(NewCliSource(ctxs[1], "cmd1"))
						conf.LoadSource(NewCliSource(ctxs[2]))

						fmt.Println(conf.GetString("opt1"))
						fmt.Println(conf.Group("cmd1").GetInt("opt2"))
						fmt.Println(conf.Group("cmd1.cmd2").GetInt("opt3"))

						return nil
					},
				},
			},
		},
	}

	// For Test
	app.Run([]string{"app", "--opt1=xyz", "cmd1", "--opt2=123", "cmd2", "--opt3=456"})

	// Output:
	// xyz
	// 123
	// 456
}

func TestNewEnvSource(t *testing.T) {
	os.Setenv("OPT1", "123")
	os.Setenv("GROUP1_OPT2", "1")
	os.Setenv("GROUP1_GROUP2_OPT3", "456")
	os.Setenv("TEST_OPT1", "456")
	os.Setenv("TEST_GROUP1_OPT2", "0")
	os.Setenv("TEST_GROUP1_GROUP2_OPT3", "789")
	os.Setenv("ABC", "xyz")

	conf := New()
	conf.RegisterOpt(IntOpt("opt1", ""))
	conf.NewGroup("group1").RegisterOpt(BoolOpt("opt2", ""))
	conf.Group("group1").NewGroup("group2").RegisterOpt(Float64Opt("opt3", ""))

	conf.LoadSource(NewEnvSource())
	if v := conf.GetInt("opt1"); v != 123 {
		t.Error(v)
	} else if v := conf.Group("group1").GetBool("opt2"); !v {
		t.Fail()
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 456 {
		t.Error(v)
	}

	conf.LoadSource(NewEnvSource("test"), true)
	if v := conf.GetInt("opt1"); v != 456 {
		t.Error(v)
	} else if v := conf.Group("group1").GetBool("opt2"); v {
		t.Fail()
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 789 {
		t.Error(v)
	}
}

func TestNewFileSource_INI(t *testing.T) {
	// Prepare the ini file
	filename := "_test_ini_file_source_.conf"
	if file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		t.Error(err)
	} else {
		file.Write([]byte(`
		[DEFAULT]
		opt1 = 1
		opt2 = 0
		opt3 = 1
		[group1]
		opt1 = 2
		opt2 = 1
		opt3 = 2
		[group1.group2]
		opt1 = 3
		opt2 = 0
		opt3 = 3
		`))
		file.Close()
	}
	defer os.Remove(filename)

	// Load the config
	conf := New()
	conf.RegisterOpt(IntOpt("opt1", ""))
	conf.NewGroup("group1").RegisterOpt(BoolOpt("opt2", ""))
	conf.Group("group1").NewGroup("group2").RegisterOpt(Float64Opt("opt3", ""))
	conf.LoadSource(NewFileSource(filename))

	// Check the config
	if v := conf.GetInt("opt1"); v != 1 {
		t.Error(v)
	} else if v := conf.Group("group1").GetBool("opt2"); !v {
		t.Fail()
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 3 {
		t.Error(v)
	}
}

func TestNewFileSource_JSON(t *testing.T) {
	// Prepare the json file
	filename := "_test_json_file_source_.json"
	if file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		t.Error(err)
	} else {
		file.Write([]byte(`{
			"opt1": 1,
			"opt2": false,
			"opt3": 1,
			"group1": {
				"opt1": 2,
				"opt2": true,
				"opt3": 2,
				"group2": {
					"opt1": 3,
					"opt2": false,
					"opt3": 3
				}
			}
		}`))
		file.Close()
	}
	defer os.Remove(filename)

	// Load the config
	conf := New()
	conf.RegisterOpt(IntOpt("opt1", ""))
	conf.NewGroup("group1").RegisterOpt(BoolOpt("opt2", ""))
	conf.Group("group1").NewGroup("group2").RegisterOpt(Float64Opt("opt3", ""))
	conf.LoadSource(NewFileSource(filename))

	// Check the config
	if v := conf.GetInt("opt1"); v != 1 {
		t.Error(v)
	} else if v := conf.Group("group1").GetBool("opt2"); !v {
		t.Fail()
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 3 {
		t.Error(v)
	}
}

func TestNewURLSource(t *testing.T) {
	first := true

	// Start the http server
	go http.ListenAndServe("127.0.0.1:12345", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if first {
			first = false
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(`{"opt": 123}`))
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(`{"opt": 456}`))
		}
	}))
	time.Sleep(time.Millisecond * 50) // Wait that the http server finishes to start.

	conf := New()
	conf.RegisterOpt(IntOpt("opt", ""))
	conf.LoadSource(NewURLSource("http://127.0.0.1:12345/", time.Millisecond*100))

	if v := conf.GetInt("opt"); v != 123 {
		t.Error(v)
	} else {
		time.Sleep(time.Millisecond * 200)
		if v := conf.GetInt("opt"); v != 456 {
			t.Error(v)
		}
	}
}
