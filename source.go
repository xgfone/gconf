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
	"time"
)

var errNoDecoder = fmt.Errorf("no decoder")

// SourceError represents an error about the data source.
type SourceError struct {
	Source string
	Format string
	Data   []byte
	Err    error
}

func (se SourceError) Error() string {
	return fmt.Sprintf("source(%s)[%s]: %s", se.Source, se.Format, se.Err.Error())
}

// NewSourceError returns a new source error.
func NewSourceError(source, format string, data []byte, err error) SourceError {
	return SourceError{Source: source, Format: format, Data: data, Err: err}
}

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

// LoadDataSet loads the DataSet ds, which will parse the data by calling the
// corresponding decoder and load it.
//
// If a certain option has been set, it will be ignored. But you can set force
// to true to reset the value of this option.
func (c *Config) LoadDataSet(ds DataSet, force ...bool) (err error) {
	if len(ds.Data) == 0 {
		return nil
	}

	decoder, ok := c.GetDecoder(ds.Format)
	if !ok {
		err = NewSourceError(ds.Source, ds.Format, ds.Data, errNoDecoder)
		c.handleError(err)
		return err
	}

	ms := make(map[string]interface{}, 32)
	if err := decoder.Decode(ds.Data, ms); err != nil {
		err = NewSourceError(ds.Source, ds.Format, ds.Data, err)
		c.handleError(err)
		return err
	}

	if err = c.LoadMap(ms, force...); err == nil && ds.Args != nil {
		if c.Args() == nil || (len(force) > 0 && force[0]) {
			c.SetArgs(ds.Args)
		}
	}
	return
}

func (c *Config) loadDataSetWithError(ds DataSet, err error, force ...bool) (ok bool) {
	switch err.(type) {
	case nil:
		ok = c.LoadDataSet(ds, force...) == nil
	case SourceError:
		c.handleError(err)
	default:
		c.handleError(NewSourceError(ds.Source, ds.Format, ds.Data, err))
	}
	return
}

// LoadDataSetCallback is a callback used by the watcher.
func (c *Config) LoadDataSetCallback(ds DataSet, err error) bool {
	return c.loadDataSetWithError(ds, err, true)
}

// Source represents a data source where the data is.
type Source interface {
	// Read reads the source data once, which should not block.
	Read() (DataSet, error)

	// Watch watches the change of the source, then send the changed data to ds.
	//
	// The source can check whether close is closed to determine whether the
	// configuration is closed and to do any cleanup.
	//
	// Notice: load returns true only the DataSet is loaded successfully.
	Watch(load func(DataSet, error) bool, close <-chan struct{})
}

func (c *Config) loadSource(source Source, force ...bool) {
	ds, err := source.Read()
	c.loadDataSetWithError(ds, err, force...)
}

// LoadSource loads the sources, and call Watch to watch the source, which is
// equal to
//
//   c.LoadSourceWithoutWatch(source, force...)
//   source.Watch(c.LoadDataSetCallback, c.CloseNotice())
//
// When loading the source, if a certain option of a certain group has been set,
// it will be ignored. But you can set force to true to reset the value of this
// option.
func (c *Config) LoadSource(source Source, force ...bool) {
	c.loadSource(source, force...)
	source.Watch(c.LoadDataSetCallback, c.exit)
}

// LoadSourceWithoutWatch is the same as LoadSource, but does not call
// the Watch method of the source to watch the source.
func (c *Config) LoadSourceWithoutWatch(source Source, force ...bool) {
	c.loadSource(source, force...)
}
