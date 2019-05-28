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
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// DefaultGroupName is the name of the default group.
const DefaultGroupName = "DEFAULT"

type option struct {
	lock sync.RWMutex

	opt   Opt
	unmut int
	isCli bool
	value interface{}
}

// OptGroup is the group of the options.
type OptGroup struct {
	cmd  *Command
	conf *Config
	lock sync.RWMutex

	name   string // short name
	fname  string
	paths  []string
	opts   map[string]*option
	fields map[string]reflect.Value
	groups map[string]*OptGroup
}

// newOptGroup returns a new OptGroup.
func newOptGroup(conf *Config, cmd *Command, name string, parents ...string) *OptGroup {
	return newOptGroup2(true, conf, cmd, name, parents...)
}

func newOptGroup2(notice bool, conf *Config, cmd *Command, name string, parents ...string) *OptGroup {
	if conf == nil {
		panic("Config must not be nil")
	} else if name = strings.Trim(name, conf.gsep); name == "" {
		panic("the group name must not be empty")
	}

	paths := append([]string{}, parents...)
	paths = append(paths, name)

	group := &OptGroup{
		cmd:   cmd,
		conf:  conf,
		name:  name,
		paths: paths,
		fname: conf.mergePaths(paths),

		opts:   make(map[string]*option, 16),
		fields: make(map[string]reflect.Value),
		groups: make(map[string]*OptGroup),
	}

	if notice {
		conf.noticeNewGroup(group)
	}

	conf.Debugf("Creating the group '%s'", group.fname)
	return group
}

func (g *OptGroup) fixOptName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

//////////////////////////////////////////////////////////////////////////////
/// MetaData

// Name returns the name of the current group.
func (g *OptGroup) Name() string {
	return g.name
}

// FullName returns the full name of the current group.
func (g *OptGroup) FullName() string {
	return g.fname
}

// OnlyGroupName returns the full name of the current group, but not contain
// the prefix of the command that the current group belongs to if exists.
//
// Return "" if the current group is the default group of a command.
func (g *OptGroup) OnlyGroupName() string {
	if g.cmd == nil {
		return g.fname
	} else if g == g.cmd.OptGroup {
		return ""
	}
	return strings.TrimPrefix(g.fname, g.cmd.OptGroup.fname+g.conf.gsep)
}

// Config returns the Config that the current group belongs to.
func (g *OptGroup) Config() *Config {
	return g.conf
}

// Command returns the Command that the current group belongs to.
//
// Return nil if the current group isn't belong to a certain command.
func (g *OptGroup) Command() *Command {
	return g.cmd
}

//////////////////////////////////////////////////////////////////////////////
/// Group

func (g *OptGroup) newGroup(name string, paths []string) (group *OptGroup) {
	if name == "" {
		panic("the group name must not be empty")
	}
	g.conf.panicIsParsed(true)

	if group = g.groups[name]; group == nil {
		group = newOptGroup(g.conf, g.cmd, name, paths...)
		g.groups[name] = group
	}
	return
}

// NewGroup returns a sub-group named name.
//
// Notice:
//   1. If the sub-group has existed, it will the old.
//   2. The command name should only contain the characters, [-_a-zA-Z0-9].
func (g *OptGroup) NewGroup(name string) *OptGroup {
	return g.newGroup(name, g.paths)
}

// IsConfigDefaultGroup reports whether the current group is the default of Config.
func (g *OptGroup) IsConfigDefaultGroup() bool {
	if g == g.conf.OptGroup {
		return true
	}
	return false
}

// IsCommandDefaultGroup reports whether the current group is the default of Command.
func (g *OptGroup) IsCommandDefaultGroup() bool {
	if g.cmd != nil && g == g.cmd.OptGroup {
		return true
	}
	return false
}

// IsDefaultGroup reports whether the current group is the default of Config or Command.
func (g *OptGroup) IsDefaultGroup() bool {
	return g.IsConfigDefaultGroup() || g.IsCommandDefaultGroup()
}

