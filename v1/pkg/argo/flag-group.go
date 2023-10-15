package argo

// A FlagGroup is a named group or category for a collection of one or more
// Flag instances.
//
// Flag groups are primarily used to categorize CLI flags when rendering help
// text.
//
// Flag groups that do not contain any flags will be filtered out at build time.
type FlagGroup interface {

	// Name of the FlagGroup.
	//
	// This value will be used when rendering help text.
	Name() string

	// Description for this FlagGroup.
	//
	// This optional value will be used when rendering help text.
	Description() string

	// HasDescription indicates whether this FlagGroup has a description attached.
	HasDescription() bool

	// Flags returns the Flag instances contained by this FlagGroup.
	Flags() []Flag

	// FindShortFlag checks this FlagGroup's contents for a flag with the given
	// short flag character.  If one could not be found, this method returns nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag checks this FlagGroup's contents for a flag with the given
	// long flag name.  If one could not be found, this method returns nil.
	FindLongFlag(name string) Flag

	// TryFlag searches this FlagGroup's contents for a Flag instance matching the
	// given FlagRef, then increments that flag's hit counter, additionally
	// attempting to set the Flag's Argument value if one is present.
	TryFlag(ref FlagRef) (bool, error)
}
