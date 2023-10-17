package argo

import (
	"fmt"
)

func uniqueCommandNames(branches []CommandBranchBuilder, leaves []CommandLeafBuilder, errs MultiError) {
	names := make(map[string]uint8, 32)
	uniqueNames(names, branches, leaves, errs)
}

func massUniqueCommandNames(groups []CommandGroupBuilder, errs MultiError) {
	names := make(map[string]uint8, 32)

	for _, group := range groups {
		uniqueNames(names, group.getBranches(), group.getLeaves(), errs)
	}
}

func uniqueNames(
	names map[string]uint8,
	branches []CommandBranchBuilder,
	leaves []CommandLeafBuilder,
	errs MultiError,
) {
	for _, branch := range branches {
		names[branch.GetName()]++
		if names[branch.GetName()] == 2 {
			errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", branch.GetName()))
		}

		for _, alias := range branch.GetAliases() {
			names[alias]++
			if names[alias] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
			}
		}
	}

	for _, leaf := range leaves {
		names[leaf.getName()]++
		if names[leaf.getName()] == 2 {
			errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", leaf.getName()))
		}

		for _, alias := range leaf.getAliases() {
			names[alias]++
			if names[alias] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
			}
		}
	}
}
