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
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var defaultDebug bool

func init() {
	Conf.RegisterOpt(ConfigFileOpt)

	for _, env := range os.Environ() {
		index := strings.IndexByte(env, '=')
		if index == -1 {
			continue
		}

		if strings.ToUpper(env[:index]) == "DEBUG" {
			if v, _ := strconv.ParseBool(env[index+1:]); v {
				defaultDebug = v
			}
			break
		}
	}
}

// debugf prints the log message only when enabling the debug mode.
func debugf(msg string, args ...interface{}) {
	if defaultDebug {
		fmt.Printf(msg, args...)
	}
}

// Conf is the default global Config.
//
// The default global Conf will register the option ConfigFileOpt.
var Conf = New()

type changedOpt struct {
	Group string
	Opt   string
	Old   interface{}
	New   interface{}
}

// Config is used to manage the configuration options.
type Config struct {
	*OptGroup  // The default group.
	errHandler func(error)

	exit chan struct{}
	lock sync.RWMutex
	gsep string // The separator between the group names.

	opts     chan changedOpt
	groups   map[string]*OptGroup // the option groups
	groups2  map[string]*OptGroup // The auxiliary groups
	decoders map[string]Decoder
	decAlias map[string]string
	watchers []Watcher
	observes []func(string, string, interface{}, interface{})
	version  Opt
}

// New returns a new Config.
func New() *Config {
	c := new(Config)
	c.gsep = "."
	c.exit = make(chan struct{})
	c.opts = make(chan changedOpt, 1024)
	c.groups = make(map[string]*OptGroup, 8)
	c.groups2 = make(map[string]*OptGroup, 8)
	c.decoders = make(map[string]Decoder, 8)
	c.decAlias = make(map[string]string, 8)
	c.watchers = make([]Watcher, 0, 8)
	c.OptGroup = newOptGroup(c, "")
	c.groups[c.OptGroup.Name()] = c.OptGroup
	c.errHandler = c.defaultErrorHandler
	c.AddDecoder(NewJSONDecoder())
	c.AddDecoder(NewIniDecoder())
	c.AddDecoderAlias("conf", "ini")
	go c.watchChangedOpt()
	return c
}

func (c *Config) newGroup(parent, name string) *OptGroup {
	if strings.Contains(name, c.gsep) {
		panic(fmt.Errorf("the group name '%s' must not contain the group separator '%s'", name, c.gsep))
	}

	c.lock.Lock()
	name = c.getGroupName(parent, name)
	group, ok := c.groups[name]
	if !ok {
		group = newOptGroup(c, name)
		c.groups[name] = group
		c.ensureGroup2(name)
	}
	c.lock.Unlock()

	if !ok {
		debugf("[Config] Creating a new group '%s'\n", name)
	}
	return group
}

func (c *Config) ensureGroup2(name string) {
	if gnames := strings.Split(name, c.gsep); len(gnames) >= 1 {
		for i := range gnames {
			gname := strings.Join(gnames[:i+1], c.gsep)
			if c.groups[gname] == nil && c.groups2[gname] == nil {
				c.groups2[gname] = newOptGroup(c, gname)
				debugf("[Config] Creating the auxiliary group '%s'\n", gname)
			}
		}
	}
}

func (c *Config) getGroupName(parent, name string) string {
	name = strings.Trim(name, c.gsep)
	if parent == "" {
		return name
	} else if name == "" {
		return parent
	}
	return strings.Join([]string{parent, name}, c.gsep)
}

func (c *Config) getGroup(parent, name string) *OptGroup {
	c.lock.RLock()
	name = c.getGroupName(parent, name)
	group, ok := c.groups[name]
	if !ok {
		group = c.groups[name]
	}
	c.lock.RUnlock()
	return group
}

func (c *Config) noticeOptChange(group, optname string, old, new interface{}) {
	select {
	case c.opts <- changedOpt{Group: group, Opt: optname, Old: old, New: new}:
	default:
	}
}

func (c *Config) watchChangedOpt() {
	observes := make([]func(string, string, interface{}, interface{}), 0, 32)
	for {
		select {
		case <-c.exit:
			return
		case opt := <-c.opts:
			c.lock.RLock()
			observes = append(observes[:0], c.observes...)
			c.lock.RUnlock()

			for _, observe := range observes {
				c.callObserver(opt, observe)
			}
		}
	}
}

