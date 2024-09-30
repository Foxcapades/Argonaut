package argo

type CommandSpec interface {
	// Name returns the name of the command.
	Name() string

	// Description returns the description text for the command.
	Description() string

	// HasDescription returns an indicator as to whether this CommandSpec has a
	// description value set.
	//
	// If this method returns `false`, the Description method will return an empty
	// string.
	HasDescription() bool

	FlagGroupSpecContainer

	ArgumentSpecContainer

	UnmappedInputContainerSpec

	PassthroughContainerSpec

	ToCommand() Command
}
