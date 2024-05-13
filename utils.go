// Copyright 2019~2023 xgfone
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
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/xgfone/go-defaults"
)

// Some type converters, all of which have a default implementation,
// but you can reset them to yourself implementations.
var (
	ToBool     = defaults.ToBool     // func(interface{}) (bool, error)
	ToInt64    = defaults.ToInt64    // func(interface{}) (int64, error)
	ToUint64   = defaults.ToUint64   // func(interface{}) (uint64, error)
	ToFloat64  = defaults.ToFloat64  // func(interface{}) (float64, error)
	ToString   = defaults.ToString   // func(interface{}) (string, error)
	ToDuration = defaults.ToDuration // func(interface{}) (time.Duration, error)
	ToTime     = defaults.ToTime     // func(interface{}) (time.Time, error)
	ToInt      = toInt               // func(interface{}) (int, error)
	ToInt16    = toInt16             // func(interface{}) (int16, error)
	ToInt32    = toInt32             // func(interface{}) (int32, error)
	ToUint     = toUint              // func(interface{}) (uint, error)
	ToUint16   = toUint16            // func(interface{}) (uint16, error)
	ToUint32   = toUint32            // func(interface{}) (uint32, error)

	// For string type, it will be split by the separator " " or ",".
	ToIntSlice      = toIntSlice      // func(interface{}) ([]int, error)
	ToUintSlice     = toUintSlice     // func(interface{}) ([]uint, error)
	ToFloat64Slice  = toFloat64Slice  // func(interface{}) ([]float64, error)
	ToStringSlice   = toStringSlice   // func(interface{}) ([]string, error)
	ToDurationSlice = toDurationSlice // func(interface{}) ([]time.Duration, error)
)

func toInt(v interface{}) (int, error) {
	return to(v, ToInt64, func(v int64) int { return int(v) })
}
func toInt16(v interface{}) (int16, error) {
	return to(v, ToInt64, func(v int64) int16 { return int16(v) })
}
func toInt32(v interface{}) (int32, error) {
	return to(v, ToInt64, func(v int64) int32 { return int32(v) })
}
func toUint(v interface{}) (uint, error) {
	return to(v, ToInt64, func(v int64) uint { return uint(v) })
}
func toUint16(v interface{}) (uint16, error) {
	return to(v, ToInt64, func(v int64) uint16 { return uint16(v) })
}
func toUint32(v interface{}) (uint32, error) {
	return to(v, ToInt64, func(v int64) uint32 { return uint32(v) })
}

func to[T1, T2 any](i interface{}, f func(interface{}) (T1, error), m func(T1) T2) (v T2, err error) {
	_v, err := f(i)
	if err == nil {
		v = m(_v)
	}
	return
}

func bytesToMd5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func bytesToSha256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func isStringSeparator(r rune) bool {
	switch r {
	case ' ', ',', '\t':
		return true
	default:
		return false
	}
}

func getStringSlice(value interface{}) []string {
	var s string
	switch v := value.(type) {
	case []string:
		return v
	case string:
		s = v
	case []byte:
		s = string(v)
	case fmt.Stringer:
		s = v.String()
	default:
		return nil
	}

	vs := strings.FieldsFunc(s, isStringSeparator)
	ss := make([]string, 0, len(vs))
	for _, s := range vs {
		if s = strings.TrimSpace(s); s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}

func toSlice[E any](value interface{}, to func(interface{}) (E, error)) ([]E, error) {
	switch v := value.(type) {
	case nil:
		return []E{}, nil
	case []E:
		return v, nil
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		var err error
		vf := reflect.ValueOf(value)
		vs := make([]E, vf.Len())
		for i, _len := 0, vf.Len(); i < _len; i++ {
			if vs[i], err = to(vf.Index(i).Interface()); err != nil {
				return []E{}, fmt.Errorf("unable to cast %#v of type %T to []int", value, value)
			}
		}
		return vs, nil
	default:
		return []E{}, fmt.Errorf("unable to cast %#v of type %T to []int", value, value)
	}
}

func toIntSlice(value interface{}) ([]int, error) {
	if ss := getStringSlice(value); ss != nil {
		var err error
		vs := make([]int, len(ss))
		for i, s := range ss {
			if vs[i], err = ToInt(s); err != nil {
				return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", value, value)
			}
		}
		return vs, nil
	}
	return toSlice(value, ToInt)
}

func toUintSlice(value interface{}) (v []uint, err error) {
	if ss := getStringSlice(value); ss != nil {
		var err error
		vs := make([]uint, len(ss))
		for i, s := range ss {
			if vs[i], err = ToUint(s); err != nil {
				return []uint{}, fmt.Errorf("unable to cast %#v of type %T to []uint", value, value)
			}
		}
		return vs, nil
	}
	return toSlice(value, ToUint)
}

func toFloat64Slice(value interface{}) ([]float64, error) {
	if ss := getStringSlice(value); ss != nil {
		var err error
		vs := make([]float64, len(ss))
		for i, s := range ss {
			if vs[i], err = ToFloat64(s); err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", value, value)
			}
		}
		return vs, nil
	}
	return toSlice(value, ToFloat64)
}

func toStringSlice(value interface{}) ([]string, error) {
	if ss := getStringSlice(value); ss != nil {
		return ss, nil
	}
	return toSlice(value, ToString)
}

func toDurationSlice(value interface{}) ([]time.Duration, error) {
	if ss := getStringSlice(value); ss != nil {
		var err error
		vs := make([]time.Duration, len(ss))
		for i, s := range ss {
			if vs[i], err = ToDuration(s); err != nil {
				return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", value, value)
			}
		}
		return vs, nil
	}
	return toSlice(value, ToDuration)
}

func inString(s string, ss []string) bool {
	for _, _s := range ss {
		if _s == s {
			return true
		}
	}
	return false
}
