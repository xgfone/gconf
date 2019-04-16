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

	"github.com/urfave/cli"
)

type cliParser struct {
	utoh bool
	app  *cli.App
	pre  func(*Config, *cli.App) error
	post func(*Config, *cli.App) error
}

// NewDefaultCliParser is equal to NewCliParser(nil, underlineToHyphen[0]).
func NewDefaultCliParser(underlineToHyphen ...bool) Parser {
	var u2h bool
	if len(underlineToHyphen) > 0 {
		u2h = underlineToHyphen[0]
	}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	return NewCliParser(app, u2h, nil, nil)
}

// NewCliParser returns a new cli parser based on "github.com/urfave/cli".
//
// Notice: You should set the action for each command and the app main,
// because "cli" does not only return an error but also call os.Exit()
// for the help and the version option. So there is no way to distinguish it,
// that's, you should not run other logic codes after calling Parse().
//
// If something is wrong, you maybe open the debug mode of Config.
func NewCliParser(app *cli.App, underlineToHyphen bool, pre, post func(*Config, *cli.App) error) Parser {
	if app == nil {
		app = cli.NewApp()
	}
	if pre == nil {
		pre = func(*Config, *cli.App) error { return nil }
	}
	if post == nil {
		post = func(*Config, *cli.App) error { return nil }
	}
	return cliParser{app: app, utoh: underlineToHyphen, pre: pre, post: post}
}

func (cp cliParser) Name() string {
	return "cli"
}

func (cp cliParser) Priority() int {
	return 0
}

func (cp cliParser) Pre(conf *Config) error {
	cp.app.Name = conf.Name()

	if help := conf.Description(); help != "" {
		cp.app.Usage = help
	}

	if _, _, version, _ := conf.GetCliVersion(); version != "" {
		cp.app.Version = version
	}

	return cp.pre(conf, cp.app)
}

func (cp cliParser) Post(conf *Config) error {
	return cp.post(conf, cp.app)
}

func (cp cliParser) updateConfigOpt(names []string, ctx *cli.Context,
	conf *Config, flag2opts map[string]*groupOpt) (err error) {

	for _, name := range names {
		gopt := flag2opts[name]
		if gopt.Ok {
			continue
		}

		var global bool
		if ctx.GlobalGeneric(name) != nil {
			global = true
		} else if ctx.Generic(name) == nil {
			conf.Printf("[%s] WARNING: the context '%s' has no value of the flag '%s'",
				cp.Name(), ctx.Command.FullName(), name)
			continue
		}

		var value interface{}
		switch gopt.Flag.(type) {
		case cli.BoolFlag:
			if global {
				value = ctx.GlobalBool(name)
			} else {
				value = ctx.Bool(name)
			}
		case cli.BoolTFlag:
			if global {
				value = ctx.GlobalBoolT(name)
			} else {
				value = ctx.BoolT(name)
			}
		case cli.Int64Flag:
			if global {
				value = ctx.GlobalInt64(name)
			} else {
				value = ctx.Int64(name)
			}
		case cli.Uint64Flag:
			if global {
				value = ctx.GlobalUint64(name)
			} else {
				value = ctx.Uint64(name)
			}
		case cli.Float64Flag:
			if global {
				value = ctx.GlobalFloat64(name)
			} else {
				value = ctx.Float64(name)
			}
		case cli.DurationFlag:
			if global {
				value = ctx.GlobalDuration(name)
			} else {
				value = ctx.Duration(name)
			}
		case cli.StringFlag:
			if global {
				value = ctx.GlobalString(name)
			} else {
				value = ctx.String(name)
			}
		}

		if value != nil {
			if err = gopt.Group.SetOptValue(cp.Priority(), gopt.Opt.Name(), value); err != nil {
				return err
			}
			gopt.Ok = true
		}
	}

	return nil
}

func (cp cliParser) updateConfig(ctx *cli.Context, conf *Config,
	flag2opts map[string]*groupOpt) (err error) {
	origCtx := ctx

	// For the current command
	if err = cp.updateConfigOpt(ctx.FlagNames(), ctx, conf, flag2opts); err != nil {
		return
	}

	// For the parent command
	for ctx.Parent() != nil {
		if err = cp.updateConfigOpt(ctx.GlobalFlagNames(), ctx, conf, flag2opts); err != nil {
			return
		}
		ctx = ctx.Parent()
	}

	// For the global, that's non-command.
	if err = cp.updateConfigOpt(ctx.GlobalFlagNames(), ctx, conf, flag2opts); err != nil {
		return
	}

	if args := origCtx.Args(); len(args) > 0 {
		conf.SetCliArgs([]string(args))
	}
	return nil
}

