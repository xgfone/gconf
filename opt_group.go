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
	"strings"
	"time"
)

// OptGroup is the proxy for a group of options.
type OptGroup struct {
	prefix string
	config *Config
}

// Group is equal to Conf.Group(name).
func Group(name string) *OptGroup { return Conf.Group(name) }

// GetGroupSep returns the separator between the option groups.
func (c *Config) GetGroupSep() (sep string) { return c.gsep }

// Group returns the option group with the group name.
func (c *Config) Group(name string) *OptGroup {
	sep := c.GetGroupSep()
	gname := strings.TrimSuffix(strings.TrimPrefix(name, sep), sep)
	if gname == "" {
		return &OptGroup{config: c}
	}

	return &OptGroup{
		prefix: gname + sep,
		config: c,
	}
}

// Group returns a sub-group with the group name.
func (g *OptGroup) Group(name string) *OptGroup {
	sep := g.config.GetGroupSep()
	gname := strings.TrimSuffix(strings.TrimPrefix(name, sep), sep)
	if gname == "" {
		return &OptGroup{prefix: g.prefix, config: g.config}
	}

	return &OptGroup{
		prefix: g.prefix + gname + sep,
		config: g.config,
	}
}

// Prefix returns the prefix of the group.
func (g *OptGroup) Prefix() string { return g.prefix }

// Self makes itself into a bool option proxy with the default value "false".
func (g *OptGroup) Self(help string) *OptProxyBool {
	if g.prefix == "" {
		panic("the group name is empty")
	}
	return g.config.NewBool(strings.TrimSuffix(g.prefix, g.config.gsep), false, help)
}

// RegisterOpts registers a set of options.
//
// Notice: if a certain option has existed, it will panic.
func (g *OptGroup) RegisterOpts(opts ...Opt) {
	_opts := make([]Opt, len(opts))
	for i, opt := range opts {
		opt.Name = g.prefix + opt.Name
		_opts[i] = opt
	}
	g.config.RegisterOpts(_opts...)
}

// UnregisterOpts unregisters the registered options.
func (g *OptGroup) UnregisterOpts(optNames ...string) {
	names := make([]string, len(optNames))
	for _len := len(optNames) - 1; _len >= 0; _len-- {
		names[_len] = g.prefix + optNames[_len]
	}
	g.config.UnregisterOpts(names...)
}

// Set resets the option named name to value.
func (g *OptGroup) Set(name string, value interface{}) error {
	return g.config.Set(g.prefix+name, value)
}

// Get returns the value of the option named name.
//
// Return nil if this option does not exist.
func (g *OptGroup) Get(name string) interface{} {
	return g.config.Get(g.prefix + name)
}

// Must is the same as Get, but panic if the returned value is nil.
func (g *OptGroup) Must(name string) interface{} {
	return g.config.Must(g.prefix + name)
}

// GetBool returns the value of the option named name as bool.
func (g *OptGroup) GetBool(name string) bool { return g.Must(name).(bool) }

// GetInt returns the value of the option named name as int.
func (g *OptGroup) GetInt(name string) int { return g.Must(name).(int) }

// GetInt32 returns the value of the option named name as int32.
func (g *OptGroup) GetInt32(name string) int32 { return g.Must(name).(int32) }

// GetInt64 returns the value of the option named name as int64.
func (g *OptGroup) GetInt64(name string) int64 { return g.Must(name).(int64) }

// GetUint returns the value of the option named name as uint.
func (g *OptGroup) GetUint(name string) uint { return g.Must(name).(uint) }

// GetUint32 returns the value of the option named name as uint32.
func (g *OptGroup) GetUint32(name string) uint32 { return g.Must(name).(uint32) }

// GetUint64 returns the value of the option named name as uint64.
func (g *OptGroup) GetUint64(name string) uint64 { return g.Must(name).(uint64) }

// GetFloat64 returns the value of the option named name as float64.
func (g *OptGroup) GetFloat64(name string) float64 { return g.Must(name).(float64) }

// GetString returns the value of the option named name as string.
func (g *OptGroup) GetString(name string) string { return g.Must(name).(string) }

// GetDuration returns the value of the option named name as time.Duration.
func (g *OptGroup) GetDuration(name string) time.Duration { return g.Must(name).(time.Duration) }

// GetTime returns the value of the option named name as time.Time.
func (g *OptGroup) GetTime(name string) time.Time { return g.Must(name).(time.Time) }

// GetIntSlice returns the value of the option named name as []int.
func (g *OptGroup) GetIntSlice(name string) []int { return g.Must(name).([]int) }

// GetUintSlice returns the value of the option named name as []uint.
func (g *OptGroup) GetUintSlice(name string) []uint { return g.Must(name).([]uint) }

// GetFloat64Slice returns the value of the option named name as []float64.
func (g *OptGroup) GetFloat64Slice(name string) []float64 { return g.Must(name).([]float64) }

// GetStringSlice returns the value of the option named name as []string.
func (g *OptGroup) GetStringSlice(name string) []string { return g.Must(name).([]string) }

// GetDurationSlice returns the value of the option named name as []time.Duration.
func (g *OptGroup) GetDurationSlice(name string) []time.Duration {
	return g.Must(name).([]time.Duration)
}
