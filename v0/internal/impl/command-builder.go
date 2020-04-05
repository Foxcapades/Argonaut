package impl

import (
	"errors"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewCommandBuilder() A.CommandBuilder {
	return &CommandBuilder{
		fGroups:     []A.FlagGroupBuilder{NewFlagGroupBuilder()},
		parser:      NewParser(),
		unmarshaler: NewDefaultedValueUnmarshaler(),
	}
}

type CommandBuilder struct {
	name        string
	desc        string
	parser      A.Parser
	unmarshaler A.ValueUnmarshaler
	fGroups     []A.FlagGroupBuilder
	args        []A.ArgumentBuilder
	warnings    []string
	examples    []string
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

func (c *CommandBuilder) Arg(arg A.ArgumentBuilder) A.CommandBuilder {
	if arg == nil {
		return c.warn("nil value passed to Arg()")
	}
	c.args = append(c.args, arg)
	return c
}

func (c *CommandBuilder) Examples(examples ...string) A.CommandBuilder {
	c.examples = examples
	return c
}

func (c *CommandBuilder) Build() (A.Command, error) {
	short := make(map[byte]bool)
	long := make(map[string]bool)

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

	// Build groups
	groups := make([]A.FlagGroup, len(c.fGroups))

	for i, fg := range c.fGroups {
		if g, err := fg.Build(); err != nil {
			return nil, err
		} else {
			groups[i] = g
		}
	}


	return &Command{
		description: c.desc,
		groups:      groups,
		unmarshal:   c.unmarshaler,
	}, nil
}

func (c *CommandBuilder) MustBuild() A.Command {
	com, err := c.Build()
	if err != nil {
		panic(err)
	}
	return com
}

func (c *CommandBuilder) Parse() (extra []string, err error) {
	panic("implement me")
}

func (c *CommandBuilder) MustParse() []string {
	panic("implement me")
}

func (c *CommandBuilder) warn(txt string) A.CommandBuilder {
	c.warnings = append(c.warnings, "CommandBuilder: "+txt)
	return c
}
