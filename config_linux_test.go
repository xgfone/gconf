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
	"os/exec"
	"time"
)

func ExampleConfig_SetHotReload() {
	filename := "test_set_hot_reload.ini"

	// The flag and cli parser will ignore the hot-reloading automatically.
	conf := NewDefault(nil)
	conf.SetHotReload(conf.Parsers()...)
	conf.RegisterOpt(Str("reload_opt", "abc", "test reload"))
	conf.Parse("--config-file", filename)

	time.Sleep(time.Millisecond * 10)
	fmt.Println(conf.String("reload_opt"))

	// Only for test: Write the test config data into the file.
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.WriteString("[DEFAULT]\nreload_opt=xyz")
	file.Close()
	defer os.Remove(filename)

	// Send the signal SIGHUP.
	cmd := exec.Command("kill", "-HUP", fmt.Sprintf("%d", os.Getpid()))
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(time.Millisecond * 10)
		fmt.Println(conf.String("reload_opt"))
	}

	// Output:
	// abc
	// xyz
}
