package flag

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// NewBuilder returns a new FlagBuilder instance.
func NewBuilder() argo.FlagBuilder {
	return &flagBuilder{}
}

type flagBuilder struct {
	short  byte
	req    bool
	isHelp bool
	long   string
	desc   string
	onHit  argo.FlagCallback
	arg    argo.ArgumentBuilder
}

func (b *flagBuilder) WithShortForm(char byte) argo.FlagBuilder {
	b.short = char
	return b
}

func (b *flagBuilder) HasShortForm() bool {
	return b.short != 0
}

func (b *flagBuilder) ShortForm() byte {
	return b.short
}

func (b *flagBuilder) WithLongForm(form string) argo.FlagBuilder {
	b.long = form
	return b
}

func (b *flagBuilder) HasLongForm() bool {
	return len(b.long) > 0
}

func (b *flagBuilder) LongForm() string {
	return b.long
}

//

func (b *flagBuilder) WithDescription(desc string) argo.FlagBuilder {
	b.desc = desc
	return b
}

func (b *flagBuilder) HasDescription() bool {
	return len(b.desc) > 0
}

func (b *flagBuilder) Description() string {
	return b.desc
}

//

func (b *flagBuilder) WithCallback(fn argo.FlagCallback) argo.FlagBuilder {
	b.onHit = fn
	return b
}

func (b *flagBuilder) HasCallback() bool {
	return b.onHit != nil
}

func (b *flagBuilder) Callback() argo.FlagCallback {
	return b.onHit
}

//

func (b *flagBuilder) WithArgument(arg argo.ArgumentBuilder) argo.FlagBuilder {
	b.arg = arg
	return b
}

func (b *flagBuilder) HasArgument() bool {
	return b.arg != nil
}

func (b *flagBuilder) Argument() argo.ArgumentBuilder {
	return b.arg
}

//

func (b *flagBuilder) Require() argo.FlagBuilder {
	b.req = true
	return b
}

func (b *flagBuilder) IsRequired() bool {
	return b.req
}

//

func (b *flagBuilder) WithBinding(pointer any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *flagBuilder) WithBindingAndDefault(pointer, def any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer).WithDefault(def)

	if required {
		b.arg.Require()
	}

	return b
}

//

func (b *flagBuilder) MarkAsHelpFlag() argo.FlagBuilder {
	b.isHelp = true
	return b
}

func (b *flagBuilder) IsHelpFlag() bool {
	return b.isHelp
}

//

func (b *flagBuilder) Build(ctx *argo.WarningContext) (argo.Flag, error) {
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

	if !b.HasShortForm() && !b.HasLongForm() {
		errs.AppendError(errors.New("flag declared with neither a long or short form"))
	}

	var arg argo.Argument

	if b.arg != nil {
		var err error
		arg, err = b.arg.Build(ctx)

		if err != nil {
			var e argo.MultiError
			if errors.As(err, &e) {
				var be argo.ArgumentBindingError
				for _, err := range e.Errors() {
					if errors.As(err, &be) {
						errs.AppendError(NewBindingError(be, b.arg, b))
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

	return &flag{
		warnings: ctx,
		short:    b.short,
		required: b.req,
		arg:      arg,
		long:     b.long,
		desc:     b.desc,
		isHelp:   b.isHelp,
		callback: b.onHit,
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
