package argo

type TabCompletionType uint8

const (
	TabCompleteFiles TabCompletionType = 1 << iota
	TabCompleteDirectories

	TabCompleteNone TabCompletionType = 0
)

type FlagCompletionHints struct {
	// Repeatable indicates that the target flag, subcommand or argument may be
	// the set number of times.
	//
	// A zero value indicates that the flag or argument may be repeated any amount
	// of times.
	//
	// This setting has no effect on subcommands.
	Repeatable uint8

	// HideFromTabCompletion indicates whether the target flag, subcommand, or
	// argument should be hidden from tab completion.
	HideFromTabCompletion bool
}

type ArgumentCompletionHints struct {
	// Repeatable indicates that the target flag, subcommand or argument may be
	// the set number of times.
	//
	// A zero value indicates that the flag or argument may be repeated any amount
	// of times.
	//
	// This setting has no effect on subcommands.
	Repeatable uint8

	// HideFromTabCompletion indicates whether the target flag, subcommand, or
	// argument should be hidden from tab completion.
	HideFromTabCompletion bool

	// CompletionType provides automatic completion hints for arguments.
	//
	// Completion types may be used in combination.
	CompletionType TabCompletionType

	// CompletionProvider is a function that may be provided to filter or produce
	// tab completion suggestions.
	//
	// The function may be passed auto-generated hints if the CompletionType value
	// is set.
	//
	// For example, if the CompletionType value is TabCompleteFiles, then this
	// function will be called with the files available based on the path the
	// user has provided (cwd files if the input path is empty).
	//
	// If no CompletionType is set, the passed argument will be nil.
	CompletionProvider func(suggestions []string) []string
}

type SubcommandCompletionHints struct {
	// HideFromTabCompletion indicates whether the target flag, subcommand, or
	// argument should be hidden from tab completion.
	HideFromTabCompletion bool
}

type CommandCompletionHints struct {
	DisableCompletion bool
}
