package argo

import "github.com/foxcapades/argonaut/v3/pkg/argoutil"

type Options struct {
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

	// DefaultCommandGroupName defines the name of the default group for
	// subcommands under a TreeCommand.
	//
	// This option is not used for simple Command instances.
	DefaultCommandGroupName string

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

	// InheritParentFlags sets whether subcommands under TreeCommand instances
	// should "inherit" their parent's flags.
	//
	// The inherited flag will be process in the context of the parent node even
	// if they appear after the child command name.
	//
	// If a child command has a flag whose short and/or long form conflict with a
	// parent's flag, the conflicting form(s) from the parent will be overridden
	// by the child's flag.
	//
	// This setting does nothing for Command instances.
	//
	// Default value is FlagInheritanceEnabled.
	//
	// Example: Inherit flags enabled
	//   // "tree"   defines the flag -a and -b
	//   // "branch" defines the flag -b and -c
	//   // "leaf"   defines the flag -c
	//
	//   $ tree branch leaf -a -b -c
	//
	//   // tree:   -a=true, -b=false
	//   // branch: -b=true, -c=false
	//   // leaf:   -c=true
	//
	// Example: Inherit flags disabled
	//   // "tree"   defines the flag -a
	//   // "branch" defines the flag -b
	//   // "leaf"   defines the flag -c
	//
	//   $ tree branch leaf -a -b -c
	//
	//   // tree:   -a=false
	//   // branch: -b=false
	//   // leaf:   -c=true
	InheritParentFlags ParentFlagInheritanceMode
}

func DefaultOptions() Options {
	consoleWidth, _ := argoutil.GetConsoleWidth()
	return Options{
		MaxDefaultFlagGroupSizeForMetaGroup: 5,
		MetaFlagGroupName:                   "General Flags",
		DefaultCommandGroupName:             "Commands",
		HelpTextMaxWidth:                    max(min(consoleWidth, 120), 60),
		InheritParentFlags:                  FlagInheritanceEnabled,
	}
}

func FixOptions(opts *Options) {
	opts.HelpTextMaxWidth = max(60, opts.HelpTextMaxWidth)
}

// ParentFlagInheritanceMode defines the behavior of parent flag inheritance in
// child commands under a TreeCommand instance.
type ParentFlagInheritanceMode uint8

const (
	// FlagInheritanceDisabled completely disables flag inheritance of parent
	// command flags in child commands.
	FlagInheritanceDisabled ParentFlagInheritanceMode = iota

	// FlagInheritanceEnabled enables flag inheritance of parent command flags in
	// child commands.
	//
	// Inherited flags will be grouped together in a single collection in rendered
	// help text.
	FlagInheritanceEnabled

	// FlagInheritanceHidden enables flag inheritance of parent command flags in
	// child commands, but does not show the inherited flags in help text.
	FlagInheritanceHidden

	// FlagInheritanceGrouped enables flag inheritance of parent command flags in
	// child commands, and groups the inherited flags by the owning parent in help
	// text.
	FlagInheritanceGrouped
)

func (p ParentFlagInheritanceMode) ShouldRender() bool {
	return p == FlagInheritanceEnabled || p == FlagInheritanceGrouped
}

func (p ParentFlagInheritanceMode) String() string {
	switch p {
	case FlagInheritanceEnabled:
		return "FlagInheritanceEnabled"
	case FlagInheritanceDisabled:
		return "FlagInheritanceDisabled"
	case FlagInheritanceHidden:
		return "FlagInheritanceHidden"
	case FlagInheritanceGrouped:
		return "FlagInheritanceGrouped"
	default:
		return "Unknown"
	}
}
