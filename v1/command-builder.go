package argo

type CommandBuilder interface {
	Description(string) CommandBuilder

	Arg(ArgumentBuilder) CommandBuilder

	Examples(...string) CommandBuilder

	Build() (Command, error)

	MustBuild() Command

	Unmarshaler(InternalUnmarshaler) CommandBuilder

	Parse() error

	MustParse()
}
