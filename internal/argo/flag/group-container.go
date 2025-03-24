package flag

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type GroupContainer interface {
	FlagGroups() []argo.FlagGroup
	HasFlagGroups() bool
	FindShortFlag(c byte) argo.Flag
	FindLongFlag(name string) argo.Flag
}

type GroupContainerBuilder[T any] interface {
	FlagGroups(includeDefault bool) []argo.FlagGroupBuilder
	WithFlagGroup(group argo.FlagGroupBuilder) T
	HasFlagGroups(includeDefault bool) bool
}
