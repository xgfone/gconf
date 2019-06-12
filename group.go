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
	"sort"
	"strings"
	"sync"
	"time"
)

// ErrNoOpt is returned when no corresponding option.
var ErrNoOpt = fmt.Errorf("no option")

// IsErrNoOpt whether reports the error is ErrNoOpt or not.
func IsErrNoOpt(err error) bool {
	if e, ok := err.(OptError); (ok && e.Err == ErrNoOpt) || err == ErrNoOpt {
		return true
	}
	return false
}

// OptError is used to represent an error about option.
type OptError struct {
	Group string
	Opt   string
	Err   error
	Value interface{}
}

// NewOptError returns a new OptError.
func NewOptError(group, opt string, err error, value interface{}) OptError {
	return OptError{Group: group, Opt: opt, Err: err, Value: value}
}

func (e OptError) Error() string {
	if e.Group == "" {
		return fmt.Sprintf("invalid value of option '%s': %s", e.Opt, e.Err)
	}
	return fmt.Sprintf("invalid value of option '%s:%s': %s", e.Group, e.Opt, e.Err)
}

type groupOpt struct {
	opt   Opt
	value interface{}
	watch func(interface{})
}

// OptGroup is the group of the options.
type OptGroup struct {
	lock sync.RWMutex
	conf *Config

	name string // short name
	opts map[string]*groupOpt
}

func newOptGroup(conf *Config, name string) *OptGroup {
	return &OptGroup{conf: conf, name: name, opts: make(map[string]*groupOpt, 16)}
}

// NewGroup returns a new sub-group with the name named `group`.
//
// Notice: If the sub-group has existed, it will return it instead.
func (g *OptGroup) NewGroup(group string) *OptGroup {
	return g.conf.newGroup(g.name, group)
}

// Config returns the Config that the current group belongs to.
func (g *OptGroup) Config() *Config {
	return g.conf
}

// Name returns the full name of the current group.
func (g *OptGroup) Name() string {
	return g.name
}

// Group returns the sub-group named group.
//
// It supports the cascaded group name, for example, the following ways are equal.
//
//     g.Group("group1.group2.group3")
//     g.Group("group1").Group("group2.group3")
//     g.Group("group1.group2").Group("group3")
//     g.Group("group1").Group("group2").Group("group3")
//
// Notice: if the group is "", it will return the current group.
func (g *OptGroup) Group(group string) *OptGroup {
	return g.conf.getGroup(g.name, group)
}

// G is the alias of Group.
func (g *OptGroup) G(group string) *OptGroup {
	return g.Group(group)
}

// AllSubGroups returns all the sub-groups of the current group.
//
// If the current group is the default, it will returns all groups,
// containing the default group.
//
// Notice: "group1.group2.group3" and "group1.group2" are the sub-group of "group1".
func (g *OptGroup) AllSubGroups() []*OptGroup {
	var gname string
	if name := g.Name(); name != "" {
		gname = name + g.conf.gsep
	} else {
		return g.conf.AllGroups()
	}

	allGroups := g.conf.AllGroups()
	groups := make([]*OptGroup, 0, len(allGroups))
	for _, group := range allGroups {
		if strings.HasPrefix(group.Name(), gname) {
			groups = append(groups, group)
		}
	}
	return groups
}

// AllOpts returns all the options in the current group.
func (g *OptGroup) AllOpts() []Opt {
	g.lock.RLock()
	opts := make([]Opt, len(g.opts))
	var index int
	for _, opt := range g.opts {
		opts[index] = opt.opt
		index++
	}
	g.lock.RUnlock()

	sort.Slice(opts, func(i, j int) bool { return opts[i].Name < opts[j].Name })
	return opts
}

func (g *OptGroup) fixOptName(name string) string {
	if name == "" {
		panic("the option name must not be empty")
	}
	return strings.Replace(name, "-", "_", -1)
}

// Opt returns the option named name.
func (g *OptGroup) Opt(name string) (opt Opt, exist bool) {
	g.lock.RLock()
	gopt := g.opts[g.fixOptName(name)]
	g.lock.RUnlock()

	if gopt != nil {
		return gopt.opt, true
	}
	return
}

// HasOpt reports whether there is an option named name in the current group.
func (g *OptGroup) HasOpt(name string) bool {
	g.lock.RLock()
	_, ok := g.opts[g.fixOptName(name)]
	g.lock.RUnlock()
	return ok
}

