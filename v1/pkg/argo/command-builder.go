package argo

type CommandBuilder interface {
	Description(string) CommandBuilder

	Examples(...string) CommandBuilder

	Arg(ArgumentBuilder) CommandBuilder

	Flag(FlagBuilder) CommandBuilder

	Unmarshaler(InternalUnmarshaler) CommandBuilder

	Build() (Command, error)

	MustBuild() Command

	Parse() error

	MustParse()
}
