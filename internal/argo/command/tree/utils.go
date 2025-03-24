package tree

import (
	"fmt"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func EnsureDefaultCommandGroup(groups []argo.CommandGroupBuilder) []argo.CommandGroupBuilder {
	if len(groups) > 0 && groups[0].Name() != DefaultCommandGroupName {
		return append(append(make([]argo.CommandGroupBuilder, 0, len(groups)+1), NewGroupBuilder(DefaultCommandGroupName)), groups...)
	}

	return groups
}

func uniqueCommandNames(branches []argo.BranchCommandBuilder, leaves []argo.LeafCommandBuilder, errs argo.MultiError) {
	names := make(map[string]uint8, 32)
	uniqueNames(names, branches, leaves, errs)
}

func massUniqueCommandNames(groups []argo.CommandGroupBuilder, errs argo.MultiError) {
	names := make(map[string]uint8, 32)

	for _, group := range groups {
		uniqueNames(names, group.Branches(), group.Leaves(), errs)
	}
}

func makeDefaultOnIncompleteHandler[T argo.ParentNode](opts argo.Options) argo.IncompleteCommandHandler[T] {
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

func uniqueNames(
	names map[string]uint8,
	branches []argo.BranchCommandBuilder,
	leaves []argo.LeafCommandBuilder,
	errs argo.MultiError,
) {
	for _, branch := range branches {
		names[branch.Name()]++
		if names[branch.Name()] == 2 {
			errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", branch.Name()))
		}

		for _, alias := range branch.Aliases() {
			names[alias]++
			if names[alias] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
			}
		}
	}

	for _, leaf := range leaves {
		names[leaf.Name()]++
		if names[leaf.Name()] == 2 {
			errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", leaf.Name()))
		}

		for _, alias := range leaf.Aliases() {
			names[alias]++
			if names[alias] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
			}
		}
	}
}

func processFlagGroups(builders []argo.FlagGroupBuilder, errs argo.MultiError) []argo.FlagGroup {
	flagGroups := make([]argo.FlagGroup, 0, len(builders))

	flag.UniqueFlagNames(builders, errs)

	for _, builder := range builders {
		if builder.HasFlags() {
			if group, err := flag.BuildGroup(builder); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	return flagGroups
}

func processCommandGroups(builders []argo.CommandGroupBuilder, opts argo.Options, parent argo.ParentNode, errs argo.MultiError) []argo.CommandGroup {
	commandGroups := make([]argo.CommandGroup, 0, len(builders))

	massUniqueCommandNames(builders, errs)

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
