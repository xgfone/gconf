// Copyright 2020 xgfone
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

package field

import (
	"sync"
	"time"

	"github.com/xgfone/gconf/v4"
)

// SafeValue is used to set and get the value safely.
//
// Notice: it is embedded in OptField in general.
type SafeValue struct {
	value interface{}
	lock  sync.RWMutex
}

// Get returns the value safely.
func (sv *SafeValue) Get(_default interface{}) (value interface{}) {
	sv.lock.RLock()
	if value = sv.value; value == nil && _default != nil {
		value = _default
	}
	sv.lock.RUnlock()
	return
}

// Set updates the value to v safely.
func (sv *SafeValue) Set(v interface{}) {
	sv.lock.Lock()
	sv.value = v
	sv.lock.Unlock()
}

// BoolOptField represents the bool option field of the struct.
//
// The default is false.
type BoolOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *BoolOptField) Default() interface{} {
	return false
}

// Parse implements OptField.Parse().
func (f *BoolOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToBool(input)
}

// Set implements OptField.Set().
func (f *BoolOptField) Set(v interface{}) {
	f.value.Set(v.(bool))
}

// Get returns the value of the option field.
func (f *BoolOptField) Get() bool {
	return f.value.Get(false).(bool)
}

// BoolTOptField represents the bool option field of the struct.
//
// The default is true.
type BoolTOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *BoolTOptField) Default() interface{} {
	return true
}

// Parse implements OptField.Parse().
func (f *BoolTOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToBool(input)
}

// Set implements OptField.Set().
func (f *BoolTOptField) Set(v interface{}) {
	f.value.Set(v.(bool))
}

// Get returns the value of the option field.
func (f *BoolTOptField) Get() bool {
	return f.value.Get(true).(bool)
}

// IntOptField represents the int option field of the struct.
type IntOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *IntOptField) Default() interface{} {
	return 0
}

// Parse implements OptField.Parse().
func (f *IntOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToInt(input)
}

// Set implements OptField.Set().
func (f *IntOptField) Set(v interface{}) {
	f.value.Set(v.(int))
}

// Get returns the value of the option field.
func (f *IntOptField) Get() int {
	return f.value.Get(0).(int)
}

// Int32OptField represents the int32 option field of the struct.
type Int32OptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Int32OptField) Default() interface{} {
	return int32(0)
}

// Parse implements OptField.Parse().
func (f *Int32OptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToInt32(input)
}

// Set implements OptField.Set().
func (f *Int32OptField) Set(v interface{}) {
	f.value.Set(v.(int32))
}

// Get returns the value of the option field.
func (f *Int32OptField) Get() int32 {
	return f.value.Get(int32(0)).(int32)
}

// Int64OptField represents the int64 option field of the struct.
type Int64OptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Int64OptField) Default() interface{} {
	return int64(0)
}

// Parse implements OptField.Parse().
func (f *Int64OptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToInt64(input)
}

// Set implements OptField.Set().
func (f *Int64OptField) Set(v interface{}) {
	f.value.Set(v.(int64))
}

// Get returns the value of the option field.
func (f *Int64OptField) Get() int64 {
	return f.value.Get(int64(0)).(int64)
}

// UintOptField represents the uint option field of the struct.
type UintOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *UintOptField) Default() interface{} {
	return uint(0)
}

// Parse implements OptField.Parse().
func (f *UintOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToUint(input)
}

// Set implements OptField.Set().
func (f *UintOptField) Set(v interface{}) {
	f.value.Set(v.(uint))
}

// Get returns the value of the option field.
func (f *UintOptField) Get() uint {
	return f.value.Get(uint(0)).(uint)
}

// Uint32OptField represents the uint32 option field of the struct.
type Uint32OptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Uint32OptField) Default() interface{} {
	return uint32(0)
}

// Parse implements OptField.Parse().
func (f *Uint32OptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToUint32(input)
}

// Set implements OptField.Set().
func (f *Uint32OptField) Set(v interface{}) {
	f.value.Set(v.(uint32))
}

// Get returns the value of the option field.
func (f *Uint32OptField) Get() uint32 {
	return f.value.Get(uint32(0)).(uint32)
}

// Uint64OptField represents the uint64 option field of the struct.
type Uint64OptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Uint64OptField) Default() interface{} {
	return uint64(0)
}

// Parse implements OptField.Parse().
func (f *Uint64OptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToUint64(input)
}

// Set implements OptField.Set().
func (f *Uint64OptField) Set(v interface{}) {
	f.value.Set(v.(uint64))
}

// Get returns the value of the option field.
func (f *Uint64OptField) Get() uint64 {
	return f.value.Get(uint64(0)).(uint64)
}

