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
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"strings"
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
	h := md5.New()
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Sha256 returns the sha256 checksum of the DataSet data
func (c DataSet) Sha256() string {
	h := sha256.New()
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
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
	Close() error
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

func (c *Config) flatMap(parent string, src, dst map[string]interface{}) {
	for key, value := range src {
		if ms, ok := value.(map[string]interface{}); ok {
			group := key
			if parent != "" {
				group = strings.Join([]string{parent, key}, c.gsep)
			}
			c.flatMap(group, ms, dst)
			continue
		}

		if parent != "" {
			key = strings.Join([]string{parent, key}, c.gsep)
		}
		dst[key] = value
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

	// Flat the map
	maps := make(map[string]interface{}, len(ms)*2)
	c.flatMap("", ms, maps)

	for key, value := range maps {
		group := c.OptGroup
		if index := strings.LastIndex(key, c.gsep); index > -1 {
			if group = c.Group(key[:index]); group == nil {
				continue
			}
			key = key[index+len(c.gsep):]
		}

		if force || group.HasOptAndIsNotSet(key) {
			group.Set(key, value)
		}
	}
	return nil
}

// AddWatcher adds some watchers to watch the change of the configuration data.
func (c *Config) AddWatcher(watchers ...Watcher) {
	ws := make([]Watcher, 0, len(watchers))
	for _, w := range watchers {
		if w != nil {
			ws = append(ws, w)
		}
	}

	c.lock.Lock()
	c.watchers = append(c.watchers, ws...)
	c.lock.Unlock()

	for _, w := range ws {
		go c.watchSource(w)
		debugf("[Config] Add source watcher '%s'", w.Source())
	}
}

// LoadSource adds the sources then loads them.
//
// If a certain option of a certain group has been set, it will be ignored.
// But you can set force to true to reset the value of this option.
func (c *Config) LoadSource(source Source, force ...bool) error {
	var _force bool
	if len(force) > 0 && force[0] {
		_force = true
	}

	// Add the watcher if having one.
	watcher, werr := source.Watch()
	if watcher != nil {
		c.AddWatcher(watcher)
	}

	// Read and parse the data from the source for the first time.
	ds, err := source.Read()
	if err != nil {
		return NewSourceError(source.String(), "", nil, err)
	} else if err = c.parseDataSet(ds, _force); err != nil {
		return err
	}

	return werr
}
