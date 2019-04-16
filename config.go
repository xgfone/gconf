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
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	// ErrParsed is an error that the config has been parsed.
	ErrParsed = fmt.Errorf("the config manager has been parsed")

	// ErrNotParsed is an error that the config has not been parsed.
	ErrNotParsed = fmt.Errorf("the config manager has not been parsed")
)

// DefaultGroupSeparator is the default separator between the group names.
var DefaultGroupSeparator = "."

// Config is used to manage the configuration options.
type Config struct {
	*OptGroup // The default group.

	lock   sync.RWMutex
	parsed int32

	// Common Settings
	stop       bool
	zero       bool
	debug      bool
	required   bool
	reregister bool
	gsep       string // The separator between the group names.
	prefix     string // ==> OptGroup.name + "."
	printf     func(string, ...interface{})

	// CLI
	name          string
	help          string
	vshort        string
	vlong         string
	vhelp         string
	version       string
	cliArgs       []string
	parsedCliArgs []string

	parsers    []Parser
	executed   *Command
	actions    map[string]func() error
	commands   map[string]*Command
	allGroups  map[string]*OptGroup
	validators []func() error
	observe    func(string, string, interface{}, interface{})
	action     func() error
	check      func() error
}

// New is equal to NewConfig("", "").
func New() *Config {
	return NewConfig("", "")
}

// NewConfig returns a new Config with the name and the description of Config.
//
// If the name is "", it will be os.Args[0] by default.
//
// The name and the description are used as the name and the usage of the program
// by the CLI parser in general.
func NewConfig(name, description string) *Config {
	if name == "" {
		name = filepath.Base(os.Args[0])
	}

	c := &Config{
		name:   name,
		help:   description,
		printf: func(f string, ss ...interface{}) { fmt.Printf(f+"\n", ss...) },

		actions:    make(map[string]func() error),
		commands:   make(map[string]*Command),
		allGroups:  make(map[string]*OptGroup),
		validators: make([]func() error, 0, 8),
	}

	c.OptGroup = newOptGroup2(false, c, nil, DefaultGroupName)
	c.SetGroupSeparator(DefaultGroupSeparator)
	c.SetCheckRequiredOption(c.checkRequiredOption)
	c.noticeNewGroup(c.OptGroup)
	return c
}

// Name returns the config name.
func (c *Config) Name() string {
	return c.name
}

// Description returns the config description.
func (c *Config) Description() string {
	return c.help
}

func (c *Config) mergePaths(paths []string) string {
	return strings.TrimPrefix(strings.Join(paths, c.gsep), c.prefix)
}

func (c *Config) getFullGroupName(parent, name string) string {
	if parent == "" {
		return name
	}
	return strings.TrimPrefix(parent+c.gsep+name, c.prefix)
}

//////////////////////////////////////////////////////////////////////////////
/// Setting

// GetDefaultGroupName returns the name of the default group.
func (c *Config) GetDefaultGroupName() string {
	return c.OptGroup.name
}

// SetDefaultGroupName resets the name of the default group.
//
// If parsed, it will panic when calling it.
func (c *Config) SetDefaultGroupName(name string) *Config {
	c.panicIsParsed(true)
	if name == "" {
		name = DefaultGroupName
	}
	if c.OptGroup.name != name {
		c.OptGroup.name = name
		c.prefix = name + c.gsep
	}
	return c
}

// GetGroupSeparator returns the separator between the group names.
func (c *Config) GetGroupSeparator() string {
	return c.gsep
}

// SetGroupSeparator sets the separator between the group names.
//
// The default separator is a dot(.).
//
// If parsed, it will panic when calling it.
func (c *Config) SetGroupSeparator(sep string) *Config {
	if sep == "" {
		panic("the separator must not be empty")
	}

	c.panicIsParsed(true)
	c.gsep = sep
	c.prefix = c.GetDefaultGroupName() + c.gsep
	return c
}

// IsDebug reports whether the debug mode is enabled.
func (c *Config) IsDebug() bool {
	return c.debug
}

// SetDebug enables or disables the debug model.
//
// If setting, when registering the option, it'll output the verbose information.
// You should set it before registering any options.
//
// If parsed, it will panic when calling it.
func (c *Config) SetDebug(debug bool) *Config {
	c.panicIsParsed(true)
	c.debug = debug
	return c
}

// SetZero enables to set the value of the option to the zero value of its type
// if the option has no value.
//
// If parsed, it will panic when calling it.
func (c *Config) SetZero(zero bool) *Config {
	c.panicIsParsed(true)
	c.zero = zero
	return c
}

