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
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

// ErrNoOpt represents an error that the option does not exist.
var ErrNoOpt = errors.New("no option")

// VersionOpt reprensents a version option.
var VersionOpt = StrOpt("version", "Print the version and exit.").S("v").D("1.0.0")

// Conf is the default global Config.
var Conf = New()

type atomicValue atomic.Value

func (v *atomicValue) Load() interface{} { return (*atomic.Value)(v).Load() }

type option struct {
	value atomicValue
	opt   Opt
}

func (o *option) GetValue() interface{} {
	return o.value.Load()
}

func (o *option) Get() (value interface{}) {
	if value = o.value.Load(); value == nil {
		value = o.opt.Default
	}
	return
}

func (o *option) Set(c *Config, newvalue interface{}) {
	oldvalue := o.value.Swap(newvalue)
	if oldvalue == nil {
		oldvalue = o.opt.Default
	}
	c.observe(o, oldvalue, newvalue)
}

// Observer is used to observe the change of the option value.
type Observer func(optName string, oldValue, newValue interface{})

// Config is used to manage the configuration options.
type Config struct {
	// Args is the CLI rest arguments.
	//
	// Default: nil
	Args []string

	// Version is the version of the application, which is used by CLI.
	Version Opt

	// Errorf is used to log the error.
	//
	// Default: log.Printf
	Errorf func(format string, args ...interface{})

	gen       uint64
	gsep      string
	ignore    bool
	options   map[string]*option
	aliases   map[string]string
	daliases  map[string]string
	decoders  map[string]Decoder
	observers []Observer
	exit      chan struct{}
}

// New returns a new Config with the "json", "yaml/yml" and "ini" decoder.
func New() *Config {
	c := &Config{
		gsep:     ".",
		ignore:   true,
		options:  make(map[string]*option, 32),
		aliases:  make(map[string]string, 8),
		daliases: make(map[string]string, 4),
		decoders: make(map[string]Decoder, 4),
		exit:     make(chan struct{}),
	}

	c.Version = VersionOpt
	c.AddDecoder("ini", NewIniDecoder())
	c.AddDecoder("yaml", NewYamlDecoder())
	c.AddDecoder("json", NewJSONDecoder())
	c.AddDecoderTypeAliases("yaml", "yml")
	return c
}

// reset clears the whole config for test.
func (c *Config) reset() {
	c.options = make(map[string]*option)
	c.aliases = make(map[string]string)
}

// Stop stops the watchers of all the sources.
func (c *Config) Stop() {
	select {
	case <-c.exit:
	default:
		close(c.exit)
	}
}

// SetVersion sets the version information.
func (c *Config) SetVersion(version string) {
	c.Version.Default = version
}

func (c *Config) fixOptionName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

func (c *Config) errorf(format string, args ...interface{}) {
	if c.Errorf == nil {
		log.Printf(format, args...)
	} else {
		c.Errorf(format, args...)
	}
}

// IgnoreNoOptError sets whether to ignore the error when updating the value
// of the option.
func (c *Config) IgnoreNoOptError(ignore bool) { c.ignore = ignore }

// Observe appends the observers to watch the change of all the option values.
func (c *Config) Observe(observers ...Observer) {
	c.observers = append(c.observers, observers...)
}

func (c *Config) observe(o *option, old, new interface{}) {
	if !reflect.DeepEqual(old, new) {
		atomic.AddUint64(&c.gen, 1)
		for _, observe := range c.observers {
			observe(o.opt.Name, old, new)
		}
		if o.opt.OnUpdate != nil {
			o.opt.OnUpdate(old, new)
		}
	}
}

func (c *Config) setOptAlias(old, new string) {
	old = c.fixOptionName(old)
	new = c.fixOptionName(new)
	if old == "" || new == "" {
		return
	}

	if opt, ok := c.options[new]; ok && !inString(old, opt.opt.Aliases) {
		opt.opt.Aliases = append(opt.opt.Aliases, old)
	}

	c.aliases[old] = new
}

func (c *Config) unsetOptAlias(name string) {
	delete(c.aliases, name)
	for oldname, newname := range c.aliases {
		if newname == name {
			delete(c.aliases, oldname)
		}
	}
}

func (c *Config) registerOpt(opt Opt) (o *option) {
	opt.check()
	if err := opt.validate(opt.Default); err != nil {
		panic(fmt.Errorf("invalid default '%v' for option named '%s': %s",
			opt.Default, opt.Name, err))
	}

	name := c.fixOptionName(opt.Name)
	if _, ok := c.options[name]; ok {
		panic(fmt.Errorf("the option named '%s' has been registered", name))
	}

	for _, alias := range opt.Aliases {
		c.setOptAlias(alias, opt.Name)
	}

	o = &option{opt: opt}
	c.options[name] = o
	return
}