// IsConfigGroup reports whether the current group belongs to Config not Command.
//
// Notice: IsConfigGroup() == !IsCommandGroup().
func (g *OptGroup) IsConfigGroup() bool {
	return g.cmd == nil
}

// IsCommandGroup reports whether the current group belongs to Command not Config.
//
// Notice: IsCommandGroup() == !IsConfigGroup().
func (g *OptGroup) IsCommandGroup() bool {
	return !g.IsConfigGroup()
}

// Groups returns all the sub-groups.
func (g *OptGroup) Groups() []*OptGroup {
	groups := make([]*OptGroup, 0, len(g.groups))
	for _, group := range g.groups {
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].name < groups[j].name })
	return groups
}

// HasGroup reports whether the group contains the sub-group named 'name'.
func (g *OptGroup) HasGroup(name string) bool {
	_, ok := g.groups[name]
	return ok
}

// Group returns the sub-group by the name.
//
// Return itself if name is "". So it can be used to get the config or command
// itself as the group.
//
// Return nil if the sub-group does not exist.
func (g *OptGroup) Group(name string) *OptGroup {
	return g.conf.getGroup(g.fname, name)
}

// G is the short for g.Group(group).
func (g *OptGroup) G(group string) *OptGroup {
	return g.Group(group)
}

//////////////////////////////////////////////////////////////////////////////
/// Options

// Priority returns the priority of the option named name.
//
// DEPRECATED!!! it always returns 0.
func (g *OptGroup) Priority(name string) int {
	return 0
}

// HasOpt reports whether the group contains the option named 'name'.
func (g *OptGroup) HasOpt(name string) bool {
	_, ok := g.opts[g.fixOptName(name)]
	return ok
}

// LockOpt locks the option and forbid it to be updated.
//
// If the option does not exist, it does nothing and returns false.
// Or return true.
//
// Notice: you should call UnlockOpt once only if calling LockOpt once.
func (g *OptGroup) LockOpt(name string) (ok bool) {
	if opt, ok := g.opts[g.fixOptName(name)]; ok {
		opt.lock.Lock()
		opt.unmut++
		opt.lock.Unlock()
		return true
	}
	return false
}

// UnlockOpt unlocks the option and allow it to be updated.
//
// If the option does not exist or is unlocked, it does nothing and returns false.
// Or return true.
func (g *OptGroup) UnlockOpt(name string) (ok bool) {
	if opt, ok := g.opts[g.fixOptName(name)]; ok {
		opt.lock.Lock()
		if opt.unmut > 0 {
			opt.unmut--
		}
		opt.lock.Unlock()
		return true
	}
	return false
}

// Opt returns the option named name.
//
// Return nil if the option does not exist.
func (g *OptGroup) Opt(name string) Opt {
	if option, ok := g.opts[g.fixOptName(name)]; ok {
		return option.opt
	}
	return nil
}

// AllOpts returns all the registered options, including the CLI options.
func (g *OptGroup) AllOpts() []Opt {
	opts := make([]Opt, 0, len(g.opts))
	for _, opt := range g.opts {
		opts = append(opts, opt.opt)
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Name() < opts[j].Name() })
	return opts
}

// NotCliOpts returns all the registered options, except the CLI options.
func (g *OptGroup) NotCliOpts() []Opt {
	opts := make([]Opt, 0, len(g.opts))
	for _, opt := range g.opts {
		if !opt.isCli {
			opts = append(opts, opt.opt)
		}
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Name() < opts[j].Name() })
	return opts
}

// CliOpts returns all the registered CLI options, except the non-CLI options.
func (g *OptGroup) CliOpts() []Opt {
	opts := make([]Opt, 0, len(g.opts))
	for _, opt := range g.opts {
		if opt.isCli {
			opts = append(opts, opt.opt)
		}
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Name() < opts[j].Name() })
	return opts
}

//////////////////////////////////////////////////////////////////////////////
/// Register Options

// RegisterOpts registers a set of options into the current group.
func (g *OptGroup) RegisterOpts(opts []Opt) *OptGroup {
	for _, opt := range opts {
		g.RegisterOpt(opt)
	}
	return g
}

