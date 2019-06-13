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

import "time"

// GetVersion is equal to Conf.GetVersion().
func GetVersion() Opt {
	return Conf.GetVersion()
}

// SetVersion is equal to Conf.SetVersion(version).
func SetVersion(version Opt) {
	Conf.SetVersion(version)
}

// SetStringVersion is equal to Conf.SetStringVersion(version).
func SetStringVersion(version string) {
	Conf.SetStringVersion(version)
}

// AddWatcher is equal to Conf.AddWatcher(watchers...).
func AddWatcher(watchers ...Watcher) {
	Conf.AddWatcher(watchers...)
}

// LoadSource is equal to Conf.LoadSource(source, force...).
func LoadSource(source Source, force ...bool) error {
	return Conf.LoadSource(source, force...)
}

// AddDecoder is equal to Conf.AddDecoder(decoder, force...).
func AddDecoder(decoder Decoder, force ...bool) (ok bool) {
	return Conf.AddDecoder(decoder, force...)
}

// GetDecoder is equal to Conf.GetDecoder(_type).
func GetDecoder(_type string) (decoder Decoder, ok bool) {
	return Conf.GetDecoder(_type)
}

// AddDecoderAlias is equal to Conf.AddDecoderAlias(_type, alias).
func AddDecoderAlias(_type, alias string) {
	Conf.AddDecoderAlias(_type, alias)
}

// AllGroups is equal to Conf.AllGroups().
func AllGroups() []*OptGroup {
	return Conf.AllGroups()
}

// Close is equal to Conf.Close().
func Close() error {
	return Conf.Close()
}

// Observe is equal to Conf.Observe(observer).
func Observe(observer func(group string, opt string, oldValue, newValue interface{})) {
	Conf.Observe(observer)
}

// SetErrHandler is equal to Conf.SetErrHandler(handler).
func SetErrHandler(handler func(error)) {
	Conf.SetErrHandler(handler)
}

// UpdateOptValue is equal to Conf.UpdateOptValue(groupName, optName, optValue).
func UpdateOptValue(groupName, optName string, optValue interface{}) {
	Conf.UpdateOptValue(groupName, optName, optValue)
}

// Group is equal to Conf.Group(group).
func Group(group string) *OptGroup {
	return Conf.Group(group)
}

// G is short for Group.
func G(group string) *OptGroup {
	return Conf.Group(group)
}

// NewGroup is eqaul to Conf.NewGroup(group).
func NewGroup(group string) *OptGroup {
	return Conf.NewGroup(group)
}

// FreezeOpt is eqaul to Conf.FreezeOpt(names...).
func FreezeOpt(names ...string) {
	Conf.FreezeOpt(names...)
}

// UnfreezeOpt is eqaul to Conf.UnfreezeOpt(names...).
func UnfreezeOpt(names ...string) {
	Conf.UnfreezeOpt(names...)
}

// OptIsFrozen is equal to Conf.OptIsFrozen(name).
func OptIsFrozen(name string) (frozen bool) {
	return Conf.OptIsFrozen(name)
}

// RegisterOpt is equal to Conf.RegisterOpt(opt, force...).
func RegisterOpt(opt Opt, force ...bool) (ok bool) {
	return Conf.RegisterOpt(opt, force...)
}

// RegisterOpts is equal to Conf.RegisterOpts(opts, force...).
func RegisterOpts(opts []Opt, force ...bool) (ok bool) {
	return Conf.RegisterOpts(opts, force...)
}

// RegisterStruct is equal to Conf.RegisterStruct(v).
func RegisterStruct(v interface{}) {
	Conf.RegisterStruct(v)
}

// Get is equal to Conf.Get(name).
func Get(name string) (value interface{}) {
	return Conf.Get(name)
}

// GetBool is equal to Conf.GetBool(name).
func GetBool(name string) bool {
	return Conf.GetBool(name)
}

// GetDuration is equal to Conf.GetDuration(name).
func GetDuration(name string) time.Duration {
	return Conf.GetDuration(name)
}

// GetDurationSlice is equal to Conf.GetDurationSlice(name).
func GetDurationSlice(name string) []time.Duration {
	return Conf.GetDurationSlice(name)
}

// GetFloat64 is equal to Conf.GetFloat64(name).
func GetFloat64(name string) float64 {
	return Conf.GetFloat64(name)
}

// GetFloat64Slice is equal to Conf.GetFloat64Slice(name).
func GetFloat64Slice(name string) []float64 {
	return Conf.GetFloat64Slice(name)
}

// GetInt is equal to Conf.GetInt(name).
func GetInt(name string) int {
	return Conf.GetInt(name)
}

// GetInt32 is equal to Conf.GetInt32(name).
func GetInt32(name string) int32 {
	return Conf.GetInt32(name)
}

// GetInt64 is equal to Conf.GetInt64(name).
func GetInt64(name string) int64 {
	return Conf.GetInt64(name)
}

// GetIntSlice is equal to Conf.GetIntSlice(name).
func GetIntSlice(name string) []int {
	return Conf.GetIntSlice(name)
}

// GetString is equal to Conf.GetString(name).
func GetString(name string) string {
	return Conf.GetString(name)
}

// GetStringSlice is equal to Conf.GetStringSlice(name).
func GetStringSlice(name string) []string {
	return Conf.GetStringSlice(name)
}

// GetTime is equal to Conf.GetTime(name).
func GetTime(name string) time.Time {
	return Conf.GetTime(name)
}

// GetUint is equal to Conf.GetUint(name).
func GetUint(name string) uint {
	return Conf.GetUint(name)
}

// GetUint32 is equal to Conf.GetUint32(name).
func GetUint32(name string) uint32 {
	return Conf.GetUint32(name)
}

// GetUint64 is equal to Conf.GetUint64(name).
func GetUint64(name string) uint64 {
	return Conf.GetUint64(name)
}

// GetUintSlice is equal to Conf.GetUintSlice(name).
func GetUintSlice(name string) []uint {
	return Conf.GetUintSlice(name)
}
