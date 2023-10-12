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
	"encoding/json"
	"os"
	"sync/atomic"
	"time"
)

// LoadBackupFile loads configuration data from the backup file if exists,
// then watches the change of the options and write them into the file.
// So you can use it as the local cache.
func (c *Config) LoadBackupFile(filename string) (err error) {
	if filename == "" {
		panic("the backup filename must not be empty")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			c.errorf("fail to read the backup file '%s': %s", filename, err)
			return
		}
	}

	if len(data) > 0 {
		ms := make(map[string]interface{}, 32)
		if err = json.Unmarshal(data, &ms); err != nil {
			c.errorf("the backup file '%s' format is error: %s", filename, err)
			return
		} else if err = c.LoadMap(ms); err != nil {
			return
		}
	}

	go c.writeSnapshotIntoFile(filename)
	return
}

func (c *Config) writeSnapshotIntoFile(filename string) {
	var lastgen uint64
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case _, ok := <-c.exit:
			if !ok {
				return
			}
		case <-ticker.C:
			if gen := atomic.LoadUint64(&c.gen); gen <= lastgen {
				continue
			}

			gen, snaps := c.Snapshot()
			if gen <= lastgen || len(snaps) == 0 {
				continue
			}

			data, err := json.Marshal(snaps)
			if err != nil {
				c.errorf("fail to marshal snapshot as json: %s", err)
				continue
			}

			if err := os.WriteFile(filename, data, os.ModePerm); err != nil {
				c.errorf("cannot write snapshot into file '%s': %s", filename, err)
			} else {
				lastgen = gen
			}
		}
	}
}

// Snapshot returns the snapshot of all the options and its generation
// which will increase with 1 each time any option value is changed.
//
// For example,
//
//	map[string]interface{} {
//	    "opt1": "value1",
//	    "opt2": "value2",
//	    "group1.opt3": "value3",
//	    "group1.group2.opt4": "value4",
//	    // ...
//	}
func (c *Config) Snapshot() (generation uint64, snap map[string]interface{}) {
	generation = atomic.LoadUint64(&c.gen)
	snap = make(map[string]interface{}, len(c.options))
	for name, opt := range c.options {
		if v := opt.GetValue(); v != nil {
			snap[name] = v
		}
	}
	return
}
