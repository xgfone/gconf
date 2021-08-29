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

// GetAllOpts is equal to Conf.GetAllOpts().
func GetAllOpts() []Opt { return Conf.GetAllOpts() }

// RegisterOpts is equal to Conf.RegisterOpts(opts...).
func RegisterOpts(opts ...Opt) { Conf.RegisterOpts(opts...) }

// UnregisterOpts is equal to Conf.UnregisterOpts(optNames...).
func UnregisterOpts(optNames ...string) { Conf.UnregisterOpts(optNames...) }

// SetVersion is equal to Conf.SetVersion(version).
func SetVersion(version string) { Conf.SetVersion(version) }

// LoadBackupFile is equal to Conf.LoadBackupFile().
func LoadBackupFile(filename string) error {
	return Conf.LoadBackupFile(filename)
}

// Snapshot is equal to Conf.Snapshot().
func Snapshot() (generation uint64, snap map[string]interface{}) {
	return Conf.Snapshot()
}

// LoadMap is equal to Conf.LoadMap(options, force...).
func LoadMap(options map[string]interface{}, force ...bool) error {
	return Conf.LoadMap(options)
}

// Set is equal to Conf.Set(name, value).
func Set(name string, value interface{}) error { return Conf.Set(name, value) }

// Get is equal to Conf.Get(name).
func Get(name string) interface{} { return Conf.Get(name) }

// Must is equal to Conf.Must(name).
func Must(name string) interface{} { return Conf.Must(name) }

// Observe is equal to Conf.Observe(observers...).
func Observe(observers ...Observer) { Conf.Observe(observers...) }

// GetGroupSep is equal to Conf.GetGroupSep().
func GetGroupSep() (sep string) { return Conf.GetGroupSep() }

// GetBool is equal to Conf.GetBool(name).
func GetBool(name string) bool { return Conf.GetBool(name) }

// GetInt is equal to Conf.GetInt(name).
func GetInt(name string) int { return Conf.GetInt(name) }

// GetInt32 is equal to Conf.GetInt32(name).
func GetInt32(name string) int32 { return Conf.GetInt32(name) }

// GetInt64 is equal to Conf.GetInt64(name).
func GetInt64(name string) int64 { return Conf.GetInt64(name) }

// GetUint is equal to Conf.GetUint(name).
func GetUint(name string) uint { return Conf.GetUint(name) }

// GetUint32 is equal to Conf.GetUint32(name).
func GetUint32(name string) uint32 { return Conf.GetUint32(name) }

// GetUint64 is equal to Conf.GetUint64(name).
func GetUint64(name string) uint64 { return Conf.GetUint64(name) }

// GetFloat64 is equal to Conf.GetFloat64(name).
func GetFloat64(name string) float64 { return Conf.GetFloat64(name) }

// GetString is equal to Conf.GetString(name).
func GetString(name string) string { return Conf.GetString(name) }

// GetDuration is equal to Conf.GetDuration(name).
func GetDuration(name string) time.Duration { return Conf.GetDuration(name) }

// GetTime is equal to Conf.GetTime(name).
func GetTime(name string) time.Time { return Conf.GetTime(name) }

// GetIntSlice is equal to Conf.GetIntSlice(name).
func GetIntSlice(name string) []int { return Conf.GetIntSlice(name) }

// GetUintSlice is equal to Conf.GetUintSlice(name).
func GetUintSlice(name string) []uint { return Conf.GetUintSlice(name) }

// GetFloat64Slice is equal to Conf.GetFloat64Slice(name).
func GetFloat64Slice(name string) []float64 { return Conf.GetFloat64Slice(name) }

// GetStringSlice is equal to Conf.GetStringSlice(name).
func GetStringSlice(name string) []string { return Conf.GetStringSlice(name) }

// GetDurationSlice is equal to Conf.GetDurationSlice(name).
func GetDurationSlice(name string) []time.Duration { return Conf.GetDurationSlice(name) }
