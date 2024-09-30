package argo

type CommandNode interface {
	Name() string

	FlagGroupContainer
}
