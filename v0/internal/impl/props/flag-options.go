package props

func DefaultFlagOptions() FlagOptions {
	return FlagOptions{
		RenderArgDescription: true,
	}
}

// FlagOptions defines the options to use when processing /
// rendering flags.
type FlagOptions struct {

	// RenderArgDescription determines whether or not the
	// description for a flag's argument should be rendered as
	// part of its parent flag's description.
	//
	// If this is not set, flag argument descriptions will not
	// be rendered in help/man text at all.
	//
	// Defaults to true
	RenderArgDescription bool
}
