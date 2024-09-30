package argo

type ArgumentSpecContainer interface {
	Arguments() []ArgumentSpec

	HasArguments() bool
}
