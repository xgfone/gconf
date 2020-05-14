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
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// NewZkSource is the same as NewZkConnSource, and conn is created by hosts,
// timeout, logger.
func NewZkSource(hosts []string, path, format string, timeout time.Duration, logger ...zk.Logger) Source {
	if format == "" {
		panic("zk source: the format must not be nil")
	}

	conn, _, err := zk.Connect(hosts, timeout)
	if err != nil {
		panic(err)
	} else if len(logger) > 0 && logger[0] != nil {
		conn.SetLogger(logger[0])
	}

	return NewZkConnSource(conn, path, format, hosts...)
}

// NewZkConnSource returns a new ZeeKeeper source with the connection to zk.
//
// path is the prefix of the zk path. If path is "", it is "/" by default.
//
// format is the format of the data of the zk source.
//
// hosts is used to generate the id of the source. If missing, it will use
// conn.Server().
func NewZkConnSource(conn *zk.Conn, path, format string, hosts ...string) Source {
	if format == "" {
		panic("zk source: the format must not be empty")
	}

	switch path = strings.TrimSpace(path); path {
	case "":
		path = "/"
	case "/":
	default:
		if path = strings.TrimRight(path, "/"); path == "" {
			path = "/"
		}
	}

	if len(hosts) == 0 {
		hosts = []string{conn.Server()}
	}

	id := fmt.Sprintf("zk:%s:%s", path, strings.Join(hosts, ","))
	return zkSource{path: path, zkconn: conn, id: id, format: format}
}

type zkSource struct {
	id     string
	path   string
	format string
	zkconn *zk.Conn
}

func (z zkSource) Read() (ds DataSet, err error) {
	data, stat, err := z.zkconn.Get(z.path)
	if err != nil {
		if err == zk.ErrNoNode {
			err = nil
		}
	} else {
		ds.Data = data
		ds.Source = z.id
		ds.Format = z.format
		ds.Checksum = ds.Md5()

		mt := stat.Mtime * int64(time.Millisecond)
		ds.Timestamp = time.Unix(mt/int64(time.Second), mt%int64(time.Second))
	}
	return
}

func (z zkSource) Watch(load func(DataSet, error) bool, exit <-chan struct{}) {
	go z.watchZkPath(load, exit)
}

func (z zkSource) watchZkPath(load func(DataSet, error) bool, exit <-chan struct{}) {
	last := DataSet{}
	interval := time.Second * 10

	for {
		switch data, stat, event, err := z.zkconn.GetW(z.path); err {
		case nil:
			ds := DataSet{Source: z.id, Format: z.format, Data: data}
			mt := stat.Mtime * int64(time.Millisecond)
			ds.Timestamp = time.Unix(mt/int64(time.Second), mt%int64(time.Second))
			ds.Checksum = ds.Md5()

			if len(ds.Data) > 0 && ds.Checksum != last.Checksum {
				if load(ds, nil) {
					last = ds
				}
			}

			select {
			case <-exit:
				z.zkconn.Close()
				return
			case _, ok := <-event:
				if ok {
					continue
				}
			}
		case zk.ErrNoNode:
		default:
			load(DataSet{Source: z.id, Format: z.format}, err)
		}
		time.Sleep(interval)
	}
}
