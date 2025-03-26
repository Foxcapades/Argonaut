package flag

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewMissingFlagError(flag argo.Flag) argo.MissingFlagError {
	return missingFlagError{flag}
}

type missingFlagError struct {
	flag argo.Flag
}

func (m missingFlagError) Error() string {
	return fmt.Sprintf("required flag %s was missing from the CLI call", PrintFlagNames(m.flag))
}

func (m missingFlagError) Flag() argo.Flag {
	return m.flag
}

//

func NewBindingError(root argo.ArgumentBindingError, arg argo.ArgumentBuilder, flag argo.FlagBuilder) argo.FlagBindingError {
	return &BindingError{root.Unwrap(), arg, flag}
}

type BindingError struct {
	root error
	arg  argo.ArgumentBuilder
	flag argo.FlagBuilder
}

func (f BindingError) Error() string {
	if f.flag.HasLongForm() {
		if f.flag.HasShortForm() {
			return "BindingError (--" + f.flag.LongForm() + "|-" + string([]byte{f.flag.ShortForm()}) + "): " + f.root.Error()
		} else {
			return "BindingError (--" + f.flag.LongForm() + "): " + f.root.Error()
		}
	} else if f.flag.HasShortForm() {
		return "BindingError (-" + string([]byte{f.flag.ShortForm()}) + "): " + f.root.Error()
	} else {
		return "BindingError: " + f.root.Error()
	}
}

func (f BindingError) ArgumentBuilder() argo.ArgumentBuilder {
	return f.arg
}

func (f BindingError) FlagBuilder() argo.FlagBuilder {
	return f.flag
}

func (f BindingError) Unwrap() error {
	return f.root
}
