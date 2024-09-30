package argo

type CommandGroup interface {
	Name() string

	SubCommandContainer
}
