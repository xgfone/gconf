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
// If no filename extension, it will use defaulFormat, or panic.
func NewFileSource(filename string, defaultFormat ...string) Source {
	id := fmt.Sprintf("file:%s", filename)
	format := strings.Trim(filepath.Ext(filename), ".")
	if format == "" {
		if len(defaultFormat) == 0 || defaultFormat[0] == "" {
			panic(fmt.Errorf("missing the file format for '%s'", filename))
		}
		format = defaultFormat[0]
	}
	return fileSource{id: id, filepath: filename, format: format}
}

type fileSource struct {
	id       string
	format   string
	filepath string
}

func (f fileSource) String() string {
	return f.id
}

func (f fileSource) Read() (DataSet, error) {
	file, err := os.Open(f.filepath)
	if err != nil {
		return DataSet{}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return DataSet{}, err
	}

	stat, err := file.Stat()
	if err != nil {
		return DataSet{}, err
	}

	ds := DataSet{
		Data:      data,
		Format:    f.format,
		Source:    f.String(),
		Timestamp: stat.ModTime(),
	}
	ds.Checksum = "md5:" + ds.Md5()

	return ds, nil
}

func (f fileSource) Watch() (Watcher, error) {
	return newFileWatcher(f)
}

func newFileWatcher(fs fileSource) (Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(fs.filepath)
	return &fileWatcher{fs: fs, fw: fw, add: true, exit: make(chan struct{})}, nil
}

type fileWatcher struct {
	fs fileSource
	fw *fsnotify.Watcher

	add  bool
	exit chan struct{}
}

func (f *fileWatcher) Close() error {
	select {
	case <-f.exit:
	default:
		close(f.exit)
	}
	return f.fw.Close()
}

func (f *fileWatcher) Source() string {
	return f.fs.String()
}

func (f *fileWatcher) Next() (DataSet, error) {
	for {
		select {
		case <-f.exit:
			return DataSet{}, ErrWatcherClosed
		default:
		}

		if _, err := os.Stat(f.fs.filepath); err != nil {
			if os.IsNotExist(err) {
				if f.add {
					f.fw.Remove(f.fs.filepath)
					f.add = false
				}
				time.Sleep(time.Second * 10)
				continue
			}
			return DataSet{}, err
		}

		if !f.add {
			if err := f.fw.Add(f.fs.filepath); err != nil {
				return DataSet{}, err
			}
			f.add = true
			return f.fs.Read()
		}

		select {
		case event, ok := <-f.fw.Events:
			if !ok {
				return DataSet{}, ErrWatcherClosed
			}

			// BUG: it will be triggered twice continuously by fsnotify on Windows.
			if event.Op&fsnotify.Write == fsnotify.Write {
				return f.fs.Read()
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				f.add = false
			}
		case err, ok := <-f.fw.Errors:
			if !ok {
				return DataSet{}, ErrWatcherClosed
			}
			return DataSet{}, err
		case <-f.exit:
			return DataSet{}, ErrWatcherClosed
		}
	}
}