func (cp cliParser) getAppFlags(groups []*OptGroup, flag2opts map[string]*groupOpt) (flags []cli.Flag) {
	for _, group := range groups {
		conf := group.Config()
		gname := group.OnlyGroupName()
		isDefault := group.IsConfigDefaultGroup()

		var cmdStr string
		if cmd := group.Command(); cmd != nil {
			cmdStr = fmt.Sprintf(" for the command '%s'", cmd.FullName())
		}

		for _, opt := range group.CliOpts() {
			// Get the name of the option.
			name := opt.Name()
			if name == "" {
				panic("the option name must not be empty")
			} else if !isDefault && gname != "" {
				name = fmt.Sprintf("%s%s%s", gname, conf.GetGroupSeparator(), name)
			}
			if cp.utoh {
				name = strings.Replace(name, "_", "-", -1)
			}
			if short := opt.Short(); short != "" {
				name = name + ", " + short
			}

			// Get the default value of the option
			help := opt.Help()

			var flag cli.Flag
			switch opt.Zero().(type) {
			case bool:
				if v := opt.Default(); v != nil && v.(bool) {
					flag = cli.BoolTFlag{Name: name, Usage: help}
				} else {
					flag = cli.BoolFlag{Name: name, Usage: help}
				}
				conf.Printf("[%s] Add the bool flag '%s'%s", cp.Name(), name, cmdStr)
			case int, int8, int16, int32, int64:
				v, _ := ToInt64(opt.Default())
				flag = cli.Int64Flag{Name: name, Usage: help, Value: v}
				conf.Printf("[%s] Add the int flag '%s'%s", cp.Name(), name, cmdStr)
			case uint, uint8, uint16, uint32, uint64:
				v, _ := ToUint64(opt.Default())
				flag = cli.Uint64Flag{Name: name, Usage: help, Value: v}
				conf.Printf("[%s] Add the uint flag '%s'%s", cp.Name(), name, cmdStr)
			case float32, float64:
				v, _ := ToFloat64(opt.Default())
				flag = cli.Float64Flag{Name: name, Usage: help, Value: v}
				conf.Printf("[%s] Add the float flag '%s'%s", cp.Name(), name, cmdStr)
			case time.Duration:
				v, _ := ToDuration(opt.Default())
				flag = cli.DurationFlag{Name: name, Usage: help, Value: v}
				conf.Printf("[%s] Add the time.Duration flag '%s'%s", cp.Name(), name, cmdStr)
			default: // Default for string
				v, _ := ToString(opt.Default())
				flag = cli.StringFlag{Name: name, Usage: help, Value: v}
				conf.Printf("[%s] Add the string flag '%s'%s", cp.Name(), name, cmdStr)
			}

			flags = append(flags, flag)
			// name = group.FullName() + conf.GetGroupSeparator() + opt.Name()
			flag2opts[name] = &groupOpt{Flag: flag, Group: group, Opt: opt}
		}
	}

	return
}

func (cp cliParser) getCmdAction(cmd *Command, flag2opts map[string]*groupOpt) func(*cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		conf := cmd.Config()
		conf.SetExecutedCommand(cmd)
		if err = cp.handleConfigOption(ctx, conf, flag2opts); err == nil {
			if action := cmd.Action(); action != nil {
				conf.Printf("[%s] Calling the action of the command '%s'",
					cp.Name(), cmd.FullName())
				err = action()
			} else {
				conf.Printf("[%s] WARNING: no action of the command '%s'",
					cp.Name(), cmd.FullName())
			}
		}
		return
	}
}

func (cp cliParser) getAppCommands(cmds []*Command, flag2opts map[string]*groupOpt) (commands []cli.Command) {
	for _, cmd := range cmds {
		commands = append(commands, cli.Command{
			Name:    cmd.Name(),
			Usage:   cmd.Description(),
			Aliases: cmd.Aliases(),
			Action:  cp.getCmdAction(cmd, flag2opts),

			Flags:       cp.getAppFlags(cmd.AllGroups(), flag2opts),
			Subcommands: cp.getAppCommands(cmd.Commands(), flag2opts),
		})
	}
	return
}

type groupOpt struct {
	Ok    bool
	Group *OptGroup
	Flag  cli.Flag
	Opt   Opt
}

func (cp cliParser) Parse(conf *Config) (err error) {
	flag2opts := make(map[string]*groupOpt, 8)
	action := conf.Action()
	if action == nil {
		conf.Printf("[%s] WARNING: Config is short of Action", cp.Name())
	}

	cp.app.Flags = cp.getAppFlags(conf.AllNotCommandGroups(), flag2opts)
	cp.app.Commands = cp.getAppCommands(conf.Commands(), flag2opts)
	cp.app.Action = func(ctx *cli.Context) (err error) {
		if err = cp.handleConfigOption(ctx, conf, flag2opts); err == nil {
			if action != nil {
				err = action()
			}
		}
		return
	}

	conf.Stop() // Stop the subsequent parsing
	return cp.app.Run(append([]string{conf.Name()}, conf.ParsedCliArgs()...))
}

func (cp cliParser) handleConfigOption(ctx *cli.Context, conf *Config,
	flag2opts map[string]*groupOpt) (err error) {
	if err = cp.updateConfig(ctx, conf, flag2opts); err == nil {
		// Call other parsers.
		for _, parser := range conf.Parsers() {
			if parser.Name() == cp.Name() {
				continue
			}
			conf.Printf("[%s] Calling the parser '%s'", cp.Name(), parser.Name())
			if err = parser.Parse(conf); err != nil {
				return
			}
		}

		// Clean all the parsers.
		for _, parser := range conf.Parsers() {
			conf.Printf("[%s] Cleaning the parser '%s'", cp.Name(), parser.Name())
			if err = parser.Post(conf); err != nil {
				return
			}
		}

		// Check and validate the options.
		return conf.CheckRequiredOption()
	}
	return nil
}