// RegisterOpts registers a set of options.
//
// Notice: if a certain option has existed, it will panic.
func (c *Config) RegisterOpts(opts ...Opt) {
	names := make([]string, len(opts))
	for i, opt := range opts {
		opt.check()
		if err := opt.validate(opt.Default); err != nil {
			panic(fmt.Errorf("invalid default '%v' for option named '%s': %s",
				opt.Default, opt.Name, err))
		}
		names[i] = c.fixOptionName(opt.Name)
	}

	for _, name := range names {
		if _, ok := c.options[name]; ok {
			panic(fmt.Errorf("the option named '%s' has been registered", name))
		}
	}

	for i, opt := range opts {
		for _, alias := range opt.Aliases {
			c.setOptAlias(alias, opt.Name)
		}
		c.options[names[i]] = &option{opt: opt}
	}
}

// UnregisterOpts unregisters the registered options.
func (c *Config) UnregisterOpts(optNames ...string) {
	for _, name := range optNames {
		c.unregisterOpt(name)
	}
}

func (c *Config) unregisterOpt(name string) {
	name = c.fixOptionName(name)
	delete(c.options, name)
	c.unsetOptAlias(name)
}

// OptIsSet reports whether the option named name is set.
//
// Return false if the option does not exist.
func (c *Config) OptIsSet(name string) (yes bool) {
	name = c.fixOptionName(name)
	if opt, ok := c.options[name]; ok {
		yes = opt.value.Load() != nil
	} else if name, ok = c.aliases[name]; ok {
		if opt, ok = c.options[name]; ok {
			yes = opt.value.Load() != nil
		}
	}
	return
}

// HasOpt reports whether the option named name has been registered.
func (c *Config) HasOpt(name string) (yes bool) {
	name = c.fixOptionName(name)
	if _, yes = c.options[name]; !yes {
		if alias, ok := c.aliases[name]; ok {
			_, yes = c.options[alias]
		}
	}
	return
}

// GetOpt returns the registered option by the name.
func (c *Config) GetOpt(name string) (opt Opt, ok bool) {
	name = c.fixOptionName(name)
	option, ok := c.options[name]
	if ok {
		opt = option.opt
	} else if name, ok = c.aliases[name]; ok {
		if option, ok = c.options[name]; ok {
			opt = option.opt
		}
	}
	return
}

// GetAllOpts returns all the registered options.
func (c *Config) GetAllOpts() []Opt {
	opts := make([]Opt, len(c.options))
	var index int
	for _, opt := range c.options {
		opts[index] = opt.opt
		index++
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Name < opts[j].Name })
	return opts
}

func (c *Config) updateOpt(name string, value interface{}, set bool) (
	*option, interface{}, error) {
	if value == nil {
		return nil, nil, nil
	}

	// Get the option by the name.
	opt, ok := c.options[name]
	if !ok {
		if alias, ok := c.aliases[name]; !ok {
			return nil, nil, ErrNoOpt
		} else if opt, ok = c.options[alias]; !ok {
			return nil, nil, ErrNoOpt
		}
	}

	// Parse the option value
	newvalue, err := opt.opt.Parser(value)
	if err != nil {
		return nil, nil, err
	} else if newvalue == nil {
		panic(fmt.Errorf("the parser of option named '%s' returns nil", name))
	}

	// Validate the option value
	if err = opt.opt.validate(newvalue); err != nil {
		return nil, nil, err
	}

	if set {
		opt.Set(c, newvalue)
	}

	return opt, newvalue, nil
}

func (c *Config) checkMultilayerMap(ms map[string]interface{}) (yes bool) {
	for _, value := range ms {
		switch value.(type) {
		case map[string]string,
			map[string]interface{},
			map[interface{}]interface{}:
			return true
		}
	}
	return false
}

func (c *Config) flatmap2(prefix string, results map[string]interface{},
	maps map[interface{}]interface{}) {
	if prefix != "" {
		prefix += c.gsep
	}

	for key, value := range maps {
		var k string
		if _key, ok := key.(string); ok {
			k = _key
		} else {
			k = fmt.Sprint(k)
		}

		switch vs := value.(type) {
		case map[string]string:
			for _k, _v := range vs {
				results[prefix+k+c.gsep+_k] = _v
			}
		case map[string]interface{}:
			c.flatmap(prefix+k, results, vs)
		case map[interface{}]interface{}:
			c.flatmap2(prefix+k, results, vs)
		default:
			results[prefix+k] = value
		}
	}
}

func (c *Config) flatmap(prefix string, results, maps map[string]interface{}) {
	if prefix != "" {
		prefix += c.gsep
	}

	for key, value := range maps {
		switch vs := value.(type) {
		case map[string]string:
			for k, v := range vs {
				results[prefix+key+c.gsep+k] = v
			}
		case map[string]interface{}:
			c.flatmap(prefix+key, results, vs)
		case map[interface{}]interface{}:
			c.flatmap2(prefix+key, results, vs)
		default:
			results[prefix+key] = value
		}
	}

	return
}

