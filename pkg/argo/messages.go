package argo

import "fmt"

var (
	ErrMsgMultiErrorHeaderLine = func(uniqueErrors, totalErrors int) string {
		return fmt.Sprintf("Encountered %d errors (%d unique):", totalErrors, uniqueErrors)
	}

	ErrMsgArgumentDefaultAndRequired = "argument was marked as required, but given a default value"

	ErrMsgFlagHasNoNames = "flag configured with neither a long-form or short-form name"

	ErrMsgInvalidShortFlagName = func(value byte, config Config) string {
		return fmt.Sprintf("invalid flag short byte: \\x%02x", value)
	}

	ErrMsgInvalidLongFlagName = func(value string, config Config) string {
		return fmt.Sprintf("invalid flag long form \"%s%s\"", config.Flags.LongFormPrefix, value)
	}
)
