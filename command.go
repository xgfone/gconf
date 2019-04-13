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

import "sort"

// Command represents the sub-command.
type Command struct {
	*OptGroup

	conf *Config
	help string

	parent    *Command
	commands  map[string]*Command
	allGroups map[string]*OptGroup
}

func newCommand(conf *Config, parent *Command, name, help string, parents ...string) *Command {
	if conf == nil {
		panic("Config must not be nil")
	} else if name == "" {
		panic("the group name must not be empty")
	}

	cmd := &Command{
		conf:   conf,
		help:   help,
		parent: parent,

		commands:  make(map[string]*Command, 8),
		allGroups: make(map[string]*OptGroup, 8),
	}
	cmd.OptGroup = newOptGroup(conf, cmd, name, parents...)

	conf.Printf("Creating the command '%s'", cmd.FullName())
	return cmd
}

// Config returns the Config that the current command belongs to.
func (cmd *Command) Config() *Config {
	return cmd.conf
}

//////////////////////////////////////////////////////////////////////////////
/// Command

// NewCommand returns a new sub-command named name with the document help.
//
// Notice: if the command has existed, it will return the old.
func (cmd *Command) NewCommand(name, help string) (c *Command) {
	if c = cmd.commands[name]; c == nil {
		c = newCommand(cmd.conf, cmd, name, help, cmd.OptGroup.paths...)
		cmd.commands[name] = c
	}
	return
}

// Command returns the sub-command named name.
//
// Return nil if the command does not exist.
func (cmd *Command) Command(name string) *Command {
	return cmd.commands[name]
}

// Commands returns all the sub-commands of the current command.
func (cmd *Command) Commands() []*Command {
	cmds := make([]*Command, 0, len(cmd.commands))
	for _, cmd := range cmd.commands {
		cmds = append(cmds, cmd)
	}
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name() < cmds[j].Name() })
	return cmds
}

// ParentCommand returns the parent command of the current command.
//
// Return nil if no parent command.
func (cmd *Command) ParentCommand() *Command {
	return cmd.parent
}

//////////////////////////////////////////////////////////////////////////////
/// Group

func (cmd *Command) noticeNewGroup(group *OptGroup) {
	if _, ok := cmd.allGroups[group.name]; !ok {
		cmd.allGroups[group.name] = group
	}
}

// AllGroups returns all the sub-groups containing the default group and sub-groups
// of the current command.
func (cmd *Command) AllGroups() []*OptGroup {
	groups := make([]*OptGroup, 0, len(cmd.allGroups))
	for _, group := range cmd.allGroups {
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].name < groups[j].name })
	return groups
}
