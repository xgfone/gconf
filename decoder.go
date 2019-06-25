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
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
)

// Decoder is used to decode the configuration data.
type Decoder struct {
	// Type is the type of decoder, such as "json", "xml", which is case insensitive.
	Type string

	// Decode is used to decode the configuration data.
	//
	// The decoder maybe decode the src data to the formats as follow:
	//
	//     map[string]interface{} {
	//         "opt1": "value1",
	//         "opt2": "value2",
	//         // ...
	//         "group1": map[string]interface{} {
	//             "opt11": "value11",
	//             "opt12": "value12",
	//             "group2": map[string]interface{} {
	//                 // ...
	//             },
	//             "group3.group4": map[string]interface{} {
	//                 // ...
	//             }
	//         },
	//         "group5.group6.group7": map[string]interface{} {
	//             "opt71": "value71",
	//             "opt72": "value72",
	//             "group8": map[string]interface{} {
	//                 // ...
	//             },
	//             "group9.group10": map[string]interface{} {
	//                 // ...
	//             }
	//         },
	//         "group11.group12.opt121": "value121",
	//         "group11.group12.opt122": "value122"
	//     }
	//
	// When loading it, it will be flatted to
	//
	//     map[string]interface{} {
	//         "opt1": "value1",
	//         "opt2": "value2",
	//         "group1.opt1": "value11",
	//         "group1.opt2": "value12",
	//         "group1.group2.XXX": "XXX",
	//         "group1.group3.group4.XXX": "XXX",
	//         "group5.group6.group7.opt71": "value71",
	//         "group5.group6.group7.opt72": "value72",
	//         "group5.group6.group7.group8.XXX": "XXX",
	//         "group5.group6.group7.group9.group10.XXX": "XXX",
	//         "group11.group12.opt121": "value121",
	//         "group11.group12.opt122": "value122"
	//     }
	//
	// So the option name must not contain the dot(.).
	Decode func(src []byte, dst map[string]interface{}) error
}

// AddDecoder adds a decoder and returns true.
//
// If the decoder has added, it will do nothing and return false.
// But you can override it by setting force to true.
func (c *Config) AddDecoder(decoder Decoder, force ...bool) (ok bool) {
	_type := strings.ToLower(decoder.Type)
	c.lock.Lock()
	if _, ok = c.decoders[_type]; !ok {
		c.decoders[_type] = decoder
		ok = true
	}
	c.lock.Unlock()
	return
}

// GetDecoder returns the decoder by the type.
func (c *Config) GetDecoder(_type string) (decoder Decoder, ok bool) {
	_type = strings.ToLower(_type)
	c.lock.RLock()
	if decoder, ok = c.decoders[_type]; !ok {
		if alias, _ok := c.decAlias[_type]; _ok {
			decoder, ok = c.decoders[alias]
		}
	}
	c.lock.RUnlock()
	return
}

// AddDecoderAlias adds the alias of the decoder typed _type. For example,
//
//   c.AddDecoderAlias("conf", "ini")
//
// When you get the "conf" decoder and it does not exist, it will try to
// return the "ini" decoder.
//
// If the alias has existed, it will override it.
func (c *Config) AddDecoderAlias(_type, alias string) {
	_type = strings.ToLower(_type)
	alias = strings.ToLower(alias)

	c.lock.Lock()
	c.decAlias[_type] = alias
	c.lock.Unlock()
}

// NewDecoder returns a new decoder.
func NewDecoder(_type string, decode func([]byte, map[string]interface{}) error) Decoder {
	return Decoder{Type: _type, Decode: decode}
}

// NewJSONDecoder returns a json decoder to decode the json data.
func NewJSONDecoder() Decoder {
	return NewDecoder("json", func(src []byte, dst map[string]interface{}) (err error) {
		return json.Unmarshal(src, &dst)
	})
}

// NewYamlDecoder returns a yaml decoder to decode the yaml data.
func NewYamlDecoder() Decoder {
	return NewDecoder("yaml", func(src []byte, dst map[string]interface{}) (err error) {
		return yaml.Unmarshal([]byte(src), &dst)
	})
}

// NewIniDecoder returns a INI decoder to decode the INI data.
//
// Notice:
//   1. The empty line will be ignored.
//   2. The spacewhite on the beginning and end of line or value will be trimmed.
//   3. The comment line starts with the character '#' or ';', which is ignored.
//   4. The name of the default group is "DEFAULT", but it is optional.
//   5. The group can nest other groups by ".", such as "group1.group2.group3".
//   6. The key must only contain the printable non-spacewhite characters.
//   7. The line can continue to the next line with the last character "\",
//      and the spacewhite on the beginning and end of the each line will be
//      trimmed, then combines them with a space.
//
func NewIniDecoder(defaultGroupName ...string) Decoder {
	defaultGroup := "DEFAULT"
	if len(defaultGroupName) > 0 && defaultGroupName[0] != "" {
		defaultGroup = defaultGroupName[0]
	}

	return NewDecoder("ini", func(src []byte, dst map[string]interface{}) (err error) {
		var gname string
		lines := strings.Split(string(src), "\n")
		for index, maxIndex := 0, len(lines); index < maxIndex; {
			line := strings.TrimSpace(lines[index])
			index++

			// Ignore the empty line and the comment line
			if len(line) == 0 || line[0] == '#' || line[0] == ';' {
				continue
			}

			// Start a new group
			if last := len(line) - 1; line[0] == '[' && line[last] == ']' {
				gname = strings.TrimSpace(line[1:last])
				if gname == defaultGroup {
					gname = ""
				}
				continue
			}

			n := strings.IndexByte(line, '=')
			if n < 0 {
				return fmt.Errorf("the %dth line misses the separator '='", index)
			}

			// Get the key
			key := strings.TrimSpace(line[:n])
			if len(key) == 0 {
				return fmt.Errorf("empty identifier key")
			}
			for _, r := range key {
				if unicode.IsSpace(r) || !unicode.IsPrint(r) {
					return fmt.Errorf("invalid identifier key '%s'", key)
				}
			}

			// Get the value
			value := strings.TrimSpace(line[n+1:])
			if value == "" { // Ignore the empty value
				continue
			} else if _len := len(value) - 1; value[_len] == '\\' { // The continuation line
				vs := []string{strings.TrimSpace(strings.TrimRight(value, "\\"))}
				for index < maxIndex {
					value = strings.TrimRight(strings.TrimSpace(lines[index]), "\\")
					if value = strings.TrimSpace(value); value == "" {
						break
					}
					index++
					vs = append(vs, value)
					if value[len(value)-1] != '\\' {
						break
					}
				}
				value = strings.Join(vs, " ")
			}

			// Add the option
			if gname != "" {
				key = strings.Join([]string{gname, key}, ".")
			}
			dst[key] = value
		}
		return
	})
}
