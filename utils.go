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

// Some converting function aliases.
var (
	IsZero    = types.IsZero
	ToBool    = types.ToBool
	ToInt64   = types.ToInt64
	ToUint64  = types.ToUint64
	ToFloat64 = types.ToFloat64
	ToString  = types.ToString
	ToTime    = types.ToTime
)

// ToDuration does the best to convert a certain value to time.Duration.
func ToDuration(v interface{}) (d time.Duration, err error) {
	switch _v := v.(type) {
	case string:
		return time.ParseDuration(_v)
	case []byte:
		return time.ParseDuration(string(_v))
	case fmt.Stringer:
		return time.ParseDuration(_v.String())
	}

	_v, err := ToInt64(v)
	if err != nil {
		return
	}
	return time.Duration(_v), nil
}

// ToStringSlice does the best to convert a certain value to []string.
//
// If the value is string, they are separated by the comma.
func ToStringSlice(_v interface{}) (v []string, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]string, 0, len(vs))
		for _, s := range vs {
			s = strings.TrimSpace(s)
			if s != "" {
				v = append(v, s)
			}
		}
	case []string:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToIntSlice does the best to convert a certain value to []int.
//
// If the value is string, they are separated by the comma.
func ToIntSlice(_v interface{}) (v []int, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]int, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := types.ToInt64(s)
			if err != nil {
				return nil, err
			}
			v = append(v, int(i))
		}
	case []int:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToInt64Slice does the best to convert a certain value to []int64.
//
// If the value is string, they are separated by the comma.
func ToInt64Slice(_v interface{}) (v []int64, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]int64, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := types.ToInt64(s)
			if err != nil {
				return nil, err
			}
			v = append(v, i)
		}
	case []int64:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToUintSlice does the best to convert a certain value to []uint.
//
// If the value is string, they are separated by the comma.
func ToUintSlice(_v interface{}) (v []uint, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]uint, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := types.ToUint64(s)
			if err != nil {
				return nil, err
			}
			v = append(v, uint(i))
		}
	case []uint:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToUint64Slice does the best to convert a certain value to []uint64.
//
// If the value is string, they are separated by the comma.
func ToUint64Slice(_v interface{}) (v []uint64, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]uint64, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := types.ToUint64(s)
			if err != nil {
				return nil, err
			}
			v = append(v, i)
		}
	case []uint64:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToFloat64Slice does the best to convert a certain value to []float64.
//
// If the value is string, they are separated by the comma.
func ToFloat64Slice(_v interface{}) (v []float64, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]float64, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := types.ToFloat64(s)
			if err != nil {
				return nil, err
			}
			v = append(v, i)
		}
	case []float64:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToTimes does the best to convert a certain value to []time.Time.
//
// If the value is string, they are separated by the comma and the each value
// is parsed by the format, layout.
func ToTimes(layout string, _v interface{}) (v []time.Time, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]time.Time, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := time.Parse(layout, s)
			if err != nil {
				return nil, err
			}
			v = append(v, i)
		}
	case []time.Time:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}

// ToDurations does the best to convert a certain value to []time.Duration.
//
// If the value is string, they are separated by the comma and the each value
// is parsed by time.ParseDuration().
func ToDurations(_v interface{}) (v []time.Duration, err error) {
	switch vv := _v.(type) {
	case string:
		vs := strings.Split(vv, ",")
		v = make([]time.Duration, 0, len(vs))
		for _, s := range vs {
			if s = strings.TrimSpace(s); s == "" {
				continue
			}

			i, err := time.ParseDuration(s)
			if err != nil {
				return nil, err
			}
			v = append(v, i)
		}
	case []time.Duration:
		v = vv
	default:
		err = types.ErrUnknownType
	}
	return
}
