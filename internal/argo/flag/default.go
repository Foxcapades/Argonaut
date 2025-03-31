package flag

import "github.com/foxcapades/argonaut/v3/pkg/argo"

const DefaultFlagGroupName = "__DEFAULT_FLAG_GROUP__"

func EnsureDefaultGroup(groups []argo.FlagGroupBuilder) []argo.FlagGroupBuilder {
	if len(groups) == 0 {
		return append(make([]argo.FlagGroupBuilder, 0, 2), NewDefaultGroupBuilder())
	}

	if !IsDefaultGroup(groups[0]) {
		return append(append(make([]argo.FlagGroupBuilder, 0, len(groups)+1), NewDefaultGroupBuilder()), groups...)
	}

	return groups
}
