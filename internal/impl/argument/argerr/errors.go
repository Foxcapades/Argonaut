package argerr

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Unexpected Argument Error                                               //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

func NewUnexpectedArgumentError(raw string, flag argo.Flag) argo.UnexpectedArgumentError {
	return unexpectedArgumentError{raw, flag}
}

type unexpectedArgumentError struct {
	raw  string
	flag argo.Flag
}

func (u unexpectedArgumentError) Error() string {
	return fmt.Sprintf("%s does not expect an argument", flagutil.PrintFlagNames(u.flag))
}

func (u unexpectedArgumentError) RawValue() string {
	return u.raw
}

func (u unexpectedArgumentError) Flag() argo.Flag {
	return u.flag
}

func (u unexpectedArgumentError) StrictOnly() bool {
	return true
}

func (u unexpectedArgumentError) HasFlag() bool {
	return u.flag != nil
}

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Argument Error                                                  //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

func MissingRequiredPositionalArgumentError(a argo.Argument, c argo.Command) error {
	return &missingArgError{arg: a, com: c}
}

func MissingRequiredFlagArgumentError(a argo.Argument, f argo.Flag, c argo.Command) error {
	return &missingArgError{arg: a, flag: f, com: c}
}

type missingArgError struct {
	arg  argo.Argument
	flag argo.Flag
	com  argo.Command
}

func (m *missingArgError) Flag() argo.Flag {
	return m.flag
}

func (m *missingArgError) HasFlag() bool {
	return m.flag != nil
}

func (m *missingArgError) StrictOnly() bool {
	return false
}

func (m *missingArgError) Error() string {
	if m.flag != nil {
		return fmt.Sprintf("Missing required argument for flag %s", flagutil.PrintFlagNames(m.flag))
	} else if m.arg.HasName() {
		return fmt.Sprintf("Missing required positional argument %s", m.arg.Name())
	} else {
		for i, a := range m.com.Arguments() {
			if m.arg == a {
				return fmt.Sprintf("Missing required positional argument #%d", i+1)
			}
		}
		return "Missing required positional argument"
	}
}

func (m *missingArgError) Argument() argo.Argument {
	return m.arg
}