// RegisterOpt registers the option into the current group.
//
// the characters "-" and "_" in the option name are equal, that's, "abcd-efg"
// is equal to "abcd_efg".
func (g *OptGroup) RegisterOpt(opt Opt) *OptGroup {
	return g.registerOpt(false, opt)
}

func (g *OptGroup) registerOpt(cli bool, opt Opt) *OptGroup {
	g.conf.panicIsParsed(true)
	if opt == nil {
		panic("the option must not be nil")
	}

	name := g.fixOptName(opt.Name())
	if _, ok := g.opts[name]; ok {
		if g.conf.reregister {
			g.conf.Debugf("WARNING: Ingore to reregister the option '%s' into the group '%s'",
				opt.Name(), g.fname)
			return g
		}
		panic(fmt.Errorf("the option '%s' has been registered into the group '%s'", opt.Name(), g.fname))
	}

	g.opts[name] = &option{isCli: cli, opt: opt}
	g.conf.Debugf("Register group=%s, option=%s, cli=%v", g.fname, opt.Name(), cli)
	return g
}

///////////////////////////////////////////////////////////////////////////////
/// Set the option value.

func (g *OptGroup) parseOptValue(name string, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	opt, ok := g.opts[name]
	if !ok {
		return nil, nil
	}

	var err error
	if value, err = opt.opt.Parse(value); err != nil {
		return nil, err
	}

	// The option has a validator.
	if v, ok := opt.opt.(Validator); ok {
		if err = v.Validate(g.fname, name, value); err != nil {
			return nil, err
		}
	}

	// The option has a validator chain.
	if vc, ok := opt.opt.(ValidatorChainOpt); ok {
		vs := vc.GetValidators()
		if len(vs) > 0 {
			for _, v := range vs {
				if err = v.Validate(g.fname, name, value); err != nil {
					return nil, err
				}
			}
		}
	}

	return value, nil
}

func (g *OptGroup) _setOptValue(name string, value interface{}) {
	var ok bool
	var old interface{}

	func() {
		option := g.opts[name]
		option.lock.Lock()
		defer option.lock.Unlock()

		if option.unmut > 0 {
			g.conf.Debugf("Ignore the option [%s]:[%s] because the option is locked", g.FullName(), name)
			return
		}

		ok = true
		old = option.value
		option.value = value

		if field, ok := g.fields[name]; ok {
			field.Set(reflect.ValueOf(value))
		}
	}()

	if ok {
		g.conf.watchChangedOption(g, name, old, value)
	}
}

// ParseOptValue parses the value of the option named name.
//
// If no the option named `name`, it will return (nil, nil).
func (g *OptGroup) ParseOptValue(name string, value interface{}) (interface{}, error) {
	g.conf.panicIsParsed(false)
	return g.parseOptValue(g.fixOptName(name), value)
}

func (g *OptGroup) setOptValue(name string, value interface{}) (err error) {
	name = g.fixOptName(name)
	if value, err = g.parseOptValue(name, value); err == nil && value != nil {
		g._setOptValue(name, value)
	}
	return
}

// UpdateOptValue parses and sets the value of the option in the current group,
// which is goroutine-safe.
//
// For the option name, the characters "-" and "_" are equal, that's, "abcd-efg"
// is equal to "abcd_efg".
//
// Notice: You cannot call UpdateOptValue() for the struct option and access them
// by the struct field, because we have no way to promise that it's goroutine-safe.
func (g *OptGroup) UpdateOptValue(name string, value interface{}) error {
	g.conf.panicIsParsed(false)
	return g.setOptValue(name, value)
}

// SetOptValue is equal to UpdateOptValue(name, value), which is deprecated.
func (g *OptGroup) SetOptValue(priority int, name string, value interface{}) error {
	return g.UpdateOptValue(name, value)
}

// CheckRequiredOption checks whether the required option has no value or a ZORE value.
func (g *OptGroup) CheckRequiredOption() (err error) {
	for name, option := range g.opts {
		if option.value != nil {
			continue
		}

		if g.conf.required {
			return fmt.Errorf("the option '%s' in the group '%s' has no value", name, g.fname)
		}
	}

	return nil
}

