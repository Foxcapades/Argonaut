package argo

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

	// CaseSensitiveLong determines whether or not the flag
	// parsing should be case sensitive for long flags.
	//
	// If set to false the following cases will be considered
	// equivalent:
	//     `--help`
	//     `--Help`
	//     `--HELP`
	//     etc.
	//
	// Defaults to true
	CaseSensitiveLong bool

	// CaseSensitiveShort determines whether or not the flag
	// parsing should be case sensitive for short flags.
	//
	// If set to false the following cases will be considered
	// equivalent:
	//     `-h`
	//     `-H`
	//
	// Defaults to true
	CaseSensitiveShort bool
}
