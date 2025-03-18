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

func NewBindingError(root argo.ArgumentBindingError, arg argo.ArgumentBuilder, flag argo.FlagBuilder) argo.BindingError {
	return &flagBindingError{root.Unwrap(), arg, flag}
}

type flagBindingError struct {
	root error
	arg  argo.ArgumentBuilder
	flag argo.FlagBuilder
}

func (f flagBindingError) Error() string {
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

func (f flagBindingError) ArgumentBuilder() argo.ArgumentBuilder {
	return f.arg
}

func (f flagBindingError) FlagBuilder() argo.FlagBuilder {
	return f.flag
}

func (f flagBindingError) Unwrap() error {
	return f.root
}