func (g *OptGroup) initAllOpts() (err error) {
	for name, option := range g.opts {
		// Set the default value if exists.
		if v := option.opt.Default(); v != nil {
			if err = g.setOptValue(name, v); err != nil {
				return
			}
			continue
		}

		// Set the zero value if it is enabled.
		if g.conf.zero {
			if v := option.opt.Zero(); v != nil {
				if err = g.setOptValue(name, v); err != nil {
					return
				}
			}
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Get the value from the current group.

// Value returns the value of the option.
//
// Return nil if the option does not exist, too.
func (g *OptGroup) Value(name string) (v interface{}) {
	if option := g.opts[g.fixOptName(name)]; option != nil {
		option.lock.RLock()
		v = option.value
		option.lock.RUnlock()
	}
	return
}

// V is the short for g.Value(name).
func (g *OptGroup) V(name string) interface{} {
	return g.Value(name)
}

func (g *OptGroup) getValue(name string, _type optType) (interface{}, error) {
	opt := g.Value(name)
	if opt == nil {
		if g.HasOpt(name) {
			return nil, fmt.Errorf("the option '%s' in the group '%s' has no value", name, g.fname)
		}
		return nil, fmt.Errorf("the group '%s' has no option '%s'", g.fname, name)
	}

	switch _type {
	case boolType:
		if v, ok := opt.(bool); ok {
			return v, nil
		}
	case stringType:
		if v, ok := opt.(string); ok {
			return v, nil
		}
	case intType:
		if v, ok := opt.(int); ok {
			return v, nil
		}
	case int8Type:
		if v, ok := opt.(int8); ok {
			return v, nil
		}
	case int16Type:
		if v, ok := opt.(int16); ok {
			return v, nil
		}
	case int32Type:
		if v, ok := opt.(int32); ok {
			return v, nil
		}
	case int64Type:
		if v, ok := opt.(int64); ok {
			return v, nil
		}
	case uintType:
		if v, ok := opt.(uint); ok {
			return v, nil
		}
	case uint8Type:
		if v, ok := opt.(uint8); ok {
			return v, nil
		}
	case uint16Type:
		if v, ok := opt.(uint16); ok {
			return v, nil
		}
	case uint32Type:
		if v, ok := opt.(uint32); ok {
			return v, nil
		}
	case uint64Type:
		if v, ok := opt.(uint64); ok {
			return v, nil
		}
	case float32Type:
		if v, ok := opt.(float32); ok {
			return v, nil
		}
	case float64Type:
		if v, ok := opt.(float64); ok {
			return v, nil
		}
	case durationType:
		if v, ok := opt.(time.Duration); ok {
			return v, nil
		}
	case timeType:
		if v, ok := opt.(time.Time); ok {
			return v, nil
		}
	case stringsType:
		if v, ok := opt.([]string); ok {
			return v, nil
		}
	case intsType:
		if v, ok := opt.([]int); ok {
			return v, nil
		}
	case int64sType:
		if v, ok := opt.([]int64); ok {
			return v, nil
		}
	case uintsType:
		if v, ok := opt.([]uint); ok {
			return v, nil
		}
	case uint64sType:
		if v, ok := opt.([]uint64); ok {
			return v, nil
		}
	case float64sType:
		if v, ok := opt.([]float64); ok {
			return v, nil
		}
	case durationsType:
		if v, ok := opt.([]time.Duration); ok {
			return v, nil
		}
	case timesType:
		if v, ok := opt.([]time.Time); ok {
			return v, nil
		}
	default:
		return nil, fmt.Errorf("don't support the type '%s'", _type)
	}
	return nil, fmt.Errorf("the option '%s' in the group '%s' is not the type '%s'",
		name, g.fname, _type)
}

// BoolE returns the option value, the type of which is bool.
//
// Return an error if no the option or the type of the option isn't bool.
func (g *OptGroup) BoolE(name string) (bool, error) {
	v, err := g.getValue(name, boolType)
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}

// BoolD is the same as BoolE, but returns the default if there is an error.
func (g *OptGroup) BoolD(name string, _default bool) bool {
	if value, err := g.BoolE(name); err == nil {
		return value
	}
	return _default
}

// Bool is the same as BoolE, but panic if there is an error.
func (g *OptGroup) Bool(name string) bool {
	value, err := g.BoolE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// StringE returns the option value, the type of which is string.
//
// Return an error if no the option or the type of the option isn't string.
func (g *OptGroup) StringE(name string) (string, error) {
	v, err := g.getValue(name, stringType)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

// StringD is the same as StringE, but returns the default if there is an error.
func (g *OptGroup) StringD(name, _default string) string {
	if value, err := g.StringE(name); err == nil {
		return value
	}
	return _default
}

// String is the same as StringE, but panic if there is an error.
func (g *OptGroup) String(name string) string {
	value, err := g.StringE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// IntE returns the option value, the type of which is int.
//
// Return an error if no the option or the type of the option isn't int.
func (g *OptGroup) IntE(name string) (int, error) {
	v, err := g.getValue(name, intType)
	if err != nil {
		return 0, err
	}
	return v.(int), nil
}

// IntD is the same as IntE, but returns the default if there is an error.
func (g *OptGroup) IntD(name string, _default int) int {
	if value, err := g.IntE(name); err == nil {
		return value
	}
	return _default
}

// Int is the same as IntE, but panic if there is an error.
func (g *OptGroup) Int(name string) int {
	value, err := g.IntE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Int8E returns the option value, the type of which is int8.
//
// Return an error if no the option or the type of the option isn't int8.
func (g *OptGroup) Int8E(name string) (int8, error) {
	v, err := g.getValue(name, int8Type)
	if err != nil {
		return 0, err
	}
	return v.(int8), nil
}

// Int8D is the same as Int8E, but returns the default if there is an error.
func (g *OptGroup) Int8D(name string, _default int8) int8 {
	if value, err := g.Int8E(name); err == nil {
		return value
	}
	return _default
}

// Int8 is the same as Int8E, but panic if there is an error.
func (g *OptGroup) Int8(name string) int8 {
	value, err := g.Int8E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Int16E returns the option value, the type of which is int16.
//
// Return an error if no the option or the type of the option isn't int16.
func (g *OptGroup) Int16E(name string) (int16, error) {
	v, err := g.getValue(name, int16Type)
	if err != nil {
		return 0, err
	}
	return v.(int16), nil
}

// Int16D is the same as Int16E, but returns the default if there is an error.
func (g *OptGroup) Int16D(name string, _default int16) int16 {
	if value, err := g.Int16E(name); err == nil {
		return value
	}
	return _default
}

// Int16 is the same as Int16E, but panic if there is an error.
func (g *OptGroup) Int16(name string) int16 {
	value, err := g.Int16E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Int32E returns the option value, the type of which is int32.
//
// Return an error if no the option or the type of the option isn't int32.
func (g *OptGroup) Int32E(name string) (int32, error) {
	v, err := g.getValue(name, int32Type)
	if err != nil {
		return 0, err
	}
	return v.(int32), nil
}

// Int32D is the same as Int32E, but returns the default if there is an error.
func (g *OptGroup) Int32D(name string, _default int32) int32 {
	if value, err := g.Int32E(name); err == nil {
		return value
	}
	return _default
}

// Int32 is the same as Int32E, but panic if there is an error.
func (g *OptGroup) Int32(name string) int32 {
	value, err := g.Int32E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Int64E returns the option value, the type of which is int64.
//
// Return an error if no the option or the type of the option isn't int64.
func (g *OptGroup) Int64E(name string) (int64, error) {
	v, err := g.getValue(name, int64Type)
	if err != nil {
		return 0, err
	}
	return v.(int64), nil
}

// Int64D is the same as Int64E, but returns the default if there is an error.
func (g *OptGroup) Int64D(name string, _default int64) int64 {
	if value, err := g.Int64E(name); err == nil {
		return value
	}
	return _default
}

// Int64 is the same as Int64E, but panic if there is an error.
func (g *OptGroup) Int64(name string) int64 {
	value, err := g.Int64E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// UintE returns the option value, the type of which is uint.
//
// Return an error if no the option or the type of the option isn't uint.
func (g *OptGroup) UintE(name string) (uint, error) {
	v, err := g.getValue(name, uintType)
	if err != nil {
		return 0, err
	}
	return v.(uint), nil
}

// UintD is the same as UintE, but returns the default if there is an error.
func (g *OptGroup) UintD(name string, _default uint) uint {
	if value, err := g.UintE(name); err == nil {
		return value
	}
	return _default
}

// Uint is the same as UintE, but panic if there is an error.
func (g *OptGroup) Uint(name string) uint {
	value, err := g.UintE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint8E returns the option value, the type of which is uint8.
//
// Return an error if no the option or the type of the option isn't uint8.
func (g *OptGroup) Uint8E(name string) (uint8, error) {
	v, err := g.getValue(name, uint8Type)
	if err != nil {
		return 0, err
	}
	return v.(uint8), nil
}

// Uint8D is the same as Uint8E, but returns the default if there is an error.
func (g *OptGroup) Uint8D(name string, _default uint8) uint8 {
	if value, err := g.Uint8E(name); err == nil {
		return value
	}
	return _default
}

// Uint8 is the same as Uint8E, but panic if there is an error.
func (g *OptGroup) Uint8(name string) uint8 {
	value, err := g.Uint8E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint16E returns the option value, the type of which is uint16.
//
// Return an error if no the option or the type of the option isn't uint16.
func (g *OptGroup) Uint16E(name string) (uint16, error) {
	v, err := g.getValue(name, uint16Type)
	if err != nil {
		return 0, err
	}
	return v.(uint16), nil
}

// Uint16D is the same as Uint16E, but returns the default if there is an error.
func (g *OptGroup) Uint16D(name string, _default uint16) uint16 {
	if value, err := g.Uint16E(name); err == nil {
		return value
	}
	return _default
}

// Uint16 is the same as Uint16E, but panic if there is an error.
func (g *OptGroup) Uint16(name string) uint16 {
	value, err := g.Uint16E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint32E returns the option value, the type of which is uint32.
//
// Return an error if no the option or the type of the option isn't uint32.
func (g *OptGroup) Uint32E(name string) (uint32, error) {
	v, err := g.getValue(name, uint32Type)
	if err != nil {
		return 0, err
	}
	return v.(uint32), nil
}

// Uint32D is the same as Uint32E, but returns the default if there is an error.
func (g *OptGroup) Uint32D(name string, _default uint32) uint32 {
	if value, err := g.Uint32E(name); err == nil {
		return value
	}
	return _default
}

// Uint32 is the same as Uint32E, but panic if there is an error.
func (g *OptGroup) Uint32(name string) uint32 {
	value, err := g.Uint32E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint64E returns the option value, the type of which is uint64.
//
// Return an error if no the option or the type of the option isn't uint64.
func (g *OptGroup) Uint64E(name string) (uint64, error) {
	v, err := g.getValue(name, uint64Type)
	if err != nil {
		return 0, err
	}
	return v.(uint64), nil
}

// Uint64D is the same as Uint64E, but returns the default if there is an error.
func (g *OptGroup) Uint64D(name string, _default uint64) uint64 {
	if value, err := g.Uint64E(name); err == nil {
		return value
	}
	return _default
}

// Uint64 is the same as Uint64E, but panic if there is an error.
func (g *OptGroup) Uint64(name string) uint64 {
	value, err := g.Uint64E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Float32E returns the option value, the type of which is float32.
//
// Return an error if no the option or the type of the option isn't float32.
func (g *OptGroup) Float32E(name string) (float32, error) {
	v, err := g.getValue(name, float32Type)
	if err != nil {
		return 0, err
	}
	return v.(float32), nil
}

// Float32D is the same as Float32E, but returns the default value if there is
// an error.
func (g *OptGroup) Float32D(name string, _default float32) float32 {
	if value, err := g.Float32E(name); err == nil {
		return value
	}
	return _default
}

// Float32 is the same as Float32E, but panic if there is an error.
func (g *OptGroup) Float32(name string) float32 {
	value, err := g.Float32E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Float64E returns the option value, the type of which is float64.
//
// Return an error if no the option or the type of the option isn't float64.
func (g *OptGroup) Float64E(name string) (float64, error) {
	v, err := g.getValue(name, float64Type)
	if err != nil {
		return 0, err
	}
	return v.(float64), nil
}

// Float64D is the same as Float64E, but returns the default value if there is
// an error.
func (g *OptGroup) Float64D(name string, _default float64) float64 {
	if value, err := g.Float64E(name); err == nil {
		return value
	}
	return _default
}

// Float64 is the same as Float64E, but panic if there is an error.
func (g *OptGroup) Float64(name string) float64 {
	value, err := g.Float64E(name)
	if err != nil {
		panic(err)
	}
	return value
}

// DurationE returns the option value, the type of which is time.Duration.
//
// Return an error if no the option or the type of the option isn't time.Duration.
func (g *OptGroup) DurationE(name string) (time.Duration, error) {
	v, err := g.getValue(name, durationType)
	if err != nil {
		return 0, err
	}
	return v.(time.Duration), nil
}

// DurationD is the same as DurationE, but returns the default value if there is
// an error.
func (g *OptGroup) DurationD(name string, _default time.Duration) time.Duration {
	if value, err := g.DurationE(name); err == nil {
		return value
	}
	return _default
}

// Duration is the same as DurationE, but panic if there is an error.
func (g *OptGroup) Duration(name string) time.Duration {
	value, err := g.DurationE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// TimeE returns the option value, the type of which is time.Time.
//
// Return an error if no the option or the type of the option isn't time.Time.
func (g *OptGroup) TimeE(name string) (time.Time, error) {
	v, err := g.getValue(name, timeType)
	if err != nil {
		return time.Time{}, err
	}
	return v.(time.Time), nil
}

// TimeD is the same as TimeE, but returns the default value if there is
// an error.
func (g *OptGroup) TimeD(name string, _default time.Time) time.Time {
	if value, err := g.TimeE(name); err == nil {
		return value
	}
	return _default
}

// Time is the same as TimeE, but panic if there is an error.
func (g *OptGroup) Time(name string) time.Time {
	value, err := g.TimeE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// StringsE returns the option value, the type of which is []string.
//
// Return an error if no the option or the type of the option isn't []string.
func (g *OptGroup) StringsE(name string) ([]string, error) {
	v, err := g.getValue(name, stringsType)
	if err != nil {
		return nil, err
	}
	return v.([]string), nil
}

// StringsD is the same as StringsE, but returns the default value if there is
// an error.
func (g *OptGroup) StringsD(name string, _default []string) []string {
	if value, err := g.StringsE(name); err == nil {
		return value
	}
	return _default
}

// Strings is the same as StringsE, but panic if there is an error.
func (g *OptGroup) Strings(name string) []string {
	value, err := g.StringsE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// IntsE returns the option value, the type of which is []int.
//
// Return an error if no the option or the type of the option isn't []int.
func (g *OptGroup) IntsE(name string) ([]int, error) {
	v, err := g.getValue(name, intsType)
	if err != nil {
		return nil, err
	}
	return v.([]int), nil
}

// IntsD is the same as IntsE, but returns the default value if there is
// an error.
func (g *OptGroup) IntsD(name string, _default []int) []int {
	if value, err := g.IntsE(name); err == nil {
		return value
	}
	return _default
}

// Ints is the same as IntsE, but panic if there is an error.
func (g *OptGroup) Ints(name string) []int {
	value, err := g.IntsE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Int64sE returns the option value, the type of which is []int64.
//
// Return an error if no the option or the type of the option isn't []int64.
func (g *OptGroup) Int64sE(name string) ([]int64, error) {
	v, err := g.getValue(name, int64sType)
	if err != nil {
		return nil, err
	}
	return v.([]int64), nil
}

// Int64sD is the same as Int64sE, but returns the default value if there is
// an error.
func (g *OptGroup) Int64sD(name string, _default []int64) []int64 {
	if value, err := g.Int64sE(name); err == nil {
		return value
	}
	return _default
}

// Int64s is the same as Int64s, but panic if there is an error.
func (g *OptGroup) Int64s(name string) []int64 {
	value, err := g.Int64sE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// UintsE returns the option value, the type of which is []uint.
//
// Return an error if no the option or the type of the option isn't []uint.
func (g *OptGroup) UintsE(name string) ([]uint, error) {
	v, err := g.getValue(name, uintsType)
	if err != nil {
		return nil, err
	}
	return v.([]uint), nil
}

// UintsD is the same as UintsE, but returns the default value if there is
// an error.
func (g *OptGroup) UintsD(name string, _default []uint) []uint {
	if value, err := g.UintsE(name); err == nil {
		return value
	}
	return _default
}

// Uints is the same as UintsE, but panic if there is an error.
func (g *OptGroup) Uints(name string) []uint {
	value, err := g.UintsE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint64sE returns the option value, the type of which is []uint64.
//
// Return an error if no the option or the type of the option isn't []uint64.
func (g *OptGroup) Uint64sE(name string) ([]uint64, error) {
	v, err := g.getValue(name, uint64sType)
	if err != nil {
		return nil, err
	}
	return v.([]uint64), nil
}

// Uint64sD is the same as Uint64sE, but returns the default value if there is
// an error.
func (g *OptGroup) Uint64sD(name string, _default []uint64) []uint64 {
	if value, err := g.Uint64sE(name); err == nil {
		return value
	}
	return _default
}

// Uint64s is the same as Uint64sE, but panic if there is an error.
func (g *OptGroup) Uint64s(name string) []uint64 {
	value, err := g.Uint64sE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// Float64sE returns the option value, the type of which is []float64.
//
// Return an error if no the option or the type of the option isn't []float64.
func (g *OptGroup) Float64sE(name string) ([]float64, error) {
	v, err := g.getValue(name, float64sType)
	if err != nil {
		return nil, err
	}
	return v.([]float64), nil
}

// Float64sD is the same as Float64sE, but returns the default value if there is
// an error.
func (g *OptGroup) Float64sD(name string, _default []float64) []float64 {
	if value, err := g.Float64sE(name); err == nil {
		return value
	}
	return _default
}

// Float64s is the same as Float64sE, but panic if there is an error.
func (g *OptGroup) Float64s(name string) []float64 {
	value, err := g.Float64sE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// DurationsE returns the option value, the type of which is []time.Duration.
//
// Return an error if no the option or the type of the option isn't []time.Duration.
func (g *OptGroup) DurationsE(name string) ([]time.Duration, error) {
	v, err := g.getValue(name, durationsType)
	if err != nil {
		return nil, err
	}
	return v.([]time.Duration), nil
}

// DurationsD is the same as DurationsE, but returns the default value if there is
// an error.
func (g *OptGroup) DurationsD(name string, _default []time.Duration) []time.Duration {
	if value, err := g.DurationsE(name); err == nil {
		return value
	}
	return _default
}

// Durations is the same as DurationsE, but panic if there is an error.
func (g *OptGroup) Durations(name string) []time.Duration {
	value, err := g.DurationsE(name)
	if err != nil {
		panic(err)
	}
	return value
}

// TimesE returns the option value, the type of which is []time.Time.
//
// Return an error if no the option or the type of the option isn't []time.Time.
func (g *OptGroup) TimesE(name string) ([]time.Time, error) {
	v, err := g.getValue(name, timesType)
	if err != nil {
		return nil, err
	}
	return v.([]time.Time), nil
}

// TimesD is the same as TimesE, but returns the default value if there is
// an error.
func (g *OptGroup) TimesD(name string, _default []time.Time) []time.Time {
	if value, err := g.TimesE(name); err == nil {
		return value
	}
	return _default
}

// Times is the same as TimesE, but panic if there is an error.
func (g *OptGroup) Times(name string) []time.Time {
	value, err := g.TimesE(name)
	if err != nil {
		panic(err)
	}
	return value
}
