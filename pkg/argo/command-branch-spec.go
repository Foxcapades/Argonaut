package argo

type CommandBranchSpec interface {
	Name() string

	FlagGroupSpecContainer
}
