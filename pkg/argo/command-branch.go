package argo

type CommandBranch interface {
	Name() string

	FlagGroupContainer
}