// IsSet reports whether the option named name has been set.
func (g *OptGroup) IsSet(name string) bool {
	var set bool
	g.lock.RLock()
	if opt, ok := g.opts[g.fixOptName(name)]; ok {
		set = opt.value != nil
	}
	g.lock.RUnlock()
	return set
}

// HasAndIsNotSet reports whether the option named name exists and
// hasn't been set, which is equal to `g.HasOpt(name) && !g.IsSet(name)`.
func (g *OptGroup) HasAndIsNotSet(name string) (yes bool) {
	g.lock.RLock()
	if opt, ok := g.opts[g.fixOptName(name)]; ok {
		yes = opt.value == nil
	}
	g.lock.RUnlock()
	return
}

func (g *OptGroup) setOptWatch(name string, watch func(interface{})) {
	g.lock.Lock()
	if opt, ok := g.opts[g.fixOptName(name)]; ok {
		opt.watch = watch
	}
	g.lock.Unlock()
}

func (g *OptGroup) registerOpt(opt Opt, force ...bool) (ok bool) {
	opt.check()
	if err := opt.validate(opt.Default); err != nil {
		panic(NewOptError(g.name, opt.Name, err, opt.Default))
	}

	name := g.fixOptName(opt.Name)
	g.lock.Lock()
	if _, exist := g.opts[name]; !exist || (len(force) > 0 && force[0]) {
		g.opts[name] = &groupOpt{opt: opt}
		ok = true
	}
	g.lock.Unlock()

	debugf("[Config] Register the option '%s' into the group '%s'\n", opt.Name, g.name)
	return
}

func (g *OptGroup) registerOpts(opts []Opt, force ...bool) (ok bool) {
	names := make([]string, len(opts))
	for i := range opts {
		opts[i].check()
		if err := opts[i].validate(opts[i].Default); err != nil {
			panic(NewOptError(g.name, opts[i].Name, err, opts[i].Default))
		}
		names[i] = g.fixOptName(opts[i].Name)
	}

	var exist bool
	g.lock.Lock()
	for _, name := range names {
		if _, ok := g.opts[name]; ok {
			exist = true
			break
		}
	}
	if !exist || (len(force) > 0 && force[0]) {
		for i, opt := range opts {
			g.opts[names[i]] = &groupOpt{opt: opt}
			debugf("[Config] Register the option '%s' into the group '%s'\n", opt.Name, g.name)
		}
		ok = true
	}
	g.lock.Unlock()

	return
}

// RegisterOpt registers an option and returns true.
//
// Notice: if the option has existed, it won't register it and return false.
// But you can set force to true to override it forcibly then to return true.
func (g *OptGroup) RegisterOpt(opt Opt, force ...bool) (ok bool) {
	return g.registerOpt(opt, force...)
}

// RegisterOpts registers a set of options and returns true.
//
// Notice: if a certain option has existed, it won't register them
// and return false. But you can set force to true to override it forcibly
// then to return true.
func (g *OptGroup) RegisterOpts(opts []Opt, force ...bool) (ok bool) {
	return g.registerOpts(opts, force...)
}

func (g *OptGroup) parseOptValue(name string, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	opt, ok := g.opts[name]
	if !ok {
		return nil, ErrNoOpt
	}

	// Parse the option value
	v, err := opt.opt.Parser(value)
	if err != nil {
		return nil, NewOptError(g.name, opt.opt.Name, err, value)
	}

	// Validate the option value
	if err = opt.opt.validate(value); err != nil {
		return nil, NewOptError(g.name, opt.opt.Name, err, v)
	}

	return v, nil
}

func (g *OptGroup) setOptValue(name string, value interface{}) {
	opt := g.opts[name]
	old := opt.value
	opt.value = value
	if old == nil {
		old = opt.opt.Default
	}
	g.conf.noticeOptChange(g.name, name, old, value)
	if opt.watch != nil {
		opt.watch(value)
	}

	if g.name == "" {
		debugf("[Config] Set [%s] to '%v'\n", name, value)
	} else {
		debugf("[Config] Set [%s:%s] to '%v'\n", g.name, name, value)
	}
}

// Parse the value of the option named name, which will call the parser and
// the validators of this option.
//
// If no the option named `name`, it will return ErrNoOpt.
// If there is an other error, it will return an OptError.
func (g *OptGroup) Parse(name string, value interface{}) (interface{}, error) {
	name = g.fixOptName(name)
	g.lock.RLock()
	value, err := g.parseOptValue(name, value)
	g.lock.RUnlock()
	return value, err
}

