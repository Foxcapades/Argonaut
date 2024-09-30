package argo

type PassthroughContainer interface {
	PassthroughInputs() []string

	HasPassthroughInputs() bool
}
