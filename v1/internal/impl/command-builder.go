package impl

import "github.com/Foxcapades/Argonaut/v1/pkg/argo"

func NewCommandBuilder() *CommandBuilder {
	return new(CommandBuilder)
}

type CommandBuilder struct {
	desc string
}

func (c *CommandBuilder) Flag(argo.FlagBuilder) argo.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Unmarshaler(argo.InternalUnmarshaler) argo.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Description(desc string) argo.CommandBuilder {
	c.desc = desc
	return c
}

func (c *CommandBuilder) Arg(argo.ArgumentBuilder) argo.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Examples(...string) argo.CommandBuilder {
	panic("implement me")
}

func (c *CommandBuilder) Build() (argo.Command, error) {
	panic("implement me")
}

func (c *CommandBuilder) MustBuild() argo.Command {
	com, err := c.Build()
	if err != nil {
		panic(err)
	}
	return com
}

func (c *CommandBuilder) Parse() error {
	//com, err := c.Build()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (c *CommandBuilder) MustParse() {
	if err := c.Parse(); err != nil {
		panic(err)
	}
}
