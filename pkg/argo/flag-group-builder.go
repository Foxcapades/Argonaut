package argo

type FlagGroupBuilder interface {
	GetName() string

	WithDescription(desc string) FlagGroupBuilder

	HasDescription() bool

	GetDescription() string

	WithFlag(flag FlagBuilder) FlagGroupBuilder

	HasFlags() bool

	GetFlags() []FlagBuilder

	Build() (FlagGroup, error)
}
