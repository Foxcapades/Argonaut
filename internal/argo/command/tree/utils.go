package tree

import (
	"fmt"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// EnsureDefaultCommandGroup tests the given slice of argo.CommandGroupBuilder
// instances to ensure that the first element in the slice is a "default"
// command group instance.
//
// If the first element in the slice is _not_ a default command group, a new
// slice will be returned with element 0 being a newly constructed default
// command group instance.
//
// If the first element in the slice _is_ a default command group, the input
// slice will be returned.
//
// This method is intended to be shared by argo.ParentNodeBuilder instances to
// verify that the first command group is a default command group before
// appending ungrouped argo.ChildNodeBuilder instances to it.
func EnsureDefaultCommandGroup(groups []argo.CommandGroupBuilder) []argo.CommandGroupBuilder {
	if len(groups) > 0 && groups[0].Name() != DefaultCommandGroupName {
		return append(
			append(
				make([]argo.CommandGroupBuilder, 0, len(groups)+1),
				NewGroupBuilder(DefaultCommandGroupName),
			),
			groups...,
		)
	}

	return groups
}

// UniqueCommandNames checks the names and aliases of all the
// argo.ChildNodeBuilder instances contained by all the argo.CommandGroupBuilder
// instances in the given slice to ensure there are no conflicts which would
// create ambiguity in CLI calls.
//
// For example, if a command in group 0 is named "foo", and a command in group 1
// has the alias "foo", an error will be generated for that conflict.
func UniqueCommandNames(groups []argo.CommandGroupBuilder, errs argo.MultiError) {
	names := make(map[string]uint8, 32)

	for _, group := range groups {
		for node := range group.Subcommands() {
			names[node.Name()]++
			if names[node.Name()] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", node.Name()))
			}

			for _, alias := range node.Aliases() {
				names[alias]++
				if names[alias] == 2 {
					errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
				}
			}
		}
	}
}

// MakeDefaultOnIncompleteHandler returns a default implementation of the
// argo.IncompleteCommandHandler function type which renders the help text of
// the last reached argo.ParentNode instance and exits the application with exit
// code 1.
func MakeDefaultOnIncompleteHandler[T argo.ParentNode](opts argo.Options) argo.IncompleteCommandHandler[T] {
	return func(parent T) {
		if tree, ok := argo.ParentNode(parent).(argo.TreeCommand); ok {
			utils.Must(RenderHelp(tree, opts, os.Stderr))
		} else if branch, ok := argo.ParentNode(parent).(argo.BranchCommand); ok {
			utils.Must(RenderBranchHelp(branch, opts, os.Stderr))
		} else {
			panic("illegal state: unrecognized command parent implementation")
		}

		os.Exit(1)
	}
}

func BuildCommandGroups(
	builders []argo.CommandGroupBuilder,
	opts argo.Options,
	parent argo.ParentNode,
	errs argo.MultiError,
) []argo.CommandGroup {
	commandGroups := make([]argo.CommandGroup, 0, len(builders))

	UniqueCommandNames(builders, errs)

	for _, build := range builders {
		if build.HasSubcommands() {
			if group, err := BuildGroup(build, opts, parent); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	return commandGroups
}
