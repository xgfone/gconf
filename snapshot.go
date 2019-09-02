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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// LoadBackupFile loads configuration data from the backup file if exists,
// then watches the change of the options and write them into the file.
//
// So you can use it as the local cache.
func (c *Config) LoadBackupFile(filename string) error {
	var data []byte
	var ms map[string]interface{}

	if _, err := os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if data, err = ioutil.ReadFile(filename); err != nil {
		return err
	} else if err = json.Unmarshal(data, &ms); err != nil {
		return err
	}

	c.updateFlattedMap(ms, false)
	go c.writeSnapshotIntoFile(bytesToMd5(data), filename)

	return nil
}

func (c *Config) writeSnapshotIntoFile(lastChecksum, filename string) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case _, ok := <-c.exit:
			if !ok {
				return
			}
		case <-ticker.C:
			data, err := c.snap.MarshalJSON()
			if err != nil {
				c.handleError(fmt.Errorf("[Config] snapshot marshal json: %s", err.Error()))
			}

			if checksum := bytesToMd5(data); checksum != lastChecksum {
				if err = ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
					c.handleError(fmt.Errorf("[Config] snapshot write file[%s]: %s", filename, err.Error()))
				} else {
					lastChecksum = checksum
					debugf("[Config] Write snapshot into file '%s'", filename)
				}
			}
		}
	}
}

// Snapshot returns the snapshot of the whole configuration options.
//
// Notice: the key includes the group name and the option name, for instance,
//
//   map[string]interface{} {
//       "opt1": "value1",
//       "opt2": "value2",
//       "group1.opt3": "value3",
//       "group1.group2.opt4": "value4",
//       // ...
//   }
func (c *Config) Snapshot() map[string]interface{} {
	return c.snap.ToMap()
}

func newSnapshot(c *Config) *snapshot {
	snap := &snapshot{conf: c, maps: make(map[string]interface{}, 64)}
	c.ObserveRegister(snap.RegisterObserver)
	c.Observe(snap.ChangeObserver)
	return snap
}

type snapshot struct {
	conf *Config
	lock sync.RWMutex
	maps map[string]interface{}
}

func (s *snapshot) ToMap() map[string]interface{} {
	s.lock.RLock()
	dst := make(map[string]interface{}, len(s.maps)*2)
	for key, value := range s.maps {
		dst[key] = value
	}
	s.lock.RUnlock()
	return dst
}

func (s *snapshot) RegisterObserver(group string, opts []Opt) {
	s.lock.Lock()
	for _, opt := range opts {
		key := s.conf.fixOptName(opt.Name)
		if group != "" {
			key = fmt.Sprintf("%s%s%s", group, s.conf.gsep, key)
		}
		s.maps[key] = opt.Default
	}
	s.lock.Unlock()
}

func (s *snapshot) ChangeObserver(group, opt string, old, new interface{}) {
	key := opt
	if group != "" {
		key = fmt.Sprintf("%s%s%s", group, s.conf.gsep, opt)
	}

	s.lock.Lock()
	s.maps[key] = new
	s.lock.Unlock()
}

func (s *snapshot) MarshalJSON() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return json.Marshal(s.maps)
}

func (s *snapshot) UnmarshalJSON(data []byte) error {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return json.Unmarshal(data, &s.maps)
}