// SetPrintf sets the printf function, which should append a newline
// after output, to print the debug log.
//
// The default printf is equal to `fmt.Printf(msg+"\n", args...)`.
func (c *Config) SetPrintf(printf func(msg string, args ...interface{})) *Config {
	if printf == nil {
		panic("the printf must not be nil")
	}

	c.panicIsParsed(true)
	c.printf = printf
	return c
}

// Printf prints the log messages if enabling debug.
//
// It is output to os.Stdout by default, and append a newline, see SetPrintf().
func (c *Config) Printf(format string, args ...interface{}) {
	if c.debug {
		c.printf(format, args...)
	}
}

// SetRequired asks that all the registered options have a value.
//
// Notice: the nil value is not considered that there is a value, but the ZERO
// value is that.
//
// If parsed, it will panic when calling it.
func (c *Config) SetRequired(required bool) *Config {
	c.panicIsParsed(true)
	c.required = required
	return c
}

// SetCheckRequiredOption sets the check function for the required options.
//
// It will check the options from all non-command groups and all the groups
// of the executed commands.
//
// The default is enough.
func (c *Config) SetCheckRequiredOption(check func() error) *Config {
	if check == nil {
		panic("the check function must not be nil")
	}

	c.panicIsParsed(true)
	c.check = check
	return c
}

// IgnoreReregister decides whether it will panic when reregistering an option
// into a certain group.
//
// The default is not to ignore it, but you can set it to false to ignore it.
func (c *Config) IgnoreReregister(ignore bool) *Config {
	c.panicIsParsed(true)
	c.reregister = ignore
	return c
}

//////////////////////////////////////////////////////////////////////////////
/// Parse

func (c *Config) panicIsParsed(p bool) {
	if c.Parsed() {
		if p {
			panic(ErrParsed)
		}
	} else {
		if !p {
			panic(ErrNotParsed)
		}
	}
}

// Parsed reports whether the config has been parsed.
func (c *Config) Parsed() bool {
	return atomic.LoadInt32(&c.parsed) == 1
}

// Stop stops the subsequent parsing.
//
// In general, it is used by the parser to stop the subsequent operation
// after its Parse() is called.
func (c *Config) Stop() *Config {
	c.panicIsParsed(false)
	c.stop = true
	return c
}

// Parse parses the options, including CLI, the config file, or others.
//
// if no any arguments, it's equal to os.Args[1:].
//
// After parsing a certain option, it will call the validators of the option
// to validate whether the option value is valid.
//
// If parsed, it will panic when calling it.
func (c *Config) Parse(args ...string) (err error) {
	c.panicIsParsed(true)

	if args == nil {
		c.parsedCliArgs = os.Args[1:]
	} else {
		c.parsedCliArgs = args
	}

	// Preprocess the parsers.
	for _, parser := range c.parsers {
		c.Printf("Initializing the parser '%s'", parser.Name())
		if err = parser.Pre(c); err != nil {
			return err
		}
	}

	// Set to have been parsed.
	atomic.StoreInt32(&c.parsed, 1)

	// Call the parsers to parse the options.
	for _, parser := range c.parsers {
		c.Printf("Calling the parser '%s'", parser.Name())
		if err = parser.Parse(c); err != nil {
			return fmt.Errorf("The '%s' parser failed: %s", parser.Name(), err)
		}
		if c.stop {
			break
		}
	}

	// Postprocess the parsers.
	for _, parser := range c.parsers {
		c.Printf("Cleaning the parser '%s'", parser.Name())
		if err = parser.Post(c); err != nil {
			return err
		}
	}

	if !c.stop {
		// Check whether all the groups have parsed all the required options.
		if err = c.check(); err != nil {
			return
		}
	}

	return
}

// CheckRequiredOption check whether all the required options have an value.
func (c *Config) CheckRequiredOption() error {
	return c.check()
}

func (c *Config) checkRequiredOptionByCmd(cmd *Command) (err error) {
	if cmd != nil {
		for _, group := range cmd.AllGroups() {
			c.Printf("Checking the required options from the group '%s' of the command '%s'",
				group.FullName(), cmd.FullName())
			if err = group.CheckRequiredOption(); err != nil {
				return
			}
		}
		return c.checkRequiredOptionByCmd(cmd.ParentCommand())
	}
	return
}

func (c *Config) checkRequiredOption() (err error) {
	for _, group := range c.AllNotCommandGroups() {
		c.Printf("Check the required options for the global group '%s'", group.FullName())
		if err = group.CheckRequiredOption(); err != nil {
			return
		}
	}

	if err = c.checkRequiredOptionByCmd(c.ExecutedCommand()); err != nil {
		return
	}

	for _, vf := range c.validators {
		if err = vf(); err != nil {
			return err
		}
	}

	return
}

