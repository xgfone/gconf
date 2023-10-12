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
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConfigFileOpt is the default option for the configuration file.
var ConfigFileOpt = StrOpt("config-file", "the config file path.")

func init() { Conf.RegisterOpts(ConfigFileOpt) }

// NewFileSource returns a new source that the data is read from the file
// named filename.
//
// The file source can watch the change of the given file.
// And it will identify the format by the filename extension automatically.
// If no filename extension, it will use defaulFormat, which is "ini" by default.
func NewFileSource(filename string, defaultFormat ...string) Source {
	format := strings.Trim(filepath.Ext(filename), ".")
	if format == "" {
		if len(defaultFormat) > 0 && defaultFormat[0] != "" {
			format = defaultFormat[0]
		} else {
			format = "ini"
		}
	}

	id := fmt.Sprintf("file:%s", filename)
	return fileSource{id: id, filepath: filename, format: format, timeout: time.Second * 10}
}

type fileSource struct {
	id       string
	format   string
	filepath string
	timeout  time.Duration
}

func (f fileSource) String() string { return f.id }

func (f fileSource) Read() (DataSet, error) {
	file, err := os.Open(f.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return DataSet{Source: f.id, Format: f.format}, nil
		}
		return DataSet{Source: f.id, Format: f.format}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
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

func (f fileSource) Watch(exit <-chan struct{}, load func(DataSet, error) bool) {
	f.watch(exit, load)
}

func (f fileSource) watch(exit <-chan struct{}, load func(DataSet, error) bool) {
	lastsize, lasttime, _ := getfileinfo(f.filepath)

	ticker := time.NewTicker(f.timeout)
	defer ticker.Stop()

	for {
		select {
		case <-exit:
			return

		case <-ticker.C:
			if size, time, err := getfileinfo(f.filepath); err != nil {
				if !os.IsNotExist(err) {
					load(DataSet{Source: f.id, Format: f.format}, err)
				}
			} else if size != lastsize || time != lasttime {
				load(f.Read())
				lastsize = size
				lasttime = time
			}
		}
	}
}

func getfileinfo(filepath string) (size, time int64, err error) {
	fi, err := os.Stat(filepath)
	if err == nil {
		time = fi.ModTime().Unix()
		size = fi.Size()
	}
	return
}
