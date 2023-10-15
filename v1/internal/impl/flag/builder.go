package flag

import (
	"errors"

	"github.com/Foxcapades/Argonaut/v1/internal/chars"
	"github.com/Foxcapades/Argonaut/v1/internal/impl/argument"
	"github.com/Foxcapades/Argonaut/v1/internal/xerr"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

func NewBuilder() argo.FlagBuilder {
	return &builder{}
}

type builder struct {
	short byte
	req   bool
	long  string
	desc  string
	onHit func(argo.Flag)
	arg   argo.ArgumentBuilder
}

func (b *builder) WithShortForm(char byte) argo.FlagBuilder {
	b.short = char
	return b
}

func (b builder) HasShortForm() bool {
	return b.short != 0
}

func (b builder) GetShortForm() byte {
	return b.short
}

func (b *builder) WithLongForm(form string) argo.FlagBuilder {
	b.long = form
	return b
}

func (b builder) HasLongForm() bool {
	return len(b.long) > 0
}

func (b builder) GetLongForm() string {
	return b.long
}

func (b *builder) WithDescription(desc string) argo.FlagBuilder {
	b.desc = desc
	return b
}

func (b builder) HasDescription() bool {
	return len(b.desc) > 0
}

func (b *builder) WithOnHitCallback(fn func(argo.Flag)) argo.FlagBuilder {
	b.onHit = fn
	return b
}

func (b *builder) WithArgument(arg argo.ArgumentBuilder) argo.FlagBuilder {
	b.arg = arg
	return b
}

func (b *builder) Require() argo.FlagBuilder {
	b.req = true
	return b
}

func (b *builder) WithBinding(pointer any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *builder) WithBindingAndDefault(pointer, def any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer).WithDefault(def)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *builder) Build() (argo.Flag, error) {
	errs := xerr.NewMultiError()

	if b.short > 0 {
		if err := validateShortForm(b.short); err != nil {
			errs.AppendError(err)
		}
	}

	if len(b.long) > 0 {
		if err := validateLongForm(b.long); err != nil {
			errs.AppendError(err)
		}
	}

	var arg argo.Argument

	if b.arg != nil {
		var err error
		arg, err = b.arg.Build()
		if err != nil {
			errs.AppendError(err)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &flag{
		short:    b.short,
		required: b.req,
		arg:      arg,
		long:     b.long,
		desc:     b.desc,
	}, nil
}

func validateShortForm(c byte) error {
	if !chars.IsAlphanumeric(c) {
		return errors.New("short-form flags must be alphanumeric")
	}

	return nil
}

func validateLongForm(f string) error {
	if !chars.IsAlphanumeric(f[0]) {
		return errors.New("long-form flags must begin with an alphanumeric character")
	}

	for i := 1; i < len(f); i++ {
		if !chars.IsFlagStringSafe(f[i]) {
			return errors.New("long-form flags must only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}
