package tree

import (
	"fmt"
	"os"

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

func defaultOnIncompleteHandler(parent argo.ParentNode) {
	if tree, ok := parent.(argo.CommandTree); ok {
		utils.Must(comTreeRenderer{}.RenderHelp(tree, os.Stderr))
	} else if branch, ok := parent.(argo.BranchCommand); ok {
		utils.Must(comBranchRenderer{}.RenderHelp(branch, os.Stderr))
	} else {
		panic("illegal state: unrecognized command parent implementation")
	}

	os.Exit(1)
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
