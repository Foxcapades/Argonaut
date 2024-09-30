package argo

type UnmappedInputContainerSpec interface {
	UnmappedInputContainer

	UnmappedLabel() string

	HasUnmappedLabel() bool
}