func (c *Config) flatMap(maps map[string]interface{}) map[string]interface{} {
	if c.checkMultilayerMap(maps) {
		tmp := make(map[string]interface{}, len(maps)*2)
		c.flatmap("", tmp, maps)
		return tmp
	}
	return maps
}

// LoadMap updates a set of the options together, but terminates to parse
// and load all if failing to parse the value of any option.
//
// If force is missing or false, ignore the assigned options.
func (c *Config) LoadMap(options map[string]interface{}, force ...bool) error {
	if len(options) == 0 {
		return nil
	}

	var _force bool
	if len(force) > 0 {
		_force = force[0]
	}

	type opt struct {
		name   string
		value  interface{}
		option *option
	}

	options = c.flatMap(options)
	opts := make([]opt, 0, len(options))

	for name, value := range options {
		name = c.fixOptionName(name)
		o, newv, err := c.updateOpt(name, value, false)
		switch err {
		case nil:
			if o.value.Load() != nil && !_force {
				continue
			}
		case ErrNoOpt:
			if c.ignore {
				continue
			}
			return fmt.Errorf("no option named '%s'", name)
		default:
			return err
		}
		opts = append(opts, opt{name: name, value: newv, option: o})
	}

	for _, opt := range opts {
		opt.option.Set(c, opt.value)
	}

	return nil
}

// Parse parses the option value named name, and returns it.
func (c *Config) Parse(name string, value interface{}) (interface{}, error) {
	name = c.fixOptionName(name)
	_, value, err := c.updateOpt(name, value, false)
	return value, err
}

// Set is used to reset the option named name to value.
func (c *Config) Set(name string, value interface{}) (err error) {
	name = c.fixOptionName(name)
	_, _, err = c.updateOpt(name, value, true)
	if err == ErrNoOpt && c.ignore {
		err = nil
	}
	return
}

// Get returns the value of the option named name.
//
// Return nil if this option does not exist.
func (c *Config) Get(name string) (value interface{}) {
	name = c.fixOptionName(name)
	if opt, ok := c.options[name]; ok {
		value = opt.Get()
	}
	return
}

// Must is the same as Get, but panic if the returned value is nil.
func (c *Config) Must(name string) (value interface{}) {
	if value = c.Get(name); value == nil {
		panic(fmt.Errorf("no option named name '%s'", name))
	}
	return
}

// GetBool returns the value of the option named name as bool.
func (c *Config) GetBool(name string) bool { return c.Must(name).(bool) }

// GetInt returns the value of the option named name as int.
func (c *Config) GetInt(name string) int { return c.Must(name).(int) }

// GetInt32 returns the value of the option named name as int32.
func (c *Config) GetInt32(name string) int32 { return c.Must(name).(int32) }

// GetInt64 returns the value of the option named name as int64.
func (c *Config) GetInt64(name string) int64 { return c.Must(name).(int64) }

// GetUint returns the value of the option named name as uint.
func (c *Config) GetUint(name string) uint { return c.Must(name).(uint) }

// GetUint32 returns the value of the option named name as uint32.
func (c *Config) GetUint32(name string) uint32 { return c.Must(name).(uint32) }

// GetUint64 returns the value of the option named name as uint64.
func (c *Config) GetUint64(name string) uint64 { return c.Must(name).(uint64) }

// GetFloat64 returns the value of the option named name as float64.
func (c *Config) GetFloat64(name string) float64 { return c.Must(name).(float64) }

// GetString returns the value of the option named name as string.
func (c *Config) GetString(name string) string { return c.Must(name).(string) }

// GetDuration returns the value of the option named name as time.Duration.
func (c *Config) GetDuration(name string) time.Duration { return c.Must(name).(time.Duration) }

// GetTime returns the value of the option named name as time.Time.
func (c *Config) GetTime(name string) time.Time { return c.Must(name).(time.Time) }

// GetIntSlice returns the value of the option named name as []int.
func (c *Config) GetIntSlice(name string) []int { return c.Must(name).([]int) }

// GetUintSlice returns the value of the option named name as []uint.
func (c *Config) GetUintSlice(name string) []uint { return c.Must(name).([]uint) }

// GetFloat64Slice returns the value of the option named name as []float64.
func (c *Config) GetFloat64Slice(name string) []float64 { return c.Must(name).([]float64) }

// GetStringSlice returns the value of the option named name as []string.
func (c *Config) GetStringSlice(name string) []string { return c.Must(name).([]string) }

// GetDurationSlice returns the value of the option named name as []time.Duration.
func (c *Config) GetDurationSlice(name string) []time.Duration { return c.Must(name).([]time.Duration) }
