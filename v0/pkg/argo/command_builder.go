package argo

type CommandBuilder interface {
	Description(string) (this CommandBuilder)

	GetDescription() string

	HasDescription() bool

	Examples(...string) (this CommandBuilder)

	Arg(ArgumentBuilder) (this CommandBuilder)

	GetArgs() []ArgumentBuilder

	Flag(FlagBuilder) (this CommandBuilder)

	FlagGroup(builder FlagGroupBuilder) (this CommandBuilder)

	GetFlagGroups() []FlagGroupBuilder

	Unmarshaler(ValueUnmarshaler) (this CommandBuilder)

	DisableHelp() (this CommandBuilder)

	Build() (Command, error)

	MustBuild() Command

	Parse() (com Command, err error)

	MustParse() Command

	Warnings() []string
}
