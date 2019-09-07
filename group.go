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

var (
	// ErrNoOpt is returned when no corresponding option.
	ErrNoOpt = fmt.Errorf("no option")

	// ErrFrozenOpt is returned when to set the value of the option but it's forzen.
	ErrFrozenOpt = fmt.Errorf("option is frozen")
)

// IsErrNoOpt whether reports the error is ErrNoOpt or not.
func IsErrNoOpt(err error) bool {
	if e, ok := err.(OptError); (ok && e.Err == ErrNoOpt) || err == ErrNoOpt {
		return true
	}
	return false
}

// IsErrFrozenOpt whether reports the error is ErrFrozenOpt or not.
func IsErrFrozenOpt(err error) bool {
	if e, ok := err.(OptError); ok && e.Err == ErrFrozenOpt {
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
		return fmt.Sprintf("[Config] invalid setting for '%s': %s", e.Opt, e.Err)
	}
	return fmt.Sprintf("[Config] invalid setting for '%s:%s': %s", e.Group, e.Opt, e.Err)
}

type groupOpt struct {
	opt    Opt
	value  interface{}
	watch  func(interface{})
	frozen bool
}

// OptGroup is the group of the options.
type OptGroup struct {
	lock sync.RWMutex
	conf *Config

	name   string // short name
	opts   map[string]*groupOpt
	alias  map[string]string
	frozen bool
}

