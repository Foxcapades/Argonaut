package argo

type FlagGroupSpec interface {

	// Name returns the name of the FlagGroupSpec.
	Name() string

	// Description returns the description text for this FlagGroupSpec.
	Description() string

	// HasDescription returns an indicator as to whether this FlagGroupSpec has
	// description text set.
	HasDescription() bool

	// Flags returns a slice containing all the FlagSpec instances contained by this
	// FlagGroupSpec.
	Flags() []FlagSpec

	// FindFlagByShortForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a short-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByShortForm(name byte) FlagSpec

	// FindFlagByLongForm searches the FlagSpec instances in this FlagGroupSpec for one
	// that has a long-form matching the given name.
	//
	// If no such FlagSpec instance can be found, this method will return `nil`.
	FindFlagByLongForm(name string) FlagSpec

	ToFlagGroup() FlagGroup
}
