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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// NewCliSource returns a new source based on "github.com/urfave/cli",
// which will reads the configuration data from the flags of cli.
//
// groups stands for the group that the context belongs on. The command name
// may be considered as the group name. The following ways are valid.
//
//   NewCliSource(ctx)                      // With the default global group
//   NewCliSource(ctx, "group1")            // With group "group1"
//   NewCliSource(ctx, "group1", "group2")  // With group "group1.group2"
//   NewCliSource(ctx, "group1.group2")     // With group "group1.group2"
//
func NewCliSource(ctx *cli.Context, groups ...string) Source {
	var group string
	if len(groups) > 0 {
		group = strings.Trim(strings.Join(groups, "."), ".")
	}
	return cliSource{ctx: ctx, group: group}
}

type cliSource struct {
	ctx   *cli.Context
	group string
}

func (c cliSource) String() string {
	return "cli"
}

func (c cliSource) Watch() (Watcher, error) {
	return nil, nil
}

func (c cliSource) Read() (DataSet, error) {
	names := c.ctx.FlagNames()
	if len(names) == 0 {
		names = c.ctx.GlobalFlagNames()
	}

	opts := make(map[string]string, 16)
	c.getFlags(c.group, c.ctx, names, opts)
	if len(opts) == 0 {
		return DataSet{}, nil
	}

	data, _ := json.Marshal(opts)
	ds := DataSet{
		Data:      data,
		Format:    "json",
		Source:    c.String(),
		Timestamp: time.Now(),
	}
	ds.Checksum = "md5:" + ds.Md5()
	return ds, nil
}

func (c cliSource) getFlags(group string, ctx *cli.Context, names []string, opts map[string]string) {
	for _, name := range names {
		key := name
		if group != "" {
			key = fmt.Sprintf("%s.%s", group, key)
		}

		if _, ok := opts[key]; ok {
			continue
		}

		opt := ctx.GlobalGeneric(name)
		if opt == nil {
			opt = ctx.Generic(name)
		}

		switch v := opt.(type) {
		case nil:
			continue
		case fmt.Stringer:
			opts[key] = v.String()
		default:
			panic(fmt.Errorf("unknown type '%T'", v))
		}
	}
}
