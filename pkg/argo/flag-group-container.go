package argo

// FlagGroupContainer defines the functionality common to all CLI component
// types that may have FlagGroup instances attached.
type FlagGroupContainer interface {
	// FlagGroups returns a slice containing all the FlagGroup instances attached
	// to this CLI component.
	FlagGroups() []FlagGroup

	// FindFlagGroup attempts to find a FlagGroup instance with a name matching
	// the given input.
	FindFlagGroup(name string) FlagGroup

	Flags() []Flag

	// FindFlagByShortForm attempts to find a Flag instance with a short-form name
	// matching the given input.
	FindFlagByShortForm(name byte) Flag

	// FindFlagByLongForm attempts to find a Flag instance with a long-form name
	// matching the given input.
	FindFlagByLongForm(name string) Flag
}
