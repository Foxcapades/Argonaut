package argo

type commandBuilder struct {

}

func (c *commandBuilder) Description(string) CommandBuilder {
	panic("implement me")
}

func (c *commandBuilder) Arg(ArgumentBuilder) CommandBuilder {
	panic("implement me")
}

func (c *commandBuilder) Examples(...string) CommandBuilder {
	panic("implement me")
}

func (c *commandBuilder) Build() (Command, error) {
	panic("implement me")
}

func (c *commandBuilder) MustBuild() Command {
	com, err := c.Build()
	if err != nil {
		panic(err)
	}
	return com
}

func (c *commandBuilder) Parse() error {
	//com, err := c.Build()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (c *commandBuilder) MustParse() {
	if err := c.Parse(); err != nil {
		panic(err)
	}
}

