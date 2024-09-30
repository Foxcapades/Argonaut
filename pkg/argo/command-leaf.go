package argo

type CommandLeaf interface {
	CommandNode

	SubCommand

	ArgumentContainer

	UnmappedInputContainer

	PassthroughContainer
}
