package comutil

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func UniqueNames(branches []argo.CommandBranchBuilder, leaves []argo.CommandLeafBuilder, errs argo.MultiError) {
	names := make(map[string]uint8, 32)
	uniqueNames(names, branches, leaves, errs)
}

func MassUniqueNames(groups []argo.CommandGroupBuilder, errs argo.MultiError) {
	names := make(map[string]uint8, 32)

	for _, group := range groups {
		uniqueNames(names, group.GetBranches(), group.GetLeaves(), errs)
	}
}

func uniqueNames(
	names map[string]uint8,
	branches []argo.CommandBranchBuilder,
	leaves []argo.CommandLeafBuilder,
	errs argo.MultiError,
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
		names[leaf.GetName()]++
		if names[leaf.GetName()] == 2 {
			errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", leaf.GetName()))
		}

		for _, alias := range leaf.GetAliases() {
			names[alias]++
			if names[alias] == 2 {
				errs.AppendError(fmt.Errorf("conflicting subcommand name/alias: %s", alias))
			}
		}
	}
}
