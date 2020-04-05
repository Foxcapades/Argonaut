package impl

import A "github.com/Foxcapades/Argonaut/v1/pkg/argo"

func NewCommandBuilder() A.CommandBuilder {
	return &CommandBuilder{
		fGroups: []A.FlagGroupBuilder{NewFlagGroupBuilder()},
		parser:  NewParser(),
	}
}

type CommandBuilder struct {
	name    string
	desc    string
	fGroups []A.FlagGroupBuilder
	parser  A.Parser
}

func (c *CommandBuilder) Flag(flag A.FlagBuilder) A.CommandBuilder {
	c.fGroups[0].Flag(flag)
	return c
}

func (c *CommandBuilder) FlagGroup(builder A.FlagGroupBuilder) (this A.CommandBuilder) {
	c.fGroups = append(c.fGroups, builder)
	return c
}

func (c *CommandBuilder) Unmarshaler(A.ValueUnmarshaler) A.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Description(desc string) A.CommandBuilder {
	c.desc = desc
	return c
}

func (c *CommandBuilder) Arg(A.ArgumentBuilder) A.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Examples(...string) A.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Build() (A.Command, error) {
	panic("implement me")
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
