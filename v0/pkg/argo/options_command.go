package argo

func DefaultCommandOptions() CommandOptions {
	return CommandOptions{
		IncludeHelp:     true,
		ArgumentOptions: DefaultArgumentOptions(),
		FlagOptions:     DefaultFlagOptions(),
	}
}

// CommandOptions defines general options for the processing
// of a command and its options.
type CommandOptions struct {

	// IncludeHelp defines whether or not a help flag
	// (-h|--help) should be autogenerated.
	//
	// Defaults to true
	IncludeHelp bool

	ArgumentOptions

	FlagOptions
}
