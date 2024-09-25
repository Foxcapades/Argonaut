package argo

type Flag interface {
	IsRequired() bool

	HasExplicitArgument() bool

	Argument() Argument

	WasUsed() bool

	UsageCount() int
}
