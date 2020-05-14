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
	Conf.RegisterOpts(ConfigFileOpt)

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
		printMsg(fmt.Sprintf(msg, args...))
	}
}

func printMsg(msg string) {
	switch DefaultWriter {
	case os.Stdout, os.Stderr:
		fmt.Fprintln(DefaultWriter, msg)
	default:
		io.WriteString(DefaultWriter, msg)
	}
}

// DefaultWriter is the default writer, which Config will write the information
// to it by default.
var DefaultWriter = os.Stdout

// Conf is the default global Config.
//
// The default global Conf will register the option ConfigFileOpt.
var Conf = New()

// Config is used to manage the configuration options.
type Config struct {
	*OptGroup  // The default group.
	errHandler func(error)

	exit chan struct{}
	lock sync.RWMutex
	gsep string // The separator between the group names.

	snap     *snapshot
	groups   map[string]*OptGroup // the option groups
	groups2  map[string]*OptGroup // The auxiliary groups
	decoders map[string]Decoder
	decAlias map[string]string
	version  Opt

	observes []func(string, string, interface{}, interface{})
}

// New returns a new Config.
//
// By default, it will add the "json", "yaml", "toml" and "ini" decoders,
// and set the aliases of "conf" and "yml" to "ini" and "yaml", for example,
//
//   c.AddDecoder(NewJSONDecoder())
//   c.AddDecoder(NewIniDecoder())
//   c.AddDecoder(NewYamlDecoder())
//   c.AddDecoder(NewTomlDecoder())
//   c.AddDecoderAlias("conf", "ini")
//   c.AddDecoderAlias("yml", "yaml")
//
func New() *Config {
	c := new(Config)
	c.gsep = "."
	c.snap = newSnapshot(c)
	c.exit = make(chan struct{})
	c.groups = make(map[string]*OptGroup, 8)
	c.groups2 = make(map[string]*OptGroup, 8)
	c.decoders = make(map[string]Decoder, 8)
	c.decAlias = make(map[string]string, 8)
	c.OptGroup = newOptGroup(c, "")
	c.groups[c.OptGroup.Name()] = c.OptGroup
	c.errHandler = ErrorHandler(func(err error) { printMsg(err.Error()) })
	c.AddDecoder(NewJSONDecoder())
	c.AddDecoder(NewIniDecoder())
	c.AddDecoder(NewYamlDecoder())
	c.AddDecoder(NewTomlDecoder())
	c.AddDecoderAlias("conf", "ini")
	c.AddDecoderAlias("yml", "yaml")
	return c
}

func (c *Config) _newGroup(parent, name string) *OptGroup {
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
		debugf("[Config] Creating a new group '%s'", name)
	}
	return group
}

func (c *Config) newGroup(parent, name string) (group *OptGroup) {
	for _, subname := range strings.Split(name, c.gsep) {
		group = c._newGroup(parent, subname)
		parent = c.getGroupName(parent, subname)
	}
	return
}