// Float64OptField represents the float64 option field of the struct.
type Float64OptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Float64OptField) Default() interface{} {
	return float64(0)
}

// Parse implements OptField.Parse().
func (f *Float64OptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToFloat64(input)
}

// Set implements OptField.Set().
func (f *Float64OptField) Set(v interface{}) {
	f.value.Set(v.(float64))
}

// Get returns the value of the option field.
func (f *Float64OptField) Get() float64 {
	return f.value.Get(float64(0)).(float64)
}

// StringOptField represents the string option field of the struct.
type StringOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *StringOptField) Default() interface{} {
	return ""
}

// Parse implements OptField.Parse().
func (f *StringOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToString(input)
}

// Set implements OptField.Set().
func (f *StringOptField) Set(v interface{}) {
	f.value.Set(v.(string))
}

// Get returns the value of the option field.
func (f *StringOptField) Get() string {
	return f.value.Get("").(string)
}

// DurationOptField represents the time.Duration option field of the struct.
type DurationOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *DurationOptField) Default() interface{} {
	return time.Duration(0)
}

// Parse implements OptField.Parse().
func (f *DurationOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToDuration(input)
}

// Set implements OptField.Set().
func (f *DurationOptField) Set(v interface{}) {
	f.value.Set(v.(time.Duration))
}

// Get returns the value of the option field.
func (f *DurationOptField) Get() time.Duration {
	return f.value.Get(time.Duration(0)).(time.Duration)
}

// TimeOptField represents the time.Time option field of the struct.
type TimeOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *TimeOptField) Default() interface{} {
	return time.Time{}
}

// Parse implements OptField.Parse().
func (f *TimeOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToTime(input)
}

// Set implements OptField.Set().
func (f *TimeOptField) Set(v interface{}) {
	f.value.Set(v.(time.Time))
}

// Get returns the value of the option field.
func (f *TimeOptField) Get() time.Time {
	return f.value.Get(time.Time{}).(time.Time)
}

// IntSliceOptField represents the []int option field of the struct.
type IntSliceOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *IntSliceOptField) Default() interface{} {
	return []int{}
}

// Parse implements OptField.Parse().
func (f *IntSliceOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToIntSlice(input)
}

// Set implements OptField.Set().
func (f *IntSliceOptField) Set(v interface{}) {
	f.value.Set(v.([]int))
}

// Get returns the value of the option field.
func (f *IntSliceOptField) Get() []int {
	return f.value.Get([]int{}).([]int)
}

// UintSliceOptField represents the []uint option field of the struct.
type UintSliceOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *UintSliceOptField) Default() interface{} {
	return []uint{}
}

// Parse implements OptField.Parse().
func (f *UintSliceOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToUintSlice(input)
}

// Set implements OptField.Set().
func (f *UintSliceOptField) Set(v interface{}) {
	f.value.Set(v.([]uint))
}

// Get returns the value of the option field.
func (f *UintSliceOptField) Get() []uint {
	return f.value.Get([]uint{}).([]uint)
}

// Float64SliceOptField represents the []float64 option field of the struct.
type Float64SliceOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *Float64SliceOptField) Default() interface{} {
	return []float64{}
}

// Parse implements OptField.Parse().
func (f *Float64SliceOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToFloat64Slice(input)
}

// Set implements OptField.Set().
func (f *Float64SliceOptField) Set(v interface{}) {
	f.value.Set(v.([]float64))
}

// Get returns the value of the option field.
func (f *Float64SliceOptField) Get() []float64 {
	return f.value.Get([]float64{}).([]float64)
}

// StringSliceOptField represents the []string option field of the struct.
type StringSliceOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *StringSliceOptField) Default() interface{} {
	return []string{}
}

// Parse implements OptField.Parse().
func (f *StringSliceOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToStringSlice(input)
}

// Set implements OptField.Set().
func (f *StringSliceOptField) Set(v interface{}) {
	f.value.Set(v.([]string))
}

// Get returns the value of the option field.
func (f *StringSliceOptField) Get() []string {
	return f.value.Get([]string{}).([]string)
}

// DurationSliceOptField represents the []time.Duration option field of the struct.
type DurationSliceOptField struct {
	value SafeValue
}

// Default implements OptField.Default().
func (f *DurationSliceOptField) Default() interface{} {
	return []time.Duration{}
}

// Parse implements OptField.Parse().
func (f *DurationSliceOptField) Parse(input interface{}) (output interface{}, err error) {
	return gconf.ToDurationSlice(input)
}

// Set implements OptField.Set().
func (f *DurationSliceOptField) Set(v interface{}) {
	f.value.Set(v.([]time.Duration))
}

// Get returns the value of the option field.
func (f *DurationSliceOptField) Get() []time.Duration {
	return f.value.Get([]time.Duration{}).([]time.Duration)
}
