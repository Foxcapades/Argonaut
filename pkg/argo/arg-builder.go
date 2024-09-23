package argo

type ArgumentSpecBuilder interface {
	Build() (Argument, error)
}
