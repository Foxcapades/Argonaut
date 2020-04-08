package com

import (
	"errors"
	"fmt"
	"os"

	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"

	"github.com/Foxcapades/Argonaut/v0/internal/util"

	"github.com/Foxcapades/Argonaut/v0/internal/impl/marsh"
	"github.com/Foxcapades/Argonaut/v0/internal/impl/parse"
)

func NewBuilder(provider A.Provider) A.CommandBuilder {
	return &CommandBuilder{
		provider:    provider,
		fGroups:     []A.FlagGroupBuilder{provider.NewFlagGroup()},
		parser:      parse.NewParser(),
		unmarshaler: marsh.NewDefaultedValueUnmarshaler(),
	}
}

type CommandBuilder struct {
	provider    A.Provider
	name        string
	desc        string
	parser      A.Parser
	unmarshaler A.ValueUnmarshaler
	fGroups     []A.FlagGroupBuilder
	args        []A.ArgumentBuilder
	warnings    []string
	examples    []string
	omitHelp    bool
}

func (c *CommandBuilder) DisableHelp() A.CommandBuilder {
	c.omitHelp = true
	return c
}

func (c *CommandBuilder) Flag(flag A.FlagBuilder) A.CommandBuilder {
	if flag == nil {
		return c.warn("nil value passed to Flag()")
	}
	c.fGroups[0].Flag(flag)
	return c
}

func (c *CommandBuilder) FlagGroup(builder A.FlagGroupBuilder) (this A.CommandBuilder) {
	if builder == nil {
		return c.warn("nil value passed to FlagGroup()")
	}
	c.fGroups = append(c.fGroups, builder)
	return c
}

func (c *CommandBuilder) GetFlagGroups() []A.FlagGroupBuilder {
	return c.fGroups
}

func (c *CommandBuilder) Unmarshaler(un A.ValueUnmarshaler) A.CommandBuilder {
	if un == nil {
		return c.warn("nil value passed to Unmarshaler()")
	}
	c.unmarshaler = un
	return c
}

func (c *CommandBuilder) Description(desc string) A.CommandBuilder {
	c.desc = desc
	return c
}

func (c *CommandBuilder) GetDescription() string {
	return c.desc
}

func (c *CommandBuilder) HasDescription() bool {
	return len(c.desc) > 0
}

func (c *CommandBuilder) Arg(arg A.ArgumentBuilder) A.CommandBuilder {
	if arg == nil {
		return c.warn("nil value passed to Arg()")
	}
	c.args = append(c.args, arg)
	return c
}

func (c *CommandBuilder) GetArgs() []A.ArgumentBuilder {
	return c.args
}

func (c *CommandBuilder) Examples(examples ...string) A.CommandBuilder {
	c.examples = examples
	return c
}

func (c *CommandBuilder) Build() (A.Command, error) {
	short := make(map[byte]bool)
	long := make(map[string]bool)
	out := new(Command)

	// Check for overlapping flags
	for _, g := range c.fGroups {
		for _, f := range g.GetFlags() {
			if f.HasShort() {
				if short[f.GetShort()] {
					// TODO: make this a real error
					return nil, errors.New("more than one flag with the identifier -" + string(f.GetShort()))
				}
				short[f.GetShort()] = true
			}
			if f.HasLong() {
				if long[f.GetLong()] {
					// TODO: make this a real error
					return nil, errors.New("more than one flag with the identifier --" + f.GetLong())
				}
				long[f.GetLong()] = true
			}
		}
	}

	if !c.omitHelp {
		c.fGroups = append(c.fGroups, c.provider.
			NewFlagGroup().
			Name("Help & Info").
			Flag(c.provider.
				NewFlag().
				Short('h').
				Long("help").
				Description("Prints this help text").
				OnHit(func(f A.Flag) {
					fmt.Println(util.RenderHelp(f))
					os.Exit(0)
				})))
	}

	// Build groups
	groups := make([]A.FlagGroup, len(c.fGroups))

	for i, fg := range c.fGroups {
		fg.Parent(out)
		if g, err := fg.Build(); err != nil {
			return nil, err
		} else {
			groups[i] = g
		}
	}

	args := make([]A.Argument, len(c.args))
	for i, arg := range c.args {
		arg.Parent(out)
		if a, err := arg.Build(); err != nil {
			return nil, err
		} else {
			args[i] = a
		}
	}

	out.description = c.desc
	out.groups = groups
	out.arguments = args
	out.unmarshal = c.unmarshaler

	return out, nil
}

func (c *CommandBuilder) MustBuild() A.Command {
	com, err := c.Build()
	if err != nil {
		panic(err)
	}
	return com
}

func (c *CommandBuilder) Parse() (extra []string, err error) {
	com, err := c.Build()
	if err != nil {
		return nil, err
	}

	err = c.parser.Parse(os.Args, com)

	return c.parser.Passthroughs(), err
}

func (c *CommandBuilder) MustParse() []string {
	if a, b := c.Parse(); b != nil {
		panic(b)
	} else {
		return a
	}
}

func (c *CommandBuilder) Warnings() []string {
	return c.warnings
}

func (c *CommandBuilder) warn(txt string) A.CommandBuilder {
	c.warnings = append(c.warnings, "CommandBuilder: "+txt)
	return c
}
