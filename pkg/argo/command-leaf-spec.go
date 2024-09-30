package argo

type CommandLeafSpec interface {
	CommandNodeSpec

	SubCommandSpec

	ArgumentSpecContainer

	UnmappedInputContainerSpec

	PassthroughContainerSpec
}
