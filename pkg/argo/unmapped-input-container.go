package argo

type UnmappedInputContainer interface {
	UnmappedInputs() []string

	HasUnmappedInputs() bool
}
