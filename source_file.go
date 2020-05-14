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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigFileOpt is the default option for the configuration file.
var ConfigFileOpt = StrOpt("config-file", "the config file path.")

// NewFileSource returns a new source that the data is read from the file
// named filename.
//
// Notice: it will identify the format by the filename extension automatically.
// If no filename extension, it will use defaulFormat, which is "ini" by default.
func NewFileSource(filename string, defaultFormat ...string) Source {
	id := fmt.Sprintf("file:%s", filename)
	format := strings.Trim(filepath.Ext(filename), ".")
	if format == "" {
		if len(defaultFormat) > 0 && defaultFormat[0] != "" {
			format = defaultFormat[0]
		} else {
			format = "ini"
		}
	}
	return fileSource{id: id, filepath: filename, format: format}
}

type fileSource struct {
	id       string
	format   string
	filepath string
}

func (f fileSource) Read() (DataSet, error) {
	file, err := os.Open(f.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return DataSet{Source: f.id, Format: f.format}, nil
		}
		return DataSet{Source: f.id, Format: f.format}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return DataSet{Source: f.id, Format: f.format}, err
	}

	stat, err := file.Stat()
	if err != nil {
		return DataSet{Source: f.id, Format: f.format, Data: data}, err
	}

	ds := DataSet{
		Data:      data,
		Format:    f.format,
		Source:    f.id,
		Timestamp: stat.ModTime(),
	}
	ds.Checksum = "md5:" + ds.Md5()

	return ds, nil
}

func (f fileSource) Watch(load func(DataSet, error) bool, exit <-chan struct{}) {
	go f.watchfile(load, exit)
}

func (f fileSource) watchfile(load func(DataSet, error) bool, exit <-chan struct{}) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		load(DataSet{Source: f.id, Format: f.format}, err)
		return
	}
	defer fw.Close()

	var add bool
	if fw.Add(f.filepath) == nil {
		add = true
	}

	interval := time.Second * 10
	for {
		time.Sleep(interval)
		if _, err := os.Stat(f.filepath); err != nil {
			if os.IsNotExist(err) {
				if add {
					fw.Remove(f.filepath)
					add = false
				}
			} else {
				load(DataSet{Source: f.id, Format: f.format}, err)
			}
			continue
		}

		if !add {
			if err := fw.Add(f.filepath); err != nil {
				load(DataSet{Source: f.id, Format: f.format}, err)
			} else {
				add = true
				load(f.Read())
			}
			continue
		}

		select {
		case <-exit:
			return
		case event, ok := <-fw.Events:
			if !ok {
				return
			}

			// BUG: it will be triggered twice continuously by fsnotify on Windows.
			if event.Op&fsnotify.Write == fsnotify.Write {
				load(f.Read())
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				add = false
			}
		case err, ok := <-fw.Errors:
			if !ok {
				return
			}
			load(DataSet{Source: f.id, Format: f.format}, err)
		}
	}

}
