package command

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Incomplete Command Error                                                //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

func IncompleteCommandError(com argo.CommandNode) argo.IncompleteCommandError {
	return incompleteCommandError{com}
}

type incompleteCommandError struct{ com argo.CommandNode }

func (i incompleteCommandError) Error() string {
	return fmt.Sprintf("incomplete command call, reached %s", i.com.Name())
}

func (i incompleteCommandError) StrictOnly() bool { return false }

func (i incompleteCommandError) LastReached() argo.CommandNode { return i.com }
