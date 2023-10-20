package argo

import (
	"fmt"
)

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Flag Error                                                      //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

// A MissingFlagError is returned on CLI parse when a flag that has been marked
// as being required was not found to be present in the CLI call.
//
// MissingFlagError is a hard error that will be returned regardless of whether
// the parser is operating in strict mode.
type MissingFlagError interface {
	error
	Flag() Flag
}

func newMissingFlagError(flag Flag) MissingFlagError {
	return missingFlagError{flag}
}

type missingFlagError struct {
	flag Flag
}

func (m missingFlagError) Error() string {
	return fmt.Sprintf("required flag %s was missing from the CLI call", printFlagNames(m.flag))
}

func (m missingFlagError) StrictOnly() bool {
	return false
}

func (m missingFlagError) Flag() Flag {
	return m.flag
}