//////////////////////////////////////////////////////////////////////////////
/// Set the option value and Observe the change of the option value

func (c *Config) watchChangedOption(group *OptGroup, opt string, old, new interface{}) {
	c.Printf("Set [%s]:[%s] from [%v] to [%v]", group.fname, opt, old, new)
	if c.observe != nil {
		c.observe(group.fname, opt, old, new)
	}
}

// Observe watches the change of values.
//
// When the option value is changed, the function f will be called.
//
// If SetOptValue() is used in the multi-thread, you should promise
// that the callback function f is thread-safe and reenterable.
//
// Notice: you can get the group by calling `config.Group(groupFullName)``
// and the option by calling `config.Group(groupFullName).Opt(optName)``.
func (c *Config) Observe(f func(groupFullName, optName string, oldOptValue, newOptValue interface{})) {
	c.panicIsParsed(true)
	c.observe = f
}

// SetOptValue sets the value of the option in the group. It's thread-safe.
//
// "priority" should be the priority of the parser. It only set the option value
// successfully for the priority higher than the last. So you can use 0
// to update it coercively.
//
// For the command or multi-groups, you should unite them using the separator.
// the command itself is considered as a group, for example,
//
//     "Group1.SubGroup2.SubSubGroup3"
//     "Command.SubGroup1.SubSubGroup2"
//     "Command1.SubCommand2.SubGroup1.SubGroup2"
//
// Notice: You cannot call SetOptValue() for the struct option, because we have
// no way to promise that it's thread-safe.
func (c *Config) SetOptValue(priority int, groupFullName, optName string, optValue interface{}) error {
	c.panicIsParsed(false)
	if priority < 0 {
		return fmt.Errorf("the priority must not be the negative")
	}

	if groupFullName == "" {
		groupFullName = c.OptGroup.name
	}

	if group := c.allGroups[groupFullName]; group != nil {
		return group.setOptValue(priority, optName, optValue)
	}
	return fmt.Errorf("no group '%s'", groupFullName)
}

//////////////////////////////////////////////////////////////////////////////
/// Parser

func (c *Config) sortParsers() {
	sort.SliceStable(c.parsers, func(i, j int) bool {
		return c.parsers[i].Priority() < c.parsers[j].Priority()
	})
}

// AddParser adds a few parsers.
func (c *Config) AddParser(parsers ...Parser) *Config {
	c.panicIsParsed(true)
	c.parsers = append(c.parsers, parsers...)
	c.sortParsers()
	return c
}

// RemoveParser removes and returns the parser named name.
//
// Return nil if the parser does not exist.
func (c *Config) RemoveParser(name string) Parser {
	c.panicIsParsed(true)
	for i, p := range c.parsers {
		if p.Name() == name {
			ps := make([]Parser, 0, len(c.parsers)-1)
			ps = append(ps, c.parsers[:i]...)
			ps = append(ps, c.parsers[i:]...)
			c.parsers = ps
			return p
		}
	}
	return nil
}

// GetParser returns the parser named name.
//
// Return nil if the parser does not exist.
func (c *Config) GetParser(name string) Parser {
	for _, p := range c.parsers {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// HasParser reports whether the parser named name exists.
func (c *Config) HasParser(name string) bool {
	return c.GetParser(name) != nil
}

// Parsers returns all the parsers.
func (c *Config) Parsers() []Parser {
	return append([]Parser{}, c.parsers...)
}

//////////////////////////////////////////////////////////////////////////////
/// Action

// Action returns the action function of Config.
//
// In general, it is used by the CLI parser.
func (c *Config) Action() func() error {
	return c.action
}

// SetAction sets the action function for Config.
func (c *Config) SetAction(action func() error) *Config {
	if action == nil {
		panic("the action must not be nil")
	}
	c.panicIsParsed(true)
	c.action = action
	return c
}

// RegisterAction registers a action of the command with the name.
//
// It may be used by the struct tag. See Config.RegisterStruct().
func (c *Config) RegisterAction(name string, action func() error) *Config {
	if name == "" {
		panic("the action name must not be empty")
	} else if action == nil {
		panic("the action must not be nil")
	}

	c.panicIsParsed(true)
	c.actions[name] = action
	c.Printf("Register the action '%s'", name)
	return c
}

// GetAction returns the action function by the name.
//
// Return nil if no action function.
func (c *Config) GetAction(name string) func() error {
	return c.actions[name]
}

//////////////////////////////////////////////////////////////////////////////
/// Command

// NewCommand news a Command to register the CLI sub-command.
//
// Notice: if the command exists, it returns the old, not a new one.
func (c *Config) NewCommand(name, help string) (cmd *Command) {
	c.panicIsParsed(true)
	if cmd = c.commands[name]; cmd == nil {
		cmd = newCommand(c, nil, name, help)
		c.commands[name] = cmd
	}
	return
}

// Command returns the command named name.
//
// Return nil if the command does not exist.
func (c *Config) Command(name string) *Command {
	return c.commands[name]
}

// Commands returns all the commands.
func (c *Config) Commands() []*Command {
	cmds := make([]*Command, 0, len(c.commands))
	for _, cmd := range c.commands {
		cmds = append(cmds, cmd)
	}
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name() < cmds[j].Name() })
	return cmds
}

