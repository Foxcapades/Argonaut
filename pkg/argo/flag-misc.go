package argo

// A FlagCallback is a function that, if set on a flag, will be called by the
// CLI parsing process if that flag is used in the CLI call.
//
// The flag callback will be called after CLI parsing has completed.
type FlagCallback = func(flag Flag)

// A MissingFlagError is returned on CLI parse when a flag that has been marked
// as being required was not found to be present in the CLI call.
//
// MissingFlagError is a hard error that will be returned regardless of whether
// the parser is operating in strict mode.
type MissingFlagError interface {
	error

	Flag() Flag
}

type FlagBindingError interface {
	ArgumentBindingError

	FlagBuilder() FlagBuilder
}
