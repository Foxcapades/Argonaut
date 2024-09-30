package argo

type FlagGroupSpecContainer interface {
	FlagGroups() []FlagGroupSpec

	FindFlagGroup(name string) FlagGroupSpec

	Flags() []FlagSpec

	FindFlagByShortForm(name byte) FlagSpec

	FindFlagByLongForm(name string) FlagSpec
}
