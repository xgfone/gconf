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

// Args is equal to Conf.Args().
func Args() []string {
	return Conf.Args()
}

// SetArgs is equal to Conf.SetArgs(args).
func SetArgs(args []string) {
	Conf.SetArgs(args)
}

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

// LoadMap is equal to Conf.LoadMap(m, force...).
func LoadMap(m map[string]interface{}, force ...bool) error {
	return Conf.LoadMap(m, force...)
}

// LoadDataSet is equal to Conf.LoadDataSet(ds, force...).
func LoadDataSet(ds DataSet, force ...bool) error {
	return Conf.LoadDataSet(ds, force...)
}

// LoadDataSetCallback is equal to Conf.LoadDataSetCallback(ds, err).
func LoadDataSetCallback(ds DataSet, err error) bool {
	return Conf.LoadDataSetCallback(ds, err)
}

// LoadSource is equal to Conf.LoadSource(source, force...).
func LoadSource(source Source, force ...bool) {
	Conf.LoadSource(source, force...)
}

// LoadSourceWithoutWatch is equal to Conf.LoadSourceWithoutWatch(source, force...).
func LoadSourceWithoutWatch(source Source, force ...bool) {
	Conf.LoadSourceWithoutWatch(source, force...)
}

// LoadBackupFile is equal to Conf.LoadBackupFile(filename).
func LoadBackupFile(filename string) error {
	return Conf.LoadBackupFile(filename)
}

