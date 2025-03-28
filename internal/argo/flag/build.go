package flag

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(b argo.FlagBuilder) (argo.Flag, error) {
	errs := xerr.NewMultiError()

	if b.HasShortForm() {
		xerr.AppendIfPresent(validateShortForm(b.ShortForm()), errs)
	}

	if b.HasLongForm() {
		xerr.AppendIfPresent(validateLongForm(b.LongForm()), errs)
	}

	if !b.HasShortForm() && !b.HasLongForm() {
		errs.AppendError(errors.New("flag declared with neither a long or short form"))
	}

	var arg argo.Argument

	if b.HasArgument() {
		var err error
		arg, err = argument.Build(b.Argument())

		if err != nil {
			var e argo.MultiError
			if errors.As(err, &e) {
				var be argo.ArgumentBindingError
				for _, err = range e.Errors() {
					if errors.As(err, &be) {
						errs.AppendError(NewBindingError(be, b.Argument(), b))
					} else {
						errs.AppendError(err)
					}
				}
			} else {
				errs.AppendError(err)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &Flag{
		short:    b.ShortForm(),
		required: b.IsRequired(),
		arg:      arg,
		long:     b.LongForm(),
		desc:     b.Description(),
		isHelp:   b.IsHelpFlag(),
		callback: b.Callback(),
	}, nil
}

func validateShortForm(c byte) error {
	if !text.IsAlphanumeric(c) {
		return errors.New("short-form flags must be alphanumeric")
	}

	return nil
}

func validateLongForm(f string) error {
	if !text.IsAlphanumeric(f[0]) {
		return errors.New("long-form flags must begin with an alphanumeric character")
	}

	for i := 1; i < len(f); i++ {
		if !(text.IsWord(f[i]) || f[i] == text.DashByte) {
			return errors.New("long-form flags must only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}
