package argo

type FlagGroupBuilder interface {
	Name(name string) (this FlagGroupBuilder)

	GetName() string

	Description(desc string) (this FlagGroupBuilder)

	GetDescription() string

	Flag(flag FlagBuilder) (this FlagGroupBuilder)

	GetFlags() []FlagBuilder

	Build() (FlagGroup, error)

	MustBuild() FlagGroup

	Parent(c Command) FlagGroupBuilder
}

type FlagGroup interface {
	Name() string

	HasName() bool

	Description() string

	HasDescription() bool

	Flags() []Flag

	HasFlags() bool

	Parent() Command
}