// ExecutedCommand returns the executed command and return nil if no command
// is executed.
func (c *Config) ExecutedCommand() *Command {
	return c.executed
}

// SetExecutedCommand sets the executed command.
func (c *Config) SetExecutedCommand(cmd *Command) *Config {
	c.panicIsParsed(false)
	c.executed = cmd
	c.Printf("Set the executed command '%s'", cmd.FullName())
	return c
}

//////////////////////////////////////////////////////////////////////////////
/// Group

func (c *Config) noticeNewGroup(group *OptGroup) {
	if _, ok := c.allGroups[group.FullName()]; ok {
		return
	}

	c.allGroups[group.FullName()] = group
	if group.cmd != nil {
		group.cmd.noticeNewGroup(group)
	}

	gnames := strings.Split(group.FullName(), c.gsep)
	if len(gnames) == 1 {
		return
	}

	for i, gname := range gnames {
		fullName := c.mergePaths(gnames[:i+1])
		if _, ok := c.allGroups[fullName]; !ok {
			group := newOptGroup2(false, c, group.cmd, gname, gnames[:i]...)
			c.allGroups[fullName] = group
		}
	}
}

func (c *Config) getGroup(parent, name string) *OptGroup {
	return c.allGroups[c.getFullGroupName(parent, name)]
}

// AllGroups returns all the groups containing the default group and
// all the sub-groups.
func (c *Config) AllGroups() []*OptGroup {
	groups := make([]*OptGroup, 0, len(c.allGroups))
	for _, group := range c.allGroups {
		if len(group.AllOpts()) > 0 {
			groups = append(groups, group)
		}
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].fname < groups[j].fname })
	return groups
}

// AllNotCommandGroups returns all the groups which don't belong to the command.
func (c *Config) AllNotCommandGroups() []*OptGroup {
	groups := make([]*OptGroup, 0, len(c.allGroups))
	for _, group := range c.allGroups {
		if group.cmd == nil && len(group.AllOpts()) > 0 {
			groups = append(groups, group)
		}
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].fname < groups[j].fname })
	return groups
}

// PrintTree prints the tree of the groups to os.Stdout.
//
// The command name is surrounded by "{" and "}" and the group name is surrounded
// by "[" and "]".
//
// Notice: it is only used to debug.
func (c *Config) PrintTree(w io.Writer) {
	indent := "|"

	// Print the default group.
	c.printGroup(w, c.OptGroup, indent)

	// Print all the commands.
	for _, cmd := range c.Commands() {
		c.printCommand(w, cmd, indent)
	}
}

func (c *Config) printGroup(w io.Writer, group *OptGroup, indent string) {
	// Print the name of the current group.
	if group == c.OptGroup {
		fmt.Fprintf(w, "[%s]\n", group.Name())
	} else {
		fmt.Fprintf(w, "%s-->[%s]\n", indent, group.Name())
		indent += "   |"
	}

	// Print the options in the current group.
	for _, opt := range group.CliOpts() {
		fmt.Fprintf(w, "%s--- %s*\n", indent, opt.Name())
	}
	for _, opt := range group.NotCliOpts() {
		fmt.Fprintf(w, "%s--- %s\n", indent, opt.Name())
	}

	// Print the sub-groups of the current group.
	for _, subGroup := range group.Groups() {
		c.printGroup(w, subGroup, indent)
	}
}

func (c *Config) printCommand(w io.Writer, cmd *Command, indent string) {
	// Print the name of the current command.
	fmt.Fprintf(w, "%s-->{%s}\n", indent, cmd.Name())
	indent += "   |"

	// Print the options in the default group of the current command.
	for _, opt := range cmd.CliOpts() {
		fmt.Fprintf(w, "%s--- %s*\n", indent, opt.Name())
	}
	for _, opt := range cmd.NotCliOpts() {
		fmt.Fprintf(w, "%s--- %s\n", indent, opt.Name())
	}

	// Print the sub-groups of the current command.
	for _, group := range cmd.Groups() {
		c.printGroup(w, group, indent)
	}

	// Print the sub-commands of the current command.
	for _, subCmd := range cmd.Commands() {
		c.printCommand(w, subCmd, indent)
	}
}
