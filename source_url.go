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
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

var errNoContentType = fmt.Errorf("http response has no the header Content-Type")

// NewURLSource returns a url source to read the configuration data
// from the url by the stdlib http.Get(url).
//
// The header "Content-Type" indicates the data format, that's, it will split
// the value by "/" and use the last part, such as "application/json" represents
// the format "json". But you can set format to override it.
//
// The url source can watch the configuration data from the url each interval
// period. If interval is equal to 0, it is defaulted to time.Minute.
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

	return urlSource{
		id:     fmt.Sprintf("url:%s", url),
		url:    url,
		format: _format,
		period: interval,
	}
}

type urlSource struct {
	id  string
	url string

	format string
	period time.Duration
}

func (u urlSource) String() string { return u.id }

func (u urlSource) Read() (DataSet, error) {
	resp, err := http.Get(u.url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return DataSet{Source: u.id, Format: u.format}, err
	}

	format := u.format
	if format == "" {
		// Get the format from the header "Content-Type".
		ct := strings.TrimSpace(resp.Header.Get("Content-Type"))
		if index := strings.IndexByte(ct, ';'); index > 0 {
			ct = strings.TrimSpace(ct[:index])
		}
		if index := strings.LastIndexByte(ct, '/'); index > 0 {
			ct = ct[index+1:]
		}
		if ct == "" {
			return DataSet{Source: u.id}, errNoContentType
		}
		format = ct
	}

	// Read the body of the response.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DataSet{Source: u.id, Format: format}, err
	}

	ds := DataSet{
		Data:      data,
		Format:    format,
		Source:    u.id,
		Timestamp: time.Now(),
	}
	ds.Checksum = "md5:" + ds.Md5()
	return ds, nil
}

func (u urlSource) Watch(exit <-chan struct{}, load func(DataSet, error) bool) {
	u.watch(exit, load)
}

func (u urlSource) watch(exit <-chan struct{}, load func(DataSet, error) bool) {
	ticker := time.NewTicker(u.period)
	defer ticker.Stop()

	var last DataSet
	for {
		select {
		case <-exit:
			return

		case <-ticker.C:
			if ds, err := u.Read(); err != nil {
				load(ds, err)
			} else if len(ds.Data) > 0 && ds.Checksum != last.Checksum {
				if load(ds, nil) {
					last = ds
				}
			}
		}
	}
}
