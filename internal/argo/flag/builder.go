package flag

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
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
