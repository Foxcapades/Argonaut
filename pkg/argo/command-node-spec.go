package argo

type CommandNodeSpec interface {
	Name() string

	FlagGroupSpecContainer
}