func (c *Config) ensureGroup2(name string) {
	if gnames := strings.Split(name, c.gsep); len(gnames) >= 1 {
		for i := range gnames {
			gname := strings.Join(gnames[:i+1], c.gsep)
			if c.groups[gname] == nil && c.groups2[gname] == nil {
				c.groups2[gname] = newOptGroup(c, gname)
				debugf("[Config] Creating the auxiliary group '%s'", gname)
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

func (c *Config) noticeOptChange(group, optname string, old, new interface{},
	observers []func(interface{})) {
	for _, observer := range observers {
		c.callOptObserver(observer, new)
	}

	c.lock.RLock()
	fs := append([]func(g, p string, o, n interface{}){}, c.observes...)
	c.lock.RUnlock()
	for _, observer := range fs {
		c.callSetObserver(group, optname, old, new, observer)
	}
}

func (c *Config) callOptObserver(observe func(interface{}), new interface{}) {
	defer c.wrapPanic("opt")
	observe(new)
}

func (c *Config) callSetObserver(group, optname string, old, new interface{},
	cb func(string, string, interface{}, interface{})) {
	defer c.wrapPanic("set")
	cb(group, optname, old, new)
}

func (c *Config) wrapPanic(s string) {
	if err := recover(); err != nil {
		c.handleError(fmt.Errorf("[Config] option %s observer panic: %v", s, err))
	}
}

// Close closes all the watchers and disables anyone to add the watcher into it.
func (c *Config) Close() {
	select {
	case <-c.exit:
	default:
		close(c.exit)
	}
}

// CloseNotice returns a close channel, which will also be closed when the config
// is closed.
func (c *Config) CloseNotice() <-chan struct{} {
	return c.exit
}

func (c *Config) handleError(err error) {
	c.lock.RLock()
	handler := c.errHandler
	c.lock.RUnlock()
	handler(err)
}

// ErrorHandler returns a error handler, which will ignore ErrNoOpt
// and ErrFrozenOpt, and pass the others to h.
func ErrorHandler(h func(err error)) func(error) {
	return func(err error) {
		if !IsErrNoOpt(err) && !IsErrFrozenOpt(err) {
			h(err)
		}
	}
}

// SetErrHandler resets the error handler to h.
//
// The default is output to DefaultWriter, but it ignores ErrNoOpt and ErrFrozenOpt.
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

// Traverse traverses all the options of all the groups.
func (c *Config) Traverse(f func(group string, opt string, value interface{})) {
	for _, group := range c.AllGroups() {
		for _, opt := range group.AllOpts() {
			name := group.fixOptName(opt.Name)
			f(group.Name(), name, group.Get(name))
		}
	}
}

/// --------------------------------------------------------------------------

// UpdateOptValue updates the value of the option of the group.
//
// If the group or the option does not exist, it will be ignored.
func (c *Config) UpdateOptValue(groupName, optName string, optValue interface{}) (err error) {
	if group := c.Group(groupName); group != nil {
		err = group.Set(optName, optValue)
	}
	return
}

// UpdateValue is the same as UpdateOptValue, but key is equal to
// `fmt.Sprintf("%s.%s", groupName, optName)`.
//
// that's,
//   c.UpdateOptValue(groupName, optName, optValue)
// is equal to
//   c.UpdateValue(fmt.Sprintf("%s.%s", groupName, optName), optValue)
func (c *Config) UpdateValue(key string, value interface{}) error {
	var group string
	if index := strings.LastIndex(key, c.gsep); index > 0 {
		group = key[:index]
		key = key[index+len(c.gsep):]
	}
	return c.UpdateOptValue(group, key, value)
}

// LoadMap loads the configuration options from the map m and returns true
// only if all options are parsed and set successfully.
//
// If a certain option has been set, it will be ignored.
// But you can set force to true to reset the value of this option.
//
// The map may be the formats as follow:
//
//     map[string]interface{} {
//         "opt1": "value1",
//         "opt2": "value2",
//         // ...
//         "group1": map[string]interface{} {
//             "opt11": "value11",
//             "opt12": "value12",
//             "group2": map[string]interface{} {
//                 // ...
//             },
//             "group3.group4": map[string]interface{} {
//                 // ...
//             }
//         },
//         "group5.group6.group7": map[string]interface{} {
//             "opt71": "value71",
//             "opt72": "value72",
//             "group8": map[string]interface{} {
//                 // ...
//             },
//             "group9.group10": map[string]interface{} {
//                 // ...
//             }
//         },
//         "group11.group12.opt121": "value121",
//         "group11.group12.opt122": "value122"
//     }
//
func (c *Config) LoadMap(m map[string]interface{}, force ...bool) error {
	opts := make([]groupOptValue, 0, len(m))
	opts, err := c.parseMap(c.OptGroup, m, opts)
	if err != nil {
		c.handleError(err)
		return err
	}

	var _force bool
	if len(force) > 0 && force[0] {
		_force = true
	}

	c.loadMap(opts, _force)
	return nil
}

type groupOptValue struct {
	Group *OptGroup
	Name  string
	Value interface{}
}

func (c *Config) parseMap(g *OptGroup, m map[string]interface{},
	opts []groupOptValue) ([]groupOptValue, error) {
	var err error
	var ms map[string]interface{}

	for key, value := range m {
		// fmt.Printf("@@@@@ %s: %#v\n", key, value)
		switch m := value.(type) {
		case map[string]interface{}:
			if _g := g.Group(key); _g != nil {
				if opts, err = c.parseMap(_g, m, opts); err != nil {
					return opts, err
				}
			}
		case map[interface{}]interface{}:
			if ms, err = toStringMap(m); err != nil {
				return opts, err
			}

			if _g := g.Group(key); _g != nil {
				if opts, err = c.parseMap(_g, ms, opts); err != nil {
					return opts, err
				}
			}
		default:
			_g := g
			if index := strings.LastIndex(key, c.gsep); index > 0 {
				if _g = _g.Group(key[:index]); _g == nil {
					continue
				}

				key = key[index+1:]
			}

			key = _g.fixOptName(key)
			switch v, err := _g.Parse(key, value); err {
			case nil:
				opts = append(opts, groupOptValue{Group: _g, Name: key, Value: v})
			case ErrNoOpt:
			default:
				return opts, err
			}
		}
	}

	return opts, nil
}

func (c *Config) loadMap(opts []groupOptValue, force bool) {
	for _, opt := range opts {
		// fmt.Println("------", opt.Group.name, opt.Name, opt.Value)
		if force || opt.Group.HasOptAndIsNotSet(opt.Name) {
			opt.Group.setOptWithLock(opt.Name, opt.Value)
		}
	}
}
