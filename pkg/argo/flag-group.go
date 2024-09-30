package argo

// FlagGroup represents an organizational, named grouping of CLI flags.
type FlagGroup interface {
	// Name returns the name of the FlagGroupSpec.
	Name() string

	// Flags returns a slice containing all the FlagSpec instances contained by this
	// FlagGroupSpec.
	Flags() []Flag

	// FindFlagByShortForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a short-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByShortForm(name byte) Flag

	// FindFlagByLongForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a long-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByLongForm(name string) Flag
}
