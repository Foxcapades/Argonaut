package argument

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewMissingRequiredPositionalArgumentError(a argo.Argument, c argo.Command) argo.MissingRequiredArgumentError {
	return &missingArgError{arg: a, com: c}
}

func NewMissingRequiredFlagArgumentError(a argo.Argument, f argo.Flag, c argo.Command) argo.MissingRequiredArgumentError {
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
		return fmt.Sprintf("Missing required argument for flag %s", m.flag)
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
