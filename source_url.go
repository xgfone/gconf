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
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/xgfone/go-tools/io2"
)

var errNoContentType = fmt.Errorf("http response has no the header Content-Type")

// NewURLSource returns a url source to read the configuration data from the url
// by the stdlib http.Get(url).
//
// The header "Content-Type" indicates the data format, that's, it will split
// the value by "/" and use the last part, such as "application/json" represents
// the format "json". But you can set format to override it.
//
// It supports the watcher, which checks whether the data is changed once
// with interval time. If interval is 0, it is time.Minute by default.
func NewURLSource(url string, interval time.Duration, format ...string) Source {
	if url == "" {
		panic("the url must not be nil")
	} else if _, err := neturl.Parse(url); err != nil {
		panic(err)
	}

	var _format string
	if len(format) > 0 && format[0] != "" {
		_format = format[0]
	}
	if interval <= 0 {
		interval = time.Minute
	}
	return urlSource{id: fmt.Sprintf("url:%s", url), url: url, format: _format, period: interval}
}

type urlSource struct {
	id  string
	url string

	format string
	period time.Duration
}

func (u urlSource) String() string {
	return u.id
}

func (u urlSource) Read() (DataSet, error) {
	resp, err := http.Get(u.url)
	if err != nil {
		return DataSet{}, err
	}
	defer resp.Body.Close()

	format := u.format
	if format == "" {
		// Get the Content-Type as the format.
		ct := strings.TrimSpace(resp.Header.Get("Content-Type"))
		if index := strings.IndexByte(ct, ';'); index > 0 {
			ct = strings.TrimSpace(ct[:index])
		}
		if index := strings.LastIndexByte(ct, '/'); index > 0 {
			ct = ct[index+1:]
		}
		if ct == "" {
			return DataSet{}, errNoContentType
		}
		format = ct
	}

	// Read the body of the response.
	data, err := io2.ReadN(resp.Body, resp.ContentLength)
	if err != nil {
		return DataSet{}, err
	}

	ds := DataSet{
		Data:      data,
		Format:    format,
		Source:    u.String(),
		Timestamp: time.Now(),
	}
	ds.Checksum = "md5:" + ds.Md5()
	return ds, nil
}

func (u urlSource) Watch() (Watcher, error) {
	return newURLWatcher(u, u.period), nil
}

func newURLWatcher(src Source, interval time.Duration) Watcher {
	w := urlWatcher{
		src:   src,
		exit:  make(chan struct{}),
		value: make(chan interface{}, 1),
		sleep: interval,
	}
	go w.loop()
	return w
}

type urlWatcher struct {
	src   Source
	last  DataSet
	exit  chan struct{}
	value chan interface{}
	sleep time.Duration
}

func (u urlWatcher) loop() {
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(u.sleep)
		}

		var value interface{}
		if ds, err := u.src.Read(); err != nil {
			value = err
		} else if len(ds.Data) == 0 || ds.Checksum == u.last.Checksum {
			continue
		} else {
			value = ds
			u.last = ds
		}

		// Maybe have the old value and consume it.
		select {
		case _, ok := <-u.exit: // be closed
			if !ok {
				return
			}
		case _, ok := <-u.value:
			if !ok { // be closed
				return
			}
		default:
		}

		select {
		case _, ok := <-u.exit: // be closed
			if !ok {
				return
			}
		case u.value <- value: // Send the new value.
		default:
		}
	}
}

func (u urlWatcher) Next() (DataSet, error) {
	if value, ok := <-u.value; !ok {
		return DataSet{}, ErrWatcherClosed
	} else if err, ok := value.(error); ok {
		return DataSet{}, err
	} else {
		return value.(DataSet), nil
	}
}

func (u urlWatcher) Source() string {
	return u.src.String()
}

func (u urlWatcher) Close() error {
	select {
	case <-u.exit:
	default:
		close(u.exit)
		close(u.value)
	}
	return nil
}
