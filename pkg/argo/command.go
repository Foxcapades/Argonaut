package argo

type Command interface {
	// Name returns the name of the command.
	Name() string

	FlagGroupContainer

	ArgumentContainer

	UnmappedInputContainer

	PassthroughContainer
}
