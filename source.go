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
	"errors"
	"time"
)

// ErrNoDecoder represents the error that there is no decoder.
var ErrNoDecoder = errors.New("no decoder")

// DataSet represents the information of the configuration data.
type DataSet struct {
	Args      []string  // The CLI arguments filled by the CLI source such as flag.
	Data      []byte    // The original data.
	Format    string    // Such as "json", "xml", etc.
	Source    string    // Such as "file:/path/to/file", "zk:127.0.0.1:2181", etc.
	Checksum  string    // Such as "md5:7d2f31e6fff478337478413ee1b70d2a", etc.
	Timestamp time.Time // The timestamp when the data is modified.
}

// Md5 returns the md5 checksum of the DataSet data
func (ds DataSet) Md5() string {
	return bytesToMd5(ds.Data)
}

// Sha256 returns the sha256 checksum of the DataSet data
func (ds DataSet) Sha256() string {
	return bytesToSha256(ds.Data)
}

// LoadDataSet is equal to Conf.LoadDataSet(source, force...).
func LoadDataSet(ds DataSet, force ...bool) error {
	return Conf.LoadDataSet(ds, force...)
}

// LoadDataSet loads the DataSet ds, which will parse the data by calling the
// corresponding decoder and load it.
//
// If failing to parse the value of any option, it terminates to parse
// and load it.
//
// If force is missing or false, ignore the assigned options.
func (c *Config) LoadDataSet(ds DataSet, force ...bool) (err error) {
	if len(ds.Data) == 0 {
		return nil
	}

	decoder := c.GetDecoder(ds.Format)
	if decoder == nil {
		return ErrNoDecoder
	}

	ms := make(map[string]interface{}, 32)
	if err = decoder(ds.Data, ms); err != nil {
		return err
	}

	if err = c.LoadMap(ms, force...); err == nil && ds.Args != nil {
		if c.Args == nil || (len(force) > 0 && force[0]) {
			c.Args = ds.Args
		}
	}

	return
}

// Source represents a data source where the data is.
type Source interface {
	// String is the description of the source, such as "env", "file:/path/to".
	String() string

	// Read reads the source data once, which should not block.
	Read() (DataSet, error)

	// Watch watches the change of the source, then call the callback load.
	//
	// close is used to notice the underlying watcher to close and clean.
	Watch(close <-chan struct{}, load func(DataSet, error) (success bool))
}

// LoadSource is equal to Conf.LoadSource(source, force...).
func LoadSource(source Source, force ...bool) error {
	return Conf.LoadSource(source, force...)
}

// LoadSource loads the options from the given source.
//
// If force is missing or false, ignore the assigned options.
func (c *Config) LoadSource(source Source, force ...bool) (err error) {
	ds, err := source.Read()
	if err != nil {
		c.errorf("fail to read the source '%s': %s", source.String(), err)
		return
	}

	if err = c.LoadDataSet(ds, force...); err != nil {
		c.errorf("fail to load the source '%s': %s", source.String(), err)
		return
	}

	return
}

// LoadAndWatchSource is equal to Conf.LoadAndWatchSource(source, force...).
func LoadAndWatchSource(source Source, force ...bool) error {
	return Conf.LoadAndWatchSource(source, force...)
}

// LoadAndWatchSource is the same as LoadSource, but also watches the source
// after loading the source successfully.
func (c *Config) LoadAndWatchSource(source Source, force ...bool) (err error) {
	if err = c.LoadSource(source, force...); err == nil {
		go source.Watch(c.exit, func(ds DataSet, err error) bool {
			if err != nil {
				c.errorf("fail to watch the source '%s': %s", source, err)
				return false
			} else if err = c.LoadDataSet(ds, true); err != nil {
				c.errorf("fail to load the source '%s': %s", source, err)
				return false
			}
			return true
		})
	}
	return
}
