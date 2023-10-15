package flag

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Flag Error                                                      //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

func MissingFlagError(flag argo.Flag) argo.MissingFlagError {
	return missingFlagError{flag}
}

type missingFlagError struct {
	flag argo.Flag
}

func (m missingFlagError) Error() string {
	return fmt.Sprintf("required flag %s was missing from the CLI call", flagutil.PrintFlagNames(m.flag))
}

func (m missingFlagError) StrictOnly() bool {
	return false
}

func (m missingFlagError) Flag() argo.Flag {
	return m.flag
}
