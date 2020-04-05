package argo

type CommandBuilder interface {
	Description(string) (this CommandBuilder)

	Examples(...string) (this CommandBuilder)

	Arg(ArgumentBuilder) (this CommandBuilder)

	Flag(FlagBuilder) (this CommandBuilder)

	FlagGroup(builder FlagGroupBuilder) (this CommandBuilder)

	Unmarshaler(ValueUnmarshaler) (this CommandBuilder)

	Build() (Command, error)

	MustBuild() Command

	Parse() (extra []string, err error)

	MustParse() []string
}
