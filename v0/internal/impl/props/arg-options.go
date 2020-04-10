package props

func DefaultArgumentOptions() ArgumentOptions {
	return ArgumentOptions{
		UseTypeAsDefault: true,
		DefaultName:      "arg",
	}
}

// ArgumentOptions defines options to use when processing
// rendering arguments
type ArgumentOptions struct {
	// UseTypeAsDefault defines whether or not an argument's
	// computed binding type should be used as a standin when
	// rendering an argument that has no name set for help/man
	// purposes
	//
	// Example:
	//     // named arg
	//     --some-flag=<arg-name>
	//
	//     // non-named arg
	//     --some-flag=<int>
	//
	// Defaults to true
	UseTypeAsDefault bool

	// DefaultName defines the standin name to use when
	// rendering a nameless argument in help / man contexts
	// that has no name.
	//
	// Defaults to "arg"
	DefaultName string
}
