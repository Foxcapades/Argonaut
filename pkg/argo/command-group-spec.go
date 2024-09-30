package argo

type CommandGroupSpec interface {
	Name() string

	Description() string

	HasDescription() bool

	SubCommandSpecContainer
}
