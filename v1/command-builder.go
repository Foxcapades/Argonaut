package argo

type CommandBuilder interface {
	Description(string) CommandBuilder

	Arg(ArgumentBuilder) CommandBuilder

	Examples(...string) CommandBuilder

	Flag(FlagBuilder) CommandBuilder

	Unmarshaler(InternalUnmarshaler) CommandBuilder

	Build() (Command, error)

	MustBuild() Command

	Parse() error

	MustParse()
}
