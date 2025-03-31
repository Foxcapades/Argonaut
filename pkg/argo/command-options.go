package argo

type CommandOptions struct {

	// MaxDefaultFlagGroupSizeForMetaGroup defines the maximum number of flags
	// that a command's default flag group may have before automatically generated
	// flags such as help flags are split off into their own group.
	//
	// Setting this to a value less than 1 disables the meta group generation.
	MaxDefaultFlagGroupSizeForMetaGroup int

	// MetaFlagGroupName defines the name of the group that will be generated for
	// meta flags if/when the number of flags in the default flag group exceeds
	// the value of MaxDefaultFlagGroupSizeForMetaGroup.
	MetaFlagGroupName string

	// HelpTextMaxWidth defines the maximum width that the help text will be
	// formated to.
	//
	// 60 is the minimum number of columns that the help text renderer will
	// accept.  If this property is set to a value below 60, the max width will be
	// set 60.
	//
	// For *nix systems, argoutil.GetConsoleWidth is used to try and automatically
	// determine the console width.  If the returned value is greater than 120,
	// the max width will be set to 120.
	HelpTextMaxWidth int
}
