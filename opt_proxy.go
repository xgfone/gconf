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

import "time"

// OptProxy is a proxy for the option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxy struct {
	config *Config
	option *option
}

// NewOptProxy registers the option into c and returns a proxy of opt.
func NewOptProxy(c *Config, opt Opt) OptProxy {
	return OptProxy{option: c.registerOpt(opt), config: c}
}

// NewOptProxy registers the option and returns a new proxy of the option.
func (c *Config) NewOptProxy(opt Opt) OptProxy {
	return OptProxy{option: c.registerOpt(opt), config: c}
}

// Name returns the name of the option.
func (o *OptProxy) Name() string { return o.option.opt.Name }

// Opt returns the registered and proxied option.
func (o *OptProxy) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxy) Get() interface{} {
	return o.config.Get(o.option.opt.Name)
}

// Set sets the value of the option to value.
func (o *OptProxy) Set(value interface{}) (err error) {
	return o.config.Set(o.option.opt.Name, value)
}

// OnUpdate resets the update callback function of the option and returns itself.
func (o *OptProxy) OnUpdate(callback func(old, new interface{})) *OptProxy {
	o.option.opt.OnUpdate = callback
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxy) IsCli(cli bool) *OptProxy {
	o.option.opt.IsCli = cli
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxy) Aliases(aliases ...string) *OptProxy {
	o.option.opt = o.option.opt.As(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxy) Short(short string) *OptProxy {
	o.option.opt = o.option.opt.S(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxy) Validators(validators ...Validator) *OptProxy {
	o.option.opt = o.option.opt.V(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxy) Default(_default interface{}) *OptProxy {
	o.option.opt = o.option.opt.D(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxy) Parser(parser Parser) *OptProxy {
	o.option.opt = o.option.opt.P(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// NewOptProxy registers the option and returns a new proxy of the option.
func (g *OptGroup) NewOptProxy(opt Opt) OptProxy {
	opt.Name = g.prefix + opt.Name
	return OptProxy{option: g.config.registerOpt(opt), config: g.config}
}

// NewBool creates and registers a bool option, then returns its proxy.
func (g *OptGroup) NewBool(name string, _default bool, help string) *OptProxyBool {
	return g.config.NewBool(g.prefix+name, _default, help)
}

// NewInt creates and registers an int option, then returns its proxy.
func (g *OptGroup) NewInt(name string, _default int, help string) *OptProxyInt {
	return g.config.NewInt(g.prefix+name, _default, help)
}

// NewInt32 creates and registers a int32 option, then returns its proxy.
func (g *OptGroup) NewInt32(name string, _default int32, help string) *OptProxyInt32 {
	return g.config.NewInt32(g.prefix+name, _default, help)
}

// NewInt64 creates and registers a int64 option, then returns its proxy.
func (g *OptGroup) NewInt64(name string, _default int64, help string) *OptProxyInt64 {
	return g.config.NewInt64(g.prefix+name, _default, help)
}

// NewUint is equal to Conf.NewUint(name, _default, help).
func (g *OptGroup) NewUint(name string, _default uint, help string) *OptProxyUint {
	return g.config.NewUint(g.prefix+name, _default, help)
}

// NewUint32 creates and registers a uint32 option, then returns its proxy.
func (g *OptGroup) NewUint32(name string, _default uint32, help string) *OptProxyUint32 {
	return g.config.NewUint32(g.prefix+name, _default, help)
}

// NewUint64 creates and registers a uint64 option, then returns its proxy.
func (g *OptGroup) NewUint64(name string, _default uint64, help string) *OptProxyUint64 {
	return g.config.NewUint64(g.prefix+name, _default, help)
}

// NewFloat64 creates and registers a float64 option, then returns its proxy.
func (g *OptGroup) NewFloat64(name string, _default float64, help string) *OptProxyFloat64 {
	return g.config.NewFloat64(g.prefix+name, _default, help)
}

// NewString creates and registers a string option, then returns its proxy.
func (g *OptGroup) NewString(name, _default, help string) *OptProxyString {
	return g.config.NewString(g.prefix+name, _default, help)
}

// NewDuration creates and registers a time.Duration option, then returns its proxy.
func (g *OptGroup) NewDuration(name string, _default time.Duration, help string) *OptProxyDuration {
	return g.config.NewDuration(g.prefix+name, _default, help)
}

// NewTime creates and registers a time.Time option, then returns its proxy.
func (g *OptGroup) NewTime(name string, _default time.Time, help string) *OptProxyTime {
	return g.config.NewTime(g.prefix+name, _default, help)
}

// NewStringSlice creates and registers a []string option, then returns its proxy.
func (g *OptGroup) NewStringSlice(name string, _default []string, help string) *OptProxyStringSlice {
	return g.config.NewStringSlice(g.prefix+name, _default, help)
}

// NewIntSlice creates and registers a []int option, then returns its proxy.
func (g *OptGroup) NewIntSlice(name string, _default []int, help string) *OptProxyIntSlice {
	return g.config.NewIntSlice(g.prefix+name, _default, help)
}

// NewUintSlice creates and registers a []uint option, then returns its proxy.
func (g *OptGroup) NewUintSlice(name string, _default []uint, help string) *OptProxyUintSlice {
	return g.config.NewUintSlice(g.prefix+name, _default, help)
}

// NewFloat64Slice creates and registers a []float64 option, then returns its proxy.
func (g *OptGroup) NewFloat64Slice(name string, _default []float64, help string) *OptProxyFloat64Slice {
	return g.config.NewFloat64Slice(g.prefix+name, _default, help)
}

// NewDurationSlice creates and registers a []time.Duration option, then returns its proxy.
func (g *OptGroup) NewDurationSlice(name string, _default []time.Duration, help string) *OptProxyDurationSlice {
	return g.config.NewDurationSlice(g.prefix+name, _default, help)
}

////////////////////////////////////////////////////////////////////////////

// OptProxyBool is a proxy for the bool option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyBool struct{ OptProxy }

// NewBool is equal to Conf.NewBool(name, _default, help).
func NewBool(name string, _default bool, help string) *OptProxyBool {
	return Conf.NewBool(name, _default, help)
}

// NewBool creates and registers a bool option, then returns its proxy.
func (c *Config) NewBool(name string, _default bool, help string) *OptProxyBool {
	return &OptProxyBool{c.NewOptProxy(BoolOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyBool) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyBool) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyBool) Get() bool { return o.OptProxy.Get().(bool) }

// Set sets the value of the option to value.
func (o *OptProxyBool) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyBool) OnUpdate(f func(old, new interface{})) *OptProxyBool {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyBool) IsCli(cli bool) *OptProxyBool {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyBool) Aliases(aliases ...string) *OptProxyBool {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyBool) Short(short string) *OptProxyBool {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyBool) Validators(validators ...Validator) *OptProxyBool {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyBool) Default(_default interface{}) *OptProxyBool {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyBool) Parser(parser Parser) *OptProxyBool {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyInt is a proxy for the int option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyInt struct{ OptProxy }

// NewInt is equal to Conf.NewInt(name, _default, help).
func NewInt(name string, _default int, help string) *OptProxyInt {
	return Conf.NewInt(name, _default, help)
}

// NewInt creates and registers an int option, then returns its proxy.
func (c *Config) NewInt(name string, _default int, help string) *OptProxyInt {
	return &OptProxyInt{c.NewOptProxy(IntOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyInt) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyInt) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyInt) Get() int { return o.OptProxy.Get().(int) }

// Set sets the value of the option to value.
func (o *OptProxyInt) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyInt) OnUpdate(f func(old, new interface{})) *OptProxyInt {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyInt) IsCli(cli bool) *OptProxyInt {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyInt) Aliases(aliases ...string) *OptProxyInt {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyInt) Short(short string) *OptProxyInt {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyInt) Validators(validators ...Validator) *OptProxyInt {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyInt) Default(_default interface{}) *OptProxyInt {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyInt) Parser(parser Parser) *OptProxyInt {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyInt32 is a proxy for the int32 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyInt32 struct{ OptProxy }

// NewInt32 is equal to Conf.NewInt32(name, _default, help).
func NewInt32(name string, _default int32, help string) *OptProxyInt32 {
	return Conf.NewInt32(name, _default, help)
}

// NewInt32 creates and registers a int32 option, then returns its proxy.
func (c *Config) NewInt32(name string, _default int32, help string) *OptProxyInt32 {
	return &OptProxyInt32{c.NewOptProxy(Int32Opt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyInt32) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyInt32) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyInt32) Get() int32 { return o.OptProxy.Get().(int32) }

// Set sets the value of the option to value.
func (o *OptProxyInt32) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyInt32) OnUpdate(f func(old, new interface{})) *OptProxyInt32 {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyInt32) IsCli(cli bool) *OptProxyInt32 {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyInt32) Aliases(aliases ...string) *OptProxyInt32 {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyInt32) Short(short string) *OptProxyInt32 {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyInt32) Validators(validators ...Validator) *OptProxyInt32 {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyInt32) Default(_default interface{}) *OptProxyInt32 {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyInt32) Parser(parser Parser) *OptProxyInt32 {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyInt64 is a proxy for the int64 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyInt64 struct{ OptProxy }

// NewInt64 is equal to Conf.NewInt64(name, _default, help).
func NewInt64(name string, _default int64, help string) *OptProxyInt64 {
	return Conf.NewInt64(name, _default, help)
}

// NewInt64 creates and registers a int64 option, then returns its proxy.
func (c *Config) NewInt64(name string, _default int64, help string) *OptProxyInt64 {
	return &OptProxyInt64{c.NewOptProxy(Int64Opt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyInt64) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyInt64) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyInt64) Get() int64 { return o.OptProxy.Get().(int64) }

// Set sets the value of the option to value.
func (o *OptProxyInt64) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyInt64) OnUpdate(f func(old, new interface{})) *OptProxyInt64 {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyInt64) IsCli(cli bool) *OptProxyInt64 {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyInt64) Aliases(aliases ...string) *OptProxyInt64 {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyInt64) Short(short string) *OptProxyInt64 {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyInt64) Validators(validators ...Validator) *OptProxyInt64 {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyInt64) Default(_default interface{}) *OptProxyInt64 {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyInt64) Parser(parser Parser) *OptProxyInt64 {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyUint is a proxy for the uint option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyUint struct{ OptProxy }

// NewUint is equal to Conf.NewUint(name, _default, help).
func NewUint(name string, _default uint, help string) *OptProxyUint {
	return Conf.NewUint(name, _default, help)
}

// NewUint creates and registers a uint option, then returns its proxy.
func (c *Config) NewUint(name string, _default uint, help string) *OptProxyUint {
	return &OptProxyUint{c.NewOptProxy(UintOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyUint) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyUint) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyUint) Get() uint { return o.OptProxy.Get().(uint) }

// Set sets the value of the option to value.
func (o *OptProxyUint) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyUint) OnUpdate(f func(old, new interface{})) *OptProxyUint {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyUint) IsCli(cli bool) *OptProxyUint {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyUint) Aliases(aliases ...string) *OptProxyUint {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyUint) Short(short string) *OptProxyUint {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyUint) Validators(validators ...Validator) *OptProxyUint {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyUint) Default(_default interface{}) *OptProxyUint {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyUint) Parser(parser Parser) *OptProxyUint {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyUint32 is a proxy for the uint32 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyUint32 struct{ OptProxy }

// NewUint32 is equal to Conf.NewUint32(name, _default, help).
func NewUint32(name string, _default uint32, help string) *OptProxyUint32 {
	return Conf.NewUint32(name, _default, help)
}

// NewUint32 creates and registers a uint32 option, then returns its proxy.
func (c *Config) NewUint32(name string, _default uint32, help string) *OptProxyUint32 {
	return &OptProxyUint32{c.NewOptProxy(Uint32Opt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyUint32) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyUint32) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyUint32) Get() uint32 { return o.OptProxy.Get().(uint32) }

// Set sets the value of the option to value.
func (o *OptProxyUint32) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyUint32) OnUpdate(f func(old, new interface{})) *OptProxyUint32 {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyUint32) IsCli(cli bool) *OptProxyUint32 {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyUint32) Aliases(aliases ...string) *OptProxyUint32 {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyUint32) Short(short string) *OptProxyUint32 {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyUint32) Validators(validators ...Validator) *OptProxyUint32 {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyUint32) Default(_default interface{}) *OptProxyUint32 {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyUint32) Parser(parser Parser) *OptProxyUint32 {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyUint64 is a proxy for the uint64 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyUint64 struct{ OptProxy }

// NewUint64 is equal to Conf.NewUint64(name, _default, help).
func NewUint64(name string, _default uint64, help string) *OptProxyUint64 {
	return Conf.NewUint64(name, _default, help)
}

// NewUint64 creates and registers a uint64 option, then returns its proxy.
func (c *Config) NewUint64(name string, _default uint64, help string) *OptProxyUint64 {
	return &OptProxyUint64{c.NewOptProxy(Uint64Opt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyUint64) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyUint64) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyUint64) Get() uint64 { return o.OptProxy.Get().(uint64) }

// Set sets the value of the option to value.
func (o *OptProxyUint64) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyUint64) OnUpdate(f func(old, new interface{})) *OptProxyUint64 {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyUint64) IsCli(cli bool) *OptProxyUint64 {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyUint64) Aliases(aliases ...string) *OptProxyUint64 {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyUint64) Short(short string) *OptProxyUint64 {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyUint64) Validators(validators ...Validator) *OptProxyUint64 {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyUint64) Default(_default interface{}) *OptProxyUint64 {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyUint64) Parser(parser Parser) *OptProxyUint64 {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyFloat64 is a proxy for the float64 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyFloat64 struct{ OptProxy }

// NewFloat64 is equal to Conf.NewFloat64(name, _default, help).
func NewFloat64(name string, _default float64, help string) *OptProxyFloat64 {
	return Conf.NewFloat64(name, _default, help)
}

// NewFloat64 creates and registers a float64 option, then returns its proxy.
func (c *Config) NewFloat64(name string, _default float64, help string) *OptProxyFloat64 {
	return &OptProxyFloat64{c.NewOptProxy(Float64Opt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyFloat64) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyFloat64) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyFloat64) Get() float64 { return o.OptProxy.Get().(float64) }

// Set sets the value of the option to value.
func (o *OptProxyFloat64) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyFloat64) OnUpdate(f func(old, new interface{})) *OptProxyFloat64 {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyFloat64) IsCli(cli bool) *OptProxyFloat64 {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyFloat64) Aliases(aliases ...string) *OptProxyFloat64 {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyFloat64) Short(short string) *OptProxyFloat64 {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyFloat64) Validators(validators ...Validator) *OptProxyFloat64 {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyFloat64) Default(_default interface{}) *OptProxyFloat64 {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyFloat64) Parser(parser Parser) *OptProxyFloat64 {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyString is a proxy for the string option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyString struct{ OptProxy }

// NewString is equal to Conf.NewString(name, _default, help).
func NewString(name, _default, help string) *OptProxyString {
	return Conf.NewString(name, _default, help)
}

// NewString creates and registers a string option, then returns its proxy.
func (c *Config) NewString(name, _default, help string) *OptProxyString {
	return &OptProxyString{c.NewOptProxy(StrOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyString) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyString) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyString) Get() string { return o.OptProxy.Get().(string) }

// Set sets the value of the option to value.
func (o *OptProxyString) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyString) OnUpdate(f func(old, new interface{})) *OptProxyString {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyString) IsCli(cli bool) *OptProxyString {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyString) Aliases(aliases ...string) *OptProxyString {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyString) Short(short string) *OptProxyString {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyString) Validators(validators ...Validator) *OptProxyString {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyString) Default(_default interface{}) *OptProxyString {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyString) Parser(parser Parser) *OptProxyString {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyDuration is a proxy for the time.Duration option registered
// into Config, which can be used to modify the attributions of the option
// and update its value directly.
type OptProxyDuration struct{ OptProxy }

// NewDuration is equal to Conf.NewDuration(name, _default, help).
func NewDuration(name string, _default time.Duration, help string) *OptProxyDuration {
	return Conf.NewDuration(name, _default, help)
}

// NewDuration creates and registers a time.Duration option, then returns its proxy.
func (c *Config) NewDuration(name string, _default time.Duration, help string) *OptProxyDuration {
	return &OptProxyDuration{c.NewOptProxy(DurationOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyDuration) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyDuration) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyDuration) Get() time.Duration { return o.OptProxy.Get().(time.Duration) }

// Set sets the value of the option to value.
func (o *OptProxyDuration) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyDuration) OnUpdate(f func(old, new interface{})) *OptProxyDuration {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyDuration) IsCli(cli bool) *OptProxyDuration {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyDuration) Aliases(aliases ...string) *OptProxyDuration {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyDuration) Short(short string) *OptProxyDuration {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyDuration) Validators(validators ...Validator) *OptProxyDuration {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyDuration) Default(_default interface{}) *OptProxyDuration {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyDuration) Parser(parser Parser) *OptProxyDuration {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyTime is a proxy for the time.Time option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyTime struct{ OptProxy }

// NewTime is equal to Conf.NewTime(name, _default, help).
func NewTime(name string, _default time.Time, help string) *OptProxyTime {
	return Conf.NewTime(name, _default, help)
}

// NewTime creates and registers a time.Time option, then returns its proxy.
func (c *Config) NewTime(name string, _default time.Time, help string) *OptProxyTime {
	return &OptProxyTime{c.NewOptProxy(TimeOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyTime) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyTime) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyTime) Get() time.Time { return o.OptProxy.Get().(time.Time) }

// Set sets the value of the option to value.
func (o *OptProxyTime) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyTime) OnUpdate(f func(old, new interface{})) *OptProxyTime {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyTime) IsCli(cli bool) *OptProxyTime {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyTime) Aliases(aliases ...string) *OptProxyTime {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyTime) Short(short string) *OptProxyTime {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyTime) Validators(validators ...Validator) *OptProxyTime {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyTime) Default(_default interface{}) *OptProxyTime {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyTime) Parser(parser Parser) *OptProxyTime {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyStringSlice is a proxy for the []string option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyStringSlice struct{ OptProxy }

// NewStringSlice is equal to Conf.NewStringSlice(name, _default, help).
func NewStringSlice(name string, _default []string, help string) *OptProxyStringSlice {
	return Conf.NewStringSlice(name, _default, help)
}

// NewStringSlice creates and registers a []string option, then returns its proxy.
func (c *Config) NewStringSlice(name string, _default []string, help string) *OptProxyStringSlice {
	return &OptProxyStringSlice{c.NewOptProxy(StrSliceOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyStringSlice) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyStringSlice) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyStringSlice) Get() []string { return o.OptProxy.Get().([]string) }

// Set sets the value of the option to value.
func (o *OptProxyStringSlice) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyStringSlice) OnUpdate(f func(old, new interface{})) *OptProxyStringSlice {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyStringSlice) IsCli(cli bool) *OptProxyStringSlice {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyStringSlice) Aliases(aliases ...string) *OptProxyStringSlice {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyStringSlice) Short(short string) *OptProxyStringSlice {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyStringSlice) Validators(validators ...Validator) *OptProxyStringSlice {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyStringSlice) Default(_default interface{}) *OptProxyStringSlice {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyStringSlice) Parser(parser Parser) *OptProxyStringSlice {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyIntSlice is a proxy for the []int option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyIntSlice struct{ OptProxy }

// NewIntSlice is equal to Conf.NewIntSlice(name, _default, help).
func NewIntSlice(name string, _default []int, help string) *OptProxyIntSlice {
	return Conf.NewIntSlice(name, _default, help)
}

// NewIntSlice creates and registers a []int option, then returns its proxy.
func (c *Config) NewIntSlice(name string, _default []int, help string) *OptProxyIntSlice {
	return &OptProxyIntSlice{c.NewOptProxy(IntSliceOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyIntSlice) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyIntSlice) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyIntSlice) Get() []int { return o.OptProxy.Get().([]int) }

// Set sets the value of the option to value.
func (o *OptProxyIntSlice) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyIntSlice) OnUpdate(f func(old, new interface{})) *OptProxyIntSlice {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyIntSlice) IsCli(cli bool) *OptProxyIntSlice {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyIntSlice) Aliases(aliases ...string) *OptProxyIntSlice {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyIntSlice) Short(short string) *OptProxyIntSlice {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyIntSlice) Validators(validators ...Validator) *OptProxyIntSlice {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyIntSlice) Default(_default interface{}) *OptProxyIntSlice {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyIntSlice) Parser(parser Parser) *OptProxyIntSlice {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyUintSlice is a proxy for the []uint option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyUintSlice struct{ OptProxy }

// NewUintSlice is equal to Conf.NewUintSlice(name, _default, help).
func NewUintSlice(name string, _default []uint, help string) *OptProxyUintSlice {
	return Conf.NewUintSlice(name, _default, help)
}

// NewUintSlice creates and registers a []uint option, then returns its proxy.
func (c *Config) NewUintSlice(name string, _default []uint, help string) *OptProxyUintSlice {
	return &OptProxyUintSlice{c.NewOptProxy(UintSliceOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyUintSlice) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyUintSlice) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyUintSlice) Get() []uint { return o.OptProxy.Get().([]uint) }

// Set sets the value of the option to value.
func (o *OptProxyUintSlice) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyUintSlice) OnUpdate(f func(old, new interface{})) *OptProxyUintSlice {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyUintSlice) IsCli(cli bool) *OptProxyUintSlice {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyUintSlice) Aliases(aliases ...string) *OptProxyUintSlice {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyUintSlice) Short(short string) *OptProxyUintSlice {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyUintSlice) Validators(validators ...Validator) *OptProxyUintSlice {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyUintSlice) Default(_default interface{}) *OptProxyUintSlice {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyUintSlice) Parser(parser Parser) *OptProxyUintSlice {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyFloat64Slice is a proxy for the []float64 option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyFloat64Slice struct{ OptProxy }

// NewFloat64Slice is equal to Conf.NewFloat64Slice(name, _default, help).
func NewFloat64Slice(name string, _default []float64, help string) *OptProxyFloat64Slice {
	return Conf.NewFloat64Slice(name, _default, help)
}

// NewFloat64Slice creates and registers a []float64 option, then returns its proxy.
func (c *Config) NewFloat64Slice(name string, _default []float64, help string) *OptProxyFloat64Slice {
	return &OptProxyFloat64Slice{c.NewOptProxy(Float64SliceOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyFloat64Slice) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyFloat64Slice) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyFloat64Slice) Get() []float64 { return o.OptProxy.Get().([]float64) }

// Set sets the value of the option to value.
func (o *OptProxyFloat64Slice) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyFloat64Slice) OnUpdate(f func(old, new interface{})) *OptProxyFloat64Slice {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyFloat64Slice) IsCli(cli bool) *OptProxyFloat64Slice {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyFloat64Slice) Aliases(aliases ...string) *OptProxyFloat64Slice {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyFloat64Slice) Short(short string) *OptProxyFloat64Slice {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyFloat64Slice) Validators(validators ...Validator) *OptProxyFloat64Slice {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyFloat64Slice) Default(_default interface{}) *OptProxyFloat64Slice {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyFloat64Slice) Parser(parser Parser) *OptProxyFloat64Slice {
	o.OptProxy.Parser(parser)
	return o
}

////////////////////////////////////////////////////////////////////////////

// OptProxyDurationSlice is a proxy for the []time.Duration option registered into Config,
// which can be used to modify the attributions of the option and
// update its value directly.
type OptProxyDurationSlice struct{ OptProxy }

// NewDurationSlice is equal to Conf.NewDurationSlice(name, _default, help).
func NewDurationSlice(name string, _default []time.Duration, help string) *OptProxyDurationSlice {
	return Conf.NewDurationSlice(name, _default, help)
}

// NewDurationSlice creates and registers a []time.Duration option, then returns its proxy.
func (c *Config) NewDurationSlice(name string, _default []time.Duration, help string) *OptProxyDurationSlice {
	return &OptProxyDurationSlice{c.NewOptProxy(DurationSliceOpt(name, help).D(_default))}
}

// Name returns the name of the option.
func (o *OptProxyDurationSlice) Name() string { return o.OptProxy.Name() }

// Opt returns the registered and proxied option.
func (o *OptProxyDurationSlice) Opt() Opt { return o.option.opt }

// Get returns the value of the option.
func (o *OptProxyDurationSlice) Get() []time.Duration { return o.OptProxy.Get().([]time.Duration) }

// Set sets the value of the option to value.
func (o *OptProxyDurationSlice) Set(value interface{}) (err error) {
	return o.OptProxy.Set(value)
}

// OnUpdate resets the update callback of the option and returns itself.
func (o *OptProxyDurationSlice) OnUpdate(f func(old, new interface{})) *OptProxyDurationSlice {
	o.OptProxy.OnUpdate(f)
	return o
}

// IsCli resets the cli flag of the option and returns itself.
func (o *OptProxyDurationSlice) IsCli(cli bool) *OptProxyDurationSlice {
	o.OptProxy.IsCli(cli)
	return o
}

// Aliases appends the aliases of the option and returns itself.
func (o *OptProxyDurationSlice) Aliases(aliases ...string) *OptProxyDurationSlice {
	o.OptProxy.Aliases(aliases...)
	return o
}

// Short resets the short name of the option and returns itself.
func (o *OptProxyDurationSlice) Short(short string) *OptProxyDurationSlice {
	o.OptProxy.Short(short)
	return o
}

// Validators appends the validators of the option and returns itself.
func (o *OptProxyDurationSlice) Validators(validators ...Validator) *OptProxyDurationSlice {
	o.OptProxy.Validators(validators...)
	return o
}

// Default resets the default value of the option and returns itself.
func (o *OptProxyDurationSlice) Default(_default interface{}) *OptProxyDurationSlice {
	o.OptProxy.Default(_default)
	return o
}

// Parser resets the parser of the option and returns itself.
func (o *OptProxyDurationSlice) Parser(parser Parser) *OptProxyDurationSlice {
	o.OptProxy.Parser(parser)
	return o
}