func (c *Config) wrapPanic() {
	if err := recover(); err != nil {
		c.handleError(fmt.Errorf("[Config] Observer panic: %v", err))
	}
}

func (c *Config) callObserver(opt changedOpt, f func(string, string, interface{}, interface{})) {
	defer c.wrapPanic()
	f(opt.Group, opt.Opt, opt.Old, opt.New)
}

// Close closes all the watchers and disables anyone to add the watcher into it.
func (c *Config) Close() {
	select {
	case <-c.exit:
	default:
		close(c.exit)
		c.lock.RLock()
		defer c.lock.RUnlock()
		for _, w := range c.watchers {
			w.Close()
		}
	}
}

func (c *Config) defaultErrorHandler(err error) {
	if !IsErrNoOpt(err) {
		fmt.Println(err)
	}
}

func (c *Config) handleError(err error) {
	c.lock.RLock()
	handler := c.errHandler
	c.lock.RUnlock()
	handler(err)
}

// SetErrHandler resets the error handler to h.
//
// The default is output to os.Stdout by fmt.Println(err), but it ignores ErrNoOpt.
func (c *Config) SetErrHandler(h func(error)) {
	if h == nil {
		panic("the error handler must not be nil")
	}

	c.lock.Lock()
	c.errHandler = h
	c.lock.Unlock()
}

// Observe appends the observer to watch the change of all the option value.
func (c *Config) Observe(observer func(group string, opt string, oldValue, newValue interface{})) {
	if observer == nil {
		panic("the observer must not be nil")
	}
	c.lock.Lock()
	c.observes = append(c.observes, observer)
	c.lock.Unlock()
}

// AllGroups returns all the groups, containing the default group.
func (c *Config) AllGroups() []*OptGroup {
	c.lock.RLock()
	groups := make([]*OptGroup, len(c.groups))
	var index int
	for _, group := range c.groups {
		groups[index] = group
		index++
	}
	c.lock.RUnlock()

	sort.Slice(groups, func(i, j int) bool { return groups[i].Name() < groups[j].Name() })
	return groups
}

// SetStringVersion is equal to c.SetVersion(VersionOpt.D(version)).
func (c *Config) SetStringVersion(version string) {
	c.SetVersion(VersionOpt.D(version))
}

// SetVersion sets the version information.
//
// Notice: the field Default must be a string.
func (c *Config) SetVersion(version Opt) {
	if v, ok := version.Default.(string); !ok {
		panic("the version is not a string value")
	} else if v == "" {
		panic("the version must not be empty")
	}

	version.check()
	c.lock.Lock()
	c.version = version
	c.lock.Unlock()
}

// GetVersion returns a the version information.
//
// Notice: the Default filed is a string representation of the version value.
// But it is "" if no version.
func (c *Config) GetVersion() (version Opt) {
	c.lock.RLock()
	version = c.version
	c.lock.RUnlock()
	return
}

// PrintGroup prints the information of all groups to w.
func (c *Config) PrintGroup(w io.Writer) error {
	for _, group := range c.AllGroups() {
		if gname := group.Name(); gname == "" {
			fmt.Fprintf(w, "[DEFAULT]\n")
		} else {
			fmt.Fprintf(w, "[%s]\n", gname)
		}

		for _, opt := range group.AllOpts() {
			fmt.Fprintf(w, "    %s\n", opt.Name)
		}
	}
	return nil
}

// UpdateOptValue updates the value of the option of the group.
//
// If the group or the option does not exist, it will be ignored.
func (c *Config) UpdateOptValue(groupName, optName string, optValue interface{}) {
	if group := c.Group(groupName); group != nil {
		group.Set(optName, optValue)
	}
}

// UpdateValue is the same as UpdateOptValue, but key is equal to
// `fmt.Sprintf("%s.%s", groupName, optName)`.
//
// that's,
//   c.UpdateOptValue(groupName, optName, optValue)
// is equal to
//   c.UpdateValue(fmt.Sprintf("%s.%s", groupName, optName), optValue)
func (c *Config) UpdateValue(key string, value interface{}) {
	var group string
	if index := strings.LastIndex(key, c.gsep); index > 0 {
		group = key[:index]
		key = key[index+len(c.gsep):]
	}
	c.UpdateOptValue(group, key, value)
}
