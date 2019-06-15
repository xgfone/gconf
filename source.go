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

// ErrWatcherClosed is returned when getting the source data after closing it.
var ErrWatcherClosed = fmt.Errorf("the watcher has been closed")

// SourceError represents an error about the data source.
type SourceError struct {
	Source string
	Format string
	Data   []byte
	Err    error
}

func (se SourceError) Error() string {
	return fmt.Sprintf("source[%s]: %s", se.Source, se.Err)
}

// NewSourceError returns a new source error.
func NewSourceError(source, format string, data []byte, err error) SourceError {
	return SourceError{Source: source, Format: format, Data: data, Err: err}
}

// DataSet represents the information of the configuration data.
type DataSet struct {
	Data      []byte // The original data.
	Format    string // Such as "json", "xml", etc.
	Source    string // Such as "file:/path/to/file", "zk:127.0.0.1:2181", etc.
	Checksum  string // Such as "md5:7d2f31e6fff478337478413ee1b70d2a", etc.
	Timestamp time.Time
}

// Md5 returns the md5 checksum of the DataSet data
func (c DataSet) Md5() string {
	return bytesToMd5(c.Data)
}

// Sha256 returns the sha256 checksum of the DataSet data
func (c DataSet) Sha256() string {
	return bytesToSha256(c.Data)
}

// Watcher is a watcher that is used to watch the data change of the source.
type Watcher interface {
	// Next returns the changed data from the source, which will be blocked
	// if no data is changed.
	//
	// After the watcher is closed, it should return ErrWatcherClosed.
	Next() (DataSet, error)

	// Source returns the source information, which is Source.String() in general.
	Source() string

	// Close closes the watcher.
	Close()
}

// Source represents a data source where the data is.
type Source interface {
	// String returns the information of the source,
	// such as "file:/path/to/file", "zk:127.0.0.1:2181", etc.
	String() string

	// Read reads the source data once, which should not block.
	Read() (DataSet, error)

	// Watch returns the watcher of the source.
	//
	// Return (nil, nil) if the source does not support the watcher.
	Watch() (Watcher, error)
}

func (c *Config) watchSource(w Watcher) {
	for {
		switch ds, err := w.Next(); err {
		case nil:
			if err = c.parseDataSet(ds, true); err != nil {
				c.handleError(err)
			}
		case ErrWatcherClosed:
			return
		default:
			c.handleError(NewSourceError(w.Source(), "", nil, err))
		}
	}
}

func (c *Config) parseDataSet(ds DataSet, force bool) error {
	if len(ds.Data) == 0 {
		return nil
	}

	decoder, ok := c.GetDecoder(ds.Format)
	if !ok {
		return NewSourceError(ds.Source, ds.Format, ds.Data, errNoDecoder)
	}

	// Decode the source data
	ms := make(map[string]interface{}, 32)
	if err := decoder.Decode(ds.Data, ms); err != nil {
		return NewSourceError(ds.Source, ds.Format, ds.Data, err)
	}

	// Load the map
	c.LoadMap(ms, force)
	return nil
}

// AddWatcher adds some watchers to watch the change of the configuration data.
//
// If the config is closed, it will ignore all the watchers.
func (c *Config) AddWatcher(watchers ...Watcher) {
	ws := make([]Watcher, 0, len(watchers))
	for _, w := range watchers {
		if w != nil {
			ws = append(ws, w)
		}
	}

	var exited bool
	c.lock.Lock()
	select {
	case <-c.exit:
		exited = true
	default:
		c.watchers = append(c.watchers, ws...)
	}
	c.lock.Unlock()

	if !exited {
		for _, w := range ws {
			go c.watchSource(w)
			debugf("[Config] Add source watcher '%s'\n", w.Source())
		}
	}
}

// LoadSource adds the sources then loads them. If the source supports
// the watcher, it will add it automatically.
//
// If a certain option of a certain group has been set, it will be ignored.
// But you can set force to true to reset the value of this option.
func (c *Config) LoadSource(source Source, force ...bool) error {
	return c.loadSource(source, true, force...)
}

// LoadSourceWithoutWatcher is the same as LoadSource, but does not add the watcher.
func (c *Config) LoadSourceWithoutWatcher(source Source, force ...bool) error {
	return c.loadSource(source, false, force...)
}

func (c *Config) loadSource(source Source, addWatcher bool, force ...bool) (err error) {
	var _force bool
	if len(force) > 0 && force[0] {
		_force = true
	}

	// Add the watcher if having one.
	var watcher Watcher
	if addWatcher {
		if watcher, err = source.Watch(); watcher != nil {
			c.AddWatcher(watcher)
		}
	}

	// Read and parse the data from the source for the first time.
	ds, rerr := source.Read()
	if rerr != nil {
		return NewSourceError(source.String(), "", nil, rerr)
	} else if rerr = c.parseDataSet(ds, _force); rerr != nil {
		return rerr
	}

	return
}
