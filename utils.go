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
	"strings"
	"time"

	"github.com/xgfone/go-tools/types"
)

// Some type converters, all of which have a default implementation,
// but you can reset them to yourself implementations.
var (
	ToBool     func(interface{}) (bool, error)
	ToInt      func(interface{}) (int, error)
	ToInt32    func(interface{}) (int32, error)
	ToInt64    func(interface{}) (int64, error)
	ToUint     func(interface{}) (uint, error)
	ToUint32   func(interface{}) (uint32, error)
	ToUint64   func(interface{}) (uint64, error)
	ToFloat64  func(interface{}) (float64, error)
	ToString   func(interface{}) (string, error)
	ToDuration func(interface{}) (time.Duration, error)
	ToTime     func(interface{}) (time.Time, error)

	// For string type, it will be separated by the comma(,) by default.
	ToIntSlice      func(interface{}) ([]int, error)
	ToUintSlice     func(interface{}) ([]uint, error)
	ToFloat64Slice  func(interface{}) ([]float64, error)
	ToStringSlice   func(interface{}) ([]string, error)
	ToDurationSlice func(interface{}) ([]time.Duration, error)
)

func init() {
	ToBool = types.ToBool
	ToInt64 = types.ToInt64
	ToUint64 = types.ToUint64
	ToFloat64 = types.ToFloat64
	ToString = types.ToString
	ToTime = func(v interface{}) (time.Time, error) { return types.ToTime(v) }
	ToInt = types.ToInt
	ToInt32 = types.ToInt32
	ToUint = types.ToUint
	ToUint32 = types.ToUint32
	ToDuration = types.ToDuration

	ToIntSlice = toIntSlice
	ToUintSlice = toUintSlice
	ToFloat64Slice = toFloat64Slice
	ToStringSlice = toStringSlice
	ToDurationSlice = toDurationSlice
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

	vs := strings.Split(s, ",")
	ss := make([]string, 0, len(vs))
	for _, s := range vs {
		if s = strings.TrimSpace(s); s != "" {
			ss = append(ss, s)
		}
	}
	return ss
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
	return types.ToIntSlice(value)
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
	return types.ToUintSlice(value)
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
	return types.ToFloat64Slice(value)
}

func toStringSlice(value interface{}) ([]string, error) {
	if ss := getStringSlice(value); ss != nil {
		return ss, nil
	}
	return types.ToStringSlice(value)
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
	return types.ToDurationSlice(value)
}