// Snapshot is equal to Conf.Snapshot().
func Snapshot() map[string]interface{} {
	return Conf.Snapshot()
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
func Close() {
	Conf.Close()
}

// CloseNotice is equal to Conf.CloseNotice().
func CloseNotice() <-chan struct{} {
	return Conf.CloseNotice()
}

// Traverse is equal to Conf.Traverse(f).
func Traverse(f func(group string, opt string, value interface{})) {
	Conf.Traverse(f)
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

// UpdateValue is equal to Conf.UpdateValue(key, value).
func UpdateValue(key string, value interface{}) {
	Conf.UpdateValue(key, value)
}

// Group is equal to Conf.Group(group).
func Group(group string) *OptGroup {
	return Conf.Group(group)
}

// G is short for Group.
func G(group string) *OptGroup {
	return Conf.Group(group)
}

// MustGroup is equal to Conf.MustGroup(group).
func MustGroup(group string) *OptGroup {
	return Conf.MustGroup(group)
}

// MustG is short for MustGroup(group).
func MustG(group string) *OptGroup {
	return Conf.MustG(group)
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

// RegisterOpts is equal to Conf.RegisterOpts(opts, force...).
func RegisterOpts(opts ...Opt) {
	Conf.RegisterOpts(opts...)
}

// UnregisterOpts is equal to Conf.UnregisterOpts(opts...).
func UnregisterOpts(opts ...Opt) {
	Conf.UnregisterOpts(opts...)
}

// RegisterStruct is equal to Conf.RegisterStruct(v).
func RegisterStruct(v interface{}) {
	Conf.RegisterStruct(v)
}

// Get is equal to Conf.Get(name).
func Get(name string) (value interface{}) {
	return Conf.Get(name)
}

// Must is equal to Conf.Must(name).
func Must(name string) (value interface{}) {
	return Conf.Must(name)
}

// GetBool is equal to Conf.GetBool(name).
func GetBool(name string) bool {
	return Conf.GetBool(name)
}

// MustBool is equal to Conf.MustBool(name).
func MustBool(name string) bool {
	return Conf.MustBool(name)
}

// GetDuration is equal to Conf.GetDuration(name).
func GetDuration(name string) time.Duration {
	return Conf.GetDuration(name)
}

// MustDuration is equal to Conf.MustDuration(name).
func MustDuration(name string) time.Duration {
	return Conf.MustDuration(name)
}

// GetDurationSlice is equal to Conf.GetDurationSlice(name).
func GetDurationSlice(name string) []time.Duration {
	return Conf.GetDurationSlice(name)
}

// MustDurationSlice is equal to Conf.MustDurationSlice(name).
func MustDurationSlice(name string) []time.Duration {
	return Conf.MustDurationSlice(name)
}

// GetFloat64 is equal to Conf.GetFloat64(name).
func GetFloat64(name string) float64 {
	return Conf.GetFloat64(name)
}

// MustFloat64 is equal to Conf.MustFloat64(name).
func MustFloat64(name string) float64 {
	return Conf.MustFloat64(name)
}

// GetFloat64Slice is equal to Conf.GetFloat64Slice(name).
func GetFloat64Slice(name string) []float64 {
	return Conf.GetFloat64Slice(name)
}

// MustFloat64Slice is equal to Conf.MustFloat64Slice(name).
func MustFloat64Slice(name string) []float64 {
	return Conf.MustFloat64Slice(name)
}

// GetInt is equal to Conf.GetInt(name).
func GetInt(name string) int {
	return Conf.GetInt(name)
}

// MustInt is equal to Conf.MustInt(name).
func MustInt(name string) int {
	return Conf.MustInt(name)
}

// GetInt32 is equal to Conf.GetInt32(name).
func GetInt32(name string) int32 {
	return Conf.GetInt32(name)
}

// MustInt32 is equal to Conf.MustInt32(name).
func MustInt32(name string) int32 {
	return Conf.MustInt32(name)
}

// GetInt64 is equal to Conf.GetInt64(name).
func GetInt64(name string) int64 {
	return Conf.GetInt64(name)
}

// MustInt64 is equal to Conf.MustInt64(name).
func MustInt64(name string) int64 {
	return Conf.MustInt64(name)
}

// GetIntSlice is equal to Conf.GetIntSlice(name).
func GetIntSlice(name string) []int {
	return Conf.GetIntSlice(name)
}

// MustIntSlice is equal to Conf.MustIntSlice(name).
func MustIntSlice(name string) []int {
	return Conf.MustIntSlice(name)
}

// GetString is equal to Conf.GetString(name).
func GetString(name string) string {
	return Conf.GetString(name)
}

// MustString is equal to Conf.MustString(name).
func MustString(name string) string {
	return Conf.MustString(name)
}

// GetStringSlice is equal to Conf.GetStringSlice(name).
func GetStringSlice(name string) []string {
	return Conf.GetStringSlice(name)
}

// MustStringSlice is equal to Conf.MustStringSlice(name).
func MustStringSlice(name string) []string {
	return Conf.MustStringSlice(name)
}

// GetTime is equal to Conf.GetTime(name).
func GetTime(name string) time.Time {
	return Conf.GetTime(name)
}

// MustTime is equal to Conf.MustTime(name).
func MustTime(name string) time.Time {
	return Conf.MustTime(name)
}

// GetUint is equal to Conf.GetUint(name).
func GetUint(name string) uint {
	return Conf.GetUint(name)
}

// MustUint is equal to Conf.MustUint(name).
func MustUint(name string) uint {
	return Conf.MustUint(name)
}

// GetUint32 is equal to Conf.GetUint32(name).
func GetUint32(name string) uint32 {
	return Conf.GetUint32(name)
}

// MustUint32 is equal to Conf.MustUint32(name).
func MustUint32(name string) uint32 {
	return Conf.MustUint32(name)
}

// GetUint64 is equal to Conf.GetUint64(name).
func GetUint64(name string) uint64 {
	return Conf.GetUint64(name)
}

// MustUint64 is equal to Conf.MustUint64(name).
func MustUint64(name string) uint64 {
	return Conf.MustUint64(name)
}

// GetUintSlice is equal to Conf.GetUintSlice(name).
func GetUintSlice(name string) []uint {
	return Conf.GetUintSlice(name)
}

// MustUintSlice is equal to Conf.MustUintSlice(name).
func MustUintSlice(name string) []uint {
	return Conf.MustUintSlice(name)
}
