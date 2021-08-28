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
	"strings"
	"time"
)

// NewEnvSource returns a new Source based on the environment variables,
// which reads the configuration from the environment variables.
//
// If giving the prefix, it only uses the environment variable name
// matching the given prefix, then removes the prefix and the rest is used
// as the option name.
//
// Notice: It will convert all the underlines("_") to the dots(".").
func NewEnvSource(prefix string) Source {
	if prefix != "" {
		if prefix = strings.Trim(prefix, "_"); prefix != "" {
			prefix += "_"
		}
	}
	return envSource{prefix: strings.ToLower(prefix)}
}

type envSource struct{ prefix string }

func (e envSource) String() string { return "env" }

func (e envSource) Watch(<-chan struct{}, func(DataSet, error) bool) {}

func (e envSource) Read() (DataSet, error) {
	vs := make(map[string]string, 32)
	for _, env := range os.Environ() {
		index := strings.IndexByte(env, '=')
		if index == -1 {
			continue
		}

		value := strings.TrimSpace(env[index+1:])
		if value == "" {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(env[:index]))
		if e.prefix != "" {
			if !strings.HasPrefix(key, e.prefix) {
				continue
			}
			key = strings.TrimPrefix(key, e.prefix)
		}

		key = strings.Replace(strings.Trim(key, "_"), "_", ".", -1)
		if key != "" {
			vs[key] = value
		}
	}

	data, err := json.Marshal(vs)
	if err != nil {
		return DataSet{Format: "json", Source: e.String()}, err
	}

	ds := DataSet{
		Data:      data,
		Format:    "json",
		Source:    e.String(),
		Timestamp: time.Now(),
	}
	ds.Checksum = "md5:" + ds.Md5()
	return ds, nil
}
