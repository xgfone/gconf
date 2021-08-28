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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
)

// Decoder is used to decode the configuration data.
type Decoder func(src []byte, dst map[string]interface{}) error

// AddDecoder adds a decoder, which will override it if it has been added.s
func (c *Config) AddDecoder(_type string, decoder Decoder) {
	c.decoders[strings.ToLower(_type)] = decoder
}

// GetDecoder returns the decoder by the type.
//
// Return nil if the decoder does not exist.
func (c *Config) GetDecoder(_type string) (decoder Decoder) {
	_type = strings.ToLower(_type)
	decoder, ok := c.decoders[_type]
	if !ok {
		if alias, ok := c.daliases[_type]; ok {
			decoder = c.decoders[alias]
		}
	}
	return
}

// AddDecoderTypeAliases adds the aliases of the decoder typed _type.
//
// For example,
//
//   c.AddDecoderTypeAliases("yaml", "yml")
//
// When acquiring the "yml" decoder and it does not exist, it will try to
// return the "yaml" decoder.
func (c *Config) AddDecoderTypeAliases(_type string, aliases ...string) {
	_type = strings.ToLower(_type)
	for _, alias := range aliases {
		c.daliases[strings.ToLower(alias)] = _type
	}
}

// NewJSONDecoder returns a json decoder to decode the json data.
//
// If the json data contains the comment line starting with "//", it will remove
// the comment line and parse the json data.
func NewJSONDecoder() Decoder {
	comment := []byte("//")
	newline := []byte("\n")
	return func(src []byte, dst map[string]interface{}) (err error) {
		if bytes.Contains(src, comment) {
			buf := bytes.NewBuffer(nil)
			buf.Grow(len(src))
			for _, line := range bytes.Split(src, newline) {
				if line = bytes.TrimSpace(line); len(line) == 0 {
					buf.WriteByte('\n')
					continue
				} else if len(line) > 1 && line[0] == '/' && line[1] == '/' {
					continue
				}
				buf.Write(line)
				buf.WriteByte('\n')
			}
			src = buf.Bytes()
		}
		return json.Unmarshal(src, &dst)
	}
}

// NewYamlDecoder returns a yaml decoder to decode the yaml data.
func NewYamlDecoder() Decoder {
	return func(src []byte, dst map[string]interface{}) (err error) {
		return yaml.Unmarshal([]byte(src), &dst)
	}
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

	return func(src []byte, dst map[string]interface{}) (err error) {
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
					value = strings.TrimSpace(lines[index])

					var goon bool
					if _len := len(value) - 1; value[_len] == '\\' {
						goon = true
					}

					if value = strings.TrimSpace(strings.TrimRight(value, "\\")); value == "" {
						break
					}
					index++
					vs = append(vs, value)

					if !goon {
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
	}
}
