package argo

type ParseResult struct {
	Warnings []ParseWarning
	Error    error
}

const (
	// UnexpectedFlagArgument parse warnings indicate that a flag was explicitly
	// passed an argument in the CLI call when it wasn't expecting one.
	UnexpectedFlagArgument ParseWarningType = iota

	// UnrecognizedFlag parse warnings indicate that a flag was passed in the CLI
	// call that was not registered when building the command or command tree.
	UnrecognizedFlag
)

type ParseWarningType uint8

func (p ParseWarningType) String() string {
	switch p {
	case UnexpectedFlagArgument:
		return "unexpected flag argument"
	case UnrecognizedFlag:
		return "unrecognized flag"
	default:
		return "invalid warning type"
	}
}

type ParseWarning struct {
	Type    ParseWarningType
	Message string
}
