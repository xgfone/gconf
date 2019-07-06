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
	"os"
	"strings"
	"time"
)

// NewEnvSource returns a new env Source to read the configuration
// from the environment variables.
//
// If giving the prefix, it will remove the prefix from the env key,
// then the rest is used for the option key.
//
// Notice: It will convert all the underlines("_") to the dots(".").
func NewEnvSource(prefix ...string) Source {
	var _prefix string
	if len(prefix) > 0 && prefix[0] != "" {
		if _prefix = strings.Trim(prefix[0], "_"); _prefix != "" {
			_prefix += "_"
		}
	}
	return envSource{prefix: strings.ToLower(_prefix)}
}

type envSource struct {
	prefix string
}

func (e envSource) Watch(load func(DataSet, error), exit <-chan struct{}) {}

func (e envSource) Read() (DataSet, error) {
	vs := make(map[string]string, 32)
	for _, env := range os.Environ() {
		index := strings.IndexByte(env, '=')
		if index == -1 {
			continue
		}

		value := strings.TrimSpace(env[index+1:])
		key := strings.ToLower(strings.TrimSpace(env[:index]))
		key = strings.Replace(strings.TrimPrefix(key, e.prefix), "_", ".", -1)
		if key = strings.Trim(key, "."); key != "" && value != "" {
			vs[key] = value
		}
	}

	data, err := json.Marshal(vs)
	if err != nil {
		return DataSet{Source: "env", Format: "json"}, err
	}
	ds := DataSet{Data: data, Format: "json", Source: "env", Timestamp: time.Now()}
	ds.Checksum = fmt.Sprintf("md5:%s", ds.Md5())
	return ds, nil
}
