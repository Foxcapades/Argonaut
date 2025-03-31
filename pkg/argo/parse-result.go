package argo

type ParseResult struct {
	InputWarnings []InputWarning

	Error error

	ErrorType ErrorType
}

type ErrorType uint8

const (
	ErrorTypeNone ErrorType = iota
	ConfigurationError
	InputError
)

func (t ErrorType) String() string {
	switch t {
	case ErrorTypeNone:
		return "ErrorTypeNone"
	case ConfigurationError:
		return "ConfigurationError"
	case InputError:
		return "InputError"
	default:
		return "InvalidErrorType"
	}
}

type InputWarningType uint8

const (
	// UnexpectedFlagArgument parse warnings indicate that a flag was explicitly
	// passed an argument in the CLI call when it wasn't expecting one.
	UnexpectedFlagArgument InputWarningType = iota

	// UnrecognizedFlag parse warnings indicate that a flag was passed in the CLI
	// call that was not registered when building the command or command tree.
	UnrecognizedFlag
)

func (p InputWarningType) String() string {
	switch p {
	case UnexpectedFlagArgument:
		return "unexpected flag argument"
	case UnrecognizedFlag:
		return "unrecognized flag"
	default:
		return "invalid warning type"
	}
}

type InputWarning struct {
	Type    InputWarningType
	Message string
}
