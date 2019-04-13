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

///////////////////////////////////////////////////////////////////////////////
/// Config

// GetCliVersion returns the CLI version information.
//
// Return ("", "", "", "") if no version information.
//
// Notice: it should only be used by the CLI parser.
func (c *Config) GetCliVersion() (short, long, version, help string) {
	if c.vlong == "" || c.version == "" {
		return
	}
	return c.vshort, c.vlong, c.version, c.vhelp
}

// SetCliVersion sets the CLI version information.
func (c *Config) SetCliVersion(short, long, version, help string) *Config {
	if long == "" {
		panic("the long name of the version must not be empty")
	} else if len(short) > 1 {
		panic("the short name of the version must be empty or a character")
	} else if version == "" {
		panic("the version must not be empty")
	}

	c.panicIsParsed(true)
	c.vshort, c.vlong, c.vhelp, c.version = short, long, help, version
	return c
}

// ParsedCliArgs returns the parsed CLI arguments.
//
// Notice: for CLI parser, it should use this, not os.Args[1:].
func (c *Config) ParsedCliArgs() []string {
	return c.parsedCliArgs
}

// CliArgs returns the rest cli arguments parsed by the CLI parser.
//
// If no CLI parser or no rest cli arguments, it will return nil.
func (c *Config) CliArgs() []string {
	return c.cliArgs
}

// SetCliArgs sets the rest cli arguments, then you can call CliArgs() to get it.
//
// Notice: this method should only be called by the CLI parser.
func (c *Config) SetCliArgs(args []string) *Config {
	c.panicIsParsed(false)
	c.cliArgs = args
	return c
}

// RegisterCliStruct is equal to RegisterStruct, but the cli mode of the option
// is enabled by default.
func (c *Config) RegisterCliStruct(s interface{}) *Config {
	return c.registerStruct(true, s)
}

///////////////////////////////////////////////////////////////////////////////
/// Group

// RegisterCliOpts registers a set of CLI options into the current group.
func (g *OptGroup) RegisterCliOpts(opts []Opt) *OptGroup {
	for _, opt := range opts {
		g.RegisterCliOpt(opt)
	}
	return g
}

// RegisterCliOpt registers the CLI option into the current group.
func (g *OptGroup) RegisterCliOpt(opt Opt) *OptGroup {
	return g.registerOpt(true, opt)
}
