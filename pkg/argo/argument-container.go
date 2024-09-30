package argo

type ArgumentContainer interface {
	Arguments() []ArgumentSpec

	HasArguments() bool
}
