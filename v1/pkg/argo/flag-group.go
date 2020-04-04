package argo

type FlagGroupBuilder interface {
	Name(name string) (this FlagGroupBuilder)

	Description(desc string) (this FlagGroupBuilder)

	Flag(flag FlagBuilder) (this FlagGroupBuilder)

	Build() (FlagGroup, error)

	MustBuild() FlagGroup
}

type FlagGroup interface {
	Name() string
	Description() string
	Flags() []Flag
}