func newOptGroup(conf *Config, name string) *OptGroup {
	return &OptGroup{
		conf:  conf,
		name:  name,
		opts:  make(map[string]*groupOpt, 16),
		alias: make(map[string]string, 16),
	}
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

// MustGroup is equal to g.Group(group), but panic if the group does not exist.
func (g *OptGroup) MustGroup(group string) *OptGroup {
	if _g := g.Group(group); _g != nil {
		return _g
	}
	panic(fmt.Errorf("no group '%s'", group))
}

// MustG is short for g.MustGroup(group).
func (g *OptGroup) MustG(group string) *OptGroup {
	return g.MustGroup(group)
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
	name = g.fixOptName(name)
	g.lock.RLock()
	gopt := g.opts[name]
	if gopt == nil {
		gopt = g.opts[g.alias[name]] // Check the alias
	}
	g.lock.RUnlock()

	if gopt != nil {
		return gopt.opt, true
	}
	return
}

// MustOpt is the same as Opt(name), but panic if the option does not exist.
func (g *OptGroup) MustOpt(name string) Opt {
	if opt, ok := g.Opt(name); ok {
		return opt
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// HasOpt reports whether there is an option named name in the current group.
func (g *OptGroup) HasOpt(name string) bool {
	name = g.fixOptName(name)
	g.lock.RLock()
	_, ok := g.opts[name]
	if !ok {
		if name = g.alias[name]; name != "" {
			_, ok = g.opts[name]
		}
	}
	g.lock.RUnlock()
	return ok
}

// OptIsSet reports whether the option named name has been set.
func (g *OptGroup) OptIsSet(name string) bool {
	name = g.fixOptName(name)

	var set bool
	g.lock.RLock()
	if opt, ok := g.opts[name]; ok {
		set = opt.value != nil
	} else if name, ok = g.alias[name]; ok {
		if opt, ok = g.opts[name]; ok {
			set = opt.value != nil
		}
	}
	g.lock.RUnlock()
	return set
}

// HasOptAndIsNotSet reports whether the option named name exists and
// hasn't been set, which is equal to `g.HasOpt(name) && !g.OptIsSet(name)`.
func (g *OptGroup) HasOptAndIsNotSet(name string) (yes bool) {
	name = g.fixOptName(name)
	g.lock.RLock()
	if opt, ok := g.opts[name]; ok {
		yes = opt.value == nil
	} else if name, ok = g.alias[name]; ok {
		if opt, ok = g.opts[name]; ok {
			yes = opt.value == nil
		}
	}
	g.lock.RUnlock()
	return
}

// SetOptAlias sets the alias of the option from old to new. So you can access
// the option named new by the old name.
//
// Notice: it does nothing if either is "".
func (g *OptGroup) SetOptAlias(old, new string) {
	old = g.fixOptName(old)
	new = g.fixOptName(new)
	if old == "" || new == "" {
		return
	}

	g.lock.Lock()
	for name, opt := range g.opts {
		if name == new {
			var exist bool
			for _, alias := range opt.opt.Aliases {
				if alias == old {
					exist = true
					break
				}
			}
			if !exist {
				opt.opt.Aliases = append(opt.opt.Aliases, old)
			}
		}
	}
	g.alias[old] = new
	g.lock.Unlock()
	debugf("[Config] Set the option alias from '%s' to '%s' in the group '%s'",
		old, new, g.name)
}

func (g *OptGroup) setOptWatch(name string, watch func(interface{})) {
	name = g.fixOptName(name)
	g.lock.Lock()
	if opt, ok := g.opts[name]; ok {
		opt.watch = watch
	} else if name, ok = g.alias[name]; ok {
		if opt, ok = g.opts[name]; ok {
			opt.watch = watch
		}
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

	if ok {
		debugf("[Config] Register the option '%s' into the group '%s'", opt.Name, g.name)
		g.conf.noticeOptRegister(g.name, []Opt{opt})

		for _, alias := range opt.Aliases {
			if alias != "" {
				g.SetOptAlias(alias, name)
			}
		}
	}

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
			debugf("[Config] Register the option '%s' into the group '%s'", opt.Name, g.name)
		}
		ok = true
	}
	g.lock.Unlock()

	if ok {
		g.conf.noticeOptRegister(g.name, opts)
		for _, opt := range opts {
			for _, alias := range opt.Aliases {
				if alias != "" {
					g.SetOptAlias(alias, opt.Name)
				}
			}
		}
	}

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

// FreezeGroup freezes the current group and disable its options to be set.
//
// If the current group has been frozen, it does nothing.
func (g *OptGroup) FreezeGroup() {
	g.lock.Lock()
	if !g.frozen {
		g.frozen = true
	}
	g.lock.Unlock()
}

// UnfreezeGroup unfreezes the current group and allows its options to be set.
//
// If the current group has been unfrozen, it does nothing.
func (g *OptGroup) UnfreezeGroup() {
	g.lock.Lock()
	if g.frozen {
		g.frozen = false
	}
	g.lock.Unlock()
}

// FreezeOpt freezes these options and disable them to be set.
//
// If the option does not exist has been frozen, it does nothing for it.
func (g *OptGroup) FreezeOpt(names ...string) {
	for i := range names {
		names[i] = g.fixOptName(names[i])
	}

	g.lock.Lock()
	for _, name := range names {
		gopt, ok := g.opts[name]
		if !ok {
			if alias, exist := g.alias[name]; exist {
				gopt, ok = g.opts[alias]
			}
		}
		if ok && !gopt.frozen {
			gopt.frozen = true
		}
	}
	g.lock.Unlock()
}

// UnfreezeOpt unfreeze these options and allows them to be set.
//
// If the option does not exist has been unfrozen, it does nothing for it.
func (g *OptGroup) UnfreezeOpt(names ...string) {
	for i := range names {
		names[i] = g.fixOptName(names[i])
	}

	g.lock.Lock()
	for _, name := range names {
		gopt, ok := g.opts[name]
		if !ok {
			if alias, exist := g.alias[name]; exist {
				gopt, ok = g.opts[alias]
			}
		}
		if ok && gopt.frozen {
			gopt.frozen = false
		}
	}
	g.lock.Unlock()
}

// GroupIsFrozen reports whether the current group named name is frozen.
func (g *OptGroup) GroupIsFrozen() (frozen bool) {
	g.lock.RLock()
	frozen = g.frozen
	g.lock.RUnlock()
	return
}

// OptIsFrozen reports whether the option named name is frozen.
//
// Return false if the option does not exist.
func (g *OptGroup) OptIsFrozen(name string) (frozen bool) {
	name = g.fixOptName(name)
	g.lock.RLock()
	frozen = g.optIsFrozen(name)
	g.lock.RUnlock()
	return
}

func (g *OptGroup) optIsFrozen(name string) bool {
	if g.frozen {
		return true
	} else if gopt, ok := g.opts[name]; ok {
		return gopt.frozen
	} else if name, ok = g.alias[name]; ok {
		if gopt, ok = g.opts[name]; ok {
			return gopt.frozen
		}
	}
	return false
}

func (g *OptGroup) parseOptValue(name string, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	opt, ok := g.opts[name]
	if !ok {
		if name, ok = g.alias[name]; !ok {
			return nil, ErrNoOpt
		} else if opt, ok = g.opts[name]; !ok {
			return nil, ErrNoOpt
		}
	}

	// Parse the option value
	v, err := opt.opt.Parser(value)
	if err != nil {
		return nil, NewOptError(g.name, opt.opt.Name, err, value)
	}

	// Validate the option value
	if err = opt.opt.validate(v); err != nil {
		return nil, NewOptError(g.name, opt.opt.Name, err, v)
	}

	return v, nil
}

func (g *OptGroup) setOptValue(name string, value interface{}) {
	opt := g.opts[name]
	if opt == nil {
		name = g.alias[name]
		opt = g.opts[name]
	}
	old := opt.value
	opt.value = value
	if old == nil {
		old = opt.opt.Default
	}
	g.conf.noticeOptChange(g.name, name, old, value, opt.opt.Observers)
	if opt.watch != nil {
		opt.watch(value)
	}

	if g.name == "" {
		debugf("[Config] Set [%s] to '%v'", name, value)
	} else {
		debugf("[Config] Set [%s:%s] to '%v'", g.name, name, value)
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
func (g *OptGroup) Set(name string, value interface{}) {
	if value == nil {
		return
	}

	name = g.fixOptName(name)
	g.lock.Lock()
	defer g.lock.Unlock()

	// Check whether the current group or the option is frozen.
	if g.optIsFrozen(name) {
		g.conf.handleError(NewOptError(g.Name(), name, ErrFrozenOpt, value))
		return
	}

	// Parse the option value
	value, err := g.parseOptValue(name, value)
	switch err {
	case nil:
		// Set the option value
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
	opt, ok := g.opts[name]
	if !ok {
		if name, ok = g.alias[name]; ok {
			opt, ok = g.opts[name]
		}
	}
	if ok {
		if value = opt.value; value == nil {
			value = opt.opt.Default
		}
	}
	g.lock.RUnlock()
	return
}

// Must is the same as Get(name), but panic if the option does not exist.
func (g *OptGroup) Must(name string) (value interface{}) {
	if value = g.Get(name); value == nil {
		panic(NewOptError(g.name, name, ErrNoOpt, nil))
	}
	return
}

// GetBool is the same as Get(name), but returns the bool value.
func (g *OptGroup) GetBool(name string) bool {
	v, _ := ToBool(g.Get(name))
	return v
}

// MustBool is the same as GetBool(name), but panic if the option does not exist.
func (g *OptGroup) MustBool(name string) bool {
	if value := g.Get(name); value != nil {
		v, _ := ToBool(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetInt is the same as Get(name), but returns the int value.
func (g *OptGroup) GetInt(name string) int {
	v, _ := ToInt(g.Get(name))
	return v
}

// MustInt is the same as GetInt(name), but panic if the option does not exist.
func (g *OptGroup) MustInt(name string) int {
	if value := g.Get(name); value != nil {
		v, _ := ToInt(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetInt32 is the same as Get(name), but returns the int32 value.
func (g *OptGroup) GetInt32(name string) int32 {
	v, _ := ToInt32(g.Get(name))
	return v
}

// MustInt32 is the same as GetInt32(name), but panic if the option does not exist.
func (g *OptGroup) MustInt32(name string) int32 {
	if value := g.Get(name); value != nil {
		v, _ := ToInt32(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetInt64 is the same as Get(name), but returns the int64 value.
func (g *OptGroup) GetInt64(name string) int64 {
	v, _ := ToInt64(g.Get(name))
	return v
}

// MustInt64 is the same as GetInt64(name), but panic if the option does not exist.
func (g *OptGroup) MustInt64(name string) int64 {
	if value := g.Get(name); value != nil {
		v, _ := ToInt64(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetUint is the same as Get(name), but returns the uint value.
func (g *OptGroup) GetUint(name string) uint {
	v, _ := ToUint(g.Get(name))
	return v
}

// MustUint is the same as GetUint(name), but panic if the option does not exist.
func (g *OptGroup) MustUint(name string) uint {
	if value := g.Get(name); value != nil {
		v, _ := ToUint(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetUint32 is the same as Get(name), but returns the uint32 value.
func (g *OptGroup) GetUint32(name string) uint32 {
	v, _ := ToUint32(g.Get(name))
	return v
}

// MustUint32 is the same as GetUint32(name), but panic if the option does not exist.
func (g *OptGroup) MustUint32(name string) uint32 {
	if value := g.Get(name); value != nil {
		v, _ := ToUint32(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetUint64 is the same as Get(name), but returns the uint64 value.
func (g *OptGroup) GetUint64(name string) uint64 {
	v, _ := ToUint64(g.Get(name))
	return v
}

// MustUint64 is the same as GetUint64(name), but panic if the option does not exist.
func (g *OptGroup) MustUint64(name string) uint64 {
	if value := g.Get(name); value != nil {
		v, _ := ToUint64(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetFloat64 is the same as Get(name), but returns the float64 value.
func (g *OptGroup) GetFloat64(name string) float64 {
	v, _ := ToFloat64(g.Get(name))
	return v
}

// MustFloat64 is the same as GetFloat64(name), but panic if the option does not exist.
func (g *OptGroup) MustFloat64(name string) float64 {
	if value := g.Get(name); value != nil {
		v, _ := ToFloat64(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetString is the same as Get(name), but returns the string value.
func (g *OptGroup) GetString(name string) string {
	v, _ := ToString(g.Get(name))
	return v
}

// MustString is the same as GetString(name), but panic if the option does not exist.
func (g *OptGroup) MustString(name string) string {
	if value := g.Get(name); value != nil {
		v, _ := ToString(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetDuration is the same as Get(name), but returns the time.Duration value.
func (g *OptGroup) GetDuration(name string) time.Duration {
	v, _ := ToDuration(g.Get(name))
	return v
}

// MustDuration is the same as GetDuration(name), but panic if the option does not exist.
func (g *OptGroup) MustDuration(name string) time.Duration {
	if value := g.Get(name); value != nil {
		v, _ := ToDuration(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetTime is the same as Get(name), but returns the time.Time value.
func (g *OptGroup) GetTime(name string) time.Time {
	v, _ := ToTime(g.Get(name))
	return v
}

// MustTime is the same as GetTime(name), but panic if the option does not exist.
func (g *OptGroup) MustTime(name string) time.Time {
	if value := g.Get(name); value != nil {
		v, _ := ToTime(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetIntSlice is the same as Get(name), but returns the []int value.
func (g *OptGroup) GetIntSlice(name string) []int {
	v, _ := ToIntSlice(g.Get(name))
	return v
}

// MustIntSlice is the same as GetIntSlice(name), but panic if the option does not exist.
func (g *OptGroup) MustIntSlice(name string) []int {
	if value := g.Get(name); value != nil {
		v, _ := ToIntSlice(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetUintSlice is the same as Get(name), but returns the []uint value.
func (g *OptGroup) GetUintSlice(name string) []uint {
	v, _ := ToUintSlice(g.Get(name))
	return v
}

// MustUintSlice is the same as GetUintSlice(name), but panic if the option does not exist.
func (g *OptGroup) MustUintSlice(name string) []uint {
	if value := g.Get(name); value != nil {
		v, _ := ToUintSlice(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetFloat64Slice is the same as Get(name), but returns the []float64 value.
func (g *OptGroup) GetFloat64Slice(name string) []float64 {
	v, _ := ToFloat64Slice(g.Get(name))
	return v
}

// MustFloat64Slice is the same as GetFloat64Slice(name), but panic if the option does not exist.
func (g *OptGroup) MustFloat64Slice(name string) []float64 {
	if value := g.Get(name); value != nil {
		v, _ := ToFloat64Slice(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetStringSlice is the same as Get(name), but returns the []string value.
func (g *OptGroup) GetStringSlice(name string) []string {
	v, _ := ToStringSlice(g.Get(name))
	return v
}

// MustStringSlice is the same as GetStringSlice(name), but panic if the option does not exist.
func (g *OptGroup) MustStringSlice(name string) []string {
	if value := g.Get(name); value != nil {
		v, _ := ToStringSlice(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}

// GetDurationSlice is the same as Get(name), but returns the []time.Duration value.
func (g *OptGroup) GetDurationSlice(name string) []time.Duration {
	v, _ := ToDurationSlice(g.Get(name))
	return v
}

// MustDurationSlice is the same as GetDurationSlice(name), but panic if the option does not exist.
func (g *OptGroup) MustDurationSlice(name string) []time.Duration {
	if value := g.Get(name); value != nil {
		v, _ := ToDurationSlice(value)
		return v
	}
	panic(NewOptError(g.name, name, ErrNoOpt, nil))
}
