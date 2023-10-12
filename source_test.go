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
	"net/http"
	"os"
	"testing"
	"time"
)

const testfileflag = os.O_APPEND | os.O_CREATE | os.O_WRONLY

func TestNewEnvSource(t *testing.T) {
	os.Setenv("ABC", "xyz")
	os.Setenv("OPT1", "111")
	os.Setenv("GROUP1_OPT2", "abc")
	os.Setenv("GROUP1_GROUP2_OPT3", "222")

	os.Setenv("TEST_ABC", "xyz")
	os.Setenv("TEST_OPT1", "333")
	os.Setenv("TEST_GROUP1_OPT2", "efg")
	os.Setenv("TEST_GROUP1_GROUP2_OPT3", "444")

	conf := New()
	conf.RegisterOpts(IntOpt("opt1", ""))
	conf.Group("group1").RegisterOpts(StrOpt("opt2", ""))
	conf.Group("group1").Group("group2").RegisterOpts(Float64Opt("opt3", ""))

	_ = conf.LoadSource(NewEnvSource(""), true)
	if v := conf.GetInt("opt1"); v != 111 {
		t.Errorf("expect '%d', but got '%d'", 111, v)
	} else if v := conf.Group("group1").GetString("opt2"); v != "abc" {
		t.Errorf("expect '%s', but got '%s'", "abc", v)
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 222 {
		t.Errorf("expect '%d', but got '%f'", 222, v)
	}

	_ = conf.LoadSource(NewEnvSource("test"), true)
	if v := conf.GetInt("opt1"); v != 333 {
		t.Errorf("expect '%d', but got '%d'", 333, v)
	} else if v := conf.Group("group1").GetString("opt2"); v != "efg" {
		t.Errorf("expect '%s', but got '%s'", "efg", v)
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 444 {
		t.Errorf("expect '%d', but got '%f'", 444, v)
	}
}

func TestNewFileSource_INI(t *testing.T) {
	// Prepare the ini file
	filename := "_test_ini_file_source_.ini"
	file, err := os.OpenFile(filename, testfileflag, os.ModePerm)
	if err != nil {
		t.Error(err)
	} else {
		_, _ = file.Write([]byte(`
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
	conf.RegisterOpts(IntOpt("opt1", ""))
	conf.Group("group1").RegisterOpts(BoolOpt("opt2", ""))
	conf.Group("group1.group2").RegisterOpts(Float64Opt("opt3", ""))
	_ = conf.LoadSource(NewFileSource(filename))

	// Check the config
	if v := conf.GetInt("opt1"); v != 1 {
		t.Error(v)
	} else if v := conf.GetBool("group1.opt2"); !v {
		t.Fail()
	} else if v := conf.GetFloat64("group1.group2.opt3"); v != 3 {
		t.Error(v)
	}
}

func TestNewFileSource_JSON(t *testing.T) {
	// Prepare the json file
	filename := "_test_json_file_source_.json"
	file, err := os.OpenFile(filename, testfileflag, os.ModePerm)
	if err != nil {
		t.Error(err)
	} else {
		_, _ = file.Write([]byte(`{
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
	conf.RegisterOpts(IntOpt("opt1", ""))
	conf.Group("group1").RegisterOpts(BoolOpt("opt2", ""))
	conf.Group("group1").Group("group2").RegisterOpts(Float64Opt("opt3", ""))
	_ = conf.LoadSource(NewFileSource(filename))

	// Check the config
	if v := conf.GetInt("opt1"); v != 1 {
		t.Error(v)
	} else if v := conf.Group("group1").GetBool("opt2"); !v {
		t.Fail()
	} else if v := conf.Group("group1.group2").GetFloat64("opt3"); v != 3 {
		t.Error(v)
	}
}

func TestFileSourceWatch(t *testing.T) {
	// Prepare the json file
	filename := "_test_file_source_watch_.json"
	defer os.Remove(filename)

	source := NewFileSource(filename).(fileSource)
	source.timeout = time.Second * 2

	exit := make(chan struct{})
	go func() {
		time.Sleep(time.Second)
		file, err := os.OpenFile(filename, testfileflag, os.ModePerm)
		if err != nil {
			t.Error(err)
		} else {
			_, _ = file.Write([]byte(`{"opt": 1}`))
			file.Close()
		}
		time.Sleep(time.Second * 2)
		close(exit)
	}()

	var data string
	start := time.Now()
	source.Watch(exit, func(ds DataSet, err error) bool {
		if err != nil {
			t.Error(err)
		} else if data == "" {
			data = string(ds.Data)
		} else {
			t.Fail()
		}
		return true
	})

	if cost := time.Since(start); cost < time.Second*3 || cost > time.Second*4 {
		t.Errorf("not wait for 3~4s")
	}

	expect := `{"opt": 1}`
	if data != expect {
		t.Errorf("expect '%s', but got '%s'", expect, data)
	}
}

func TestNewURLSource(t *testing.T) {
	first := true

	// Start the http server
	go func() {
		_ = http.ListenAndServe("127.0.0.1:12345", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if first {
				first = false
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				_, _ = w.Write([]byte(`{"opt": 123}`))
			} else {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				_, _ = w.Write([]byte(`{"opt": 456}`))
			}
		}))
	}()

	time.Sleep(time.Millisecond * 50) // Wait that the http server finishes to start.

	conf := New()
	conf.RegisterOpts(IntOpt("opt", ""))
	_ = conf.LoadAndWatchSource(NewURLSource("http://127.0.0.1:12345/", time.Millisecond*100))
	defer conf.Stop()

	if v := conf.GetInt("opt"); v != 123 {
		t.Error(v)
	} else {
		time.Sleep(time.Millisecond * 400)
		if v := conf.GetInt("opt"); v != 456 {
			t.Error(v)
		}
	}
}
