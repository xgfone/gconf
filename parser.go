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

// Parser is an parser interface to parse the configurations.
type Parser interface {
	// Name returns the name of the parser to identify it.
	Name() string

	// Priority reports the priority of the current parser, which should be
	// a natural number.
	//
	// The smaller the number, the higher the priority. And the higher priority
	// parser will be called to parse the option.
	//
	// For the cli parser, it maybe return 0 to indicate the highest priority.
	Priority() int

	// Pre is called before parsing the configuration, so it may be used to
	// initialize the parser, such as registering the itself options.
	Pre(*Config) error

	// Parse the value of the registered options.
	//
	// The parser can get any information from the argument, config.
	//
	// When the parser parsed out the option value, it should call
	// config.UpdateOptValue(), which will set the group option.
	// For the default group, the group name may be "" instead,
	//
	// For the CLI parser, it should get the parsed CLI argument by calling
	// config.ParsedCliArgs(), which is a string slice, not nil, but it maybe
	// have no elements. The CLI parser should not use os.Args[1:]
	// as the parsed CLI arguments. After parsing, If there are the rest CLI
	// arguments, which are those that does not start with the prefix "-", "--",
	// the CLI parser should call config.SetCliArgs() to set them.
	//
	// If there is any error, the parser should stop to parse and return it.
	//
	// If a certain option has no value, the parser should not return a default
	// one instead. Also, the parser has no need to convert the value to the
	// corresponding specific type, and just string is ok. Because the Config
	// will convert the value to the specific type automatically. Certainly,
	// it's not harmless for the parser to convert the value to the specific type.
	Parse(*Config) error

	// Pre is called before parsing the configuration, so it may be used to
	// clean the parser.
	Post(*Config) error
}

type parserOpt struct {
	Opt   Opt
	Group *OptGroup

	OptName   string
	GroupName string

	Value interface{}
	Other interface{}
}