// Set parses and sets the option value named name in the current group to value.
//
// For the option name, the characters "-" and "_" are equal, that's, "abcd-efg"
// is equal to "abcd_efg".
//
// If there is not the option or the value is nil, it will ignore it.
//
// Notice: You should not call the method for the struct option and access them
// by the struct field, because we have no way to promise that it's goroutine-safe.
func (g *OptGroup) Set(name string, value interface{}) {
	if value == nil {
		return
	}

	name = g.fixOptName(name)
	g.lock.Lock()
	defer g.lock.Unlock()

	value, err := g.parseOptValue(name, value)
	switch err {
	case nil:
		g.setOptValue(name, value)
	case ErrNoOpt:
		g.conf.handleError(NewOptError(g.Name(), name, err, value))
	default:
		g.conf.handleError(err)
	}
}

// Get returns the value of the option named name.
//
// Return nil if this option does not exist.
func (g *OptGroup) Get(name string) (value interface{}) {
	name = g.fixOptName(name)
	g.lock.RLock()
	if opt, ok := g.opts[name]; ok {
		if value = opt.value; value == nil {
			value = opt.opt.Default
		}
	}
	g.lock.RUnlock()
	return
}

// GetBool is the same as Get(name), but returns the bool value.
func (g *OptGroup) GetBool(name string) bool {
	v, _ := ToBool(g.Get(name))
	return v
}

// GetInt is the same as Get(name), but returns the int value.
func (g *OptGroup) GetInt(name string) int {
	v, _ := ToInt(g.Get(name))
	return v
}

// GetInt32 is the same as Get(name), but returns the int32 value.
func (g *OptGroup) GetInt32(name string) int32 {
	v, _ := ToInt32(g.Get(name))
	return v
}

// GetInt64 is the same as Get(name), but returns the int64 value.
func (g *OptGroup) GetInt64(name string) int64 {
	v, _ := ToInt64(g.Get(name))
	return v
}

// GetUint is the same as Get(name), but returns the uint value.
func (g *OptGroup) GetUint(name string) uint {
	v, _ := ToUint(g.Get(name))
	return v
}

// GetUint32 is the same as Get(name), but returns the uint32 value.
func (g *OptGroup) GetUint32(name string) uint32 {
	v, _ := ToUint32(g.Get(name))
	return v
}

// GetUint64 is the same as Get(name), but returns the uint64 value.
func (g *OptGroup) GetUint64(name string) uint64 {
	v, _ := ToUint64(g.Get(name))
	return v
}

// GetFloat64 is the same as Get(name), but returns the float64 value.
func (g *OptGroup) GetFloat64(name string) float64 {
	v, _ := ToFloat64(g.Get(name))
	return v
}

// GetString is the same as Get(name), but returns the string value.
func (g *OptGroup) GetString(name string) string {
	v, _ := ToString(g.Get(name))
	return v
}

// GetDuration is the same as Get(name), but returns the time.Duration value.
func (g *OptGroup) GetDuration(name string) time.Duration {
	v, _ := ToDuration(g.Get(name))
	return v
}

// GetTime is the same as Get(name), but returns the time.Time value.
func (g *OptGroup) GetTime(name string) time.Time {
	v, _ := ToTime(g.Get(name))
	return v
}

// GetIntSlice is the same as Get(name), but returns the []int value.
func (g *OptGroup) GetIntSlice(name string) []int {
	v, _ := ToIntSlice(g.Get(name))
	return v
}

// GetUintSlice is the same as Get(name), but returns the []uint value.
func (g *OptGroup) GetUintSlice(name string) []uint {
	v, _ := ToUintSlice(g.Get(name))
	return v
}

// GetFloat64Slice is the same as Get(name), but returns the []float64 value.
func (g *OptGroup) GetFloat64Slice(name string) []float64 {
	v, _ := ToFloat64Slice(g.Get(name))
	return v
}

// GetStringSlice is the same as Get(name), but returns the []string value.
func (g *OptGroup) GetStringSlice(name string) []string {
	v, _ := ToStringSlice(g.Get(name))
	return v
}

// GetDurationSlice is the same as Get(name), but returns the []time.Duration value.
func (g *OptGroup) GetDurationSlice(name string) []time.Duration {
	v, _ := ToDurationSlice(g.Get(name))
	return v
}
