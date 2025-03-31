package flag

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// NewBuilder returns a new FlagBuilder instance.
func NewBuilder() argo.FlagBuilder {
	return new(Builder)
}

type Builder struct {
	short  byte
	req    bool
	isHelp bool
	long   string
	desc   string
	onHit  argo.FlagCallback
	arg    argo.ArgumentBuilder
}

func (b *Builder) WithShortForm(char byte) argo.FlagBuilder {
	b.short = char
	return b
}

func (b *Builder) HasShortForm() bool {
	return b.short != 0
}

func (b *Builder) ShortForm() byte {
	return b.short
}

func (b *Builder) WithLongForm(form string) argo.FlagBuilder {
	b.long = form
	return b
}

func (b *Builder) HasLongForm() bool {
	return len(b.long) > 0
}

func (b *Builder) LongForm() string {
	return b.long
}

func (b *Builder) WithDescription(desc string) argo.FlagBuilder {
	b.desc = desc
	return b
}

func (b *Builder) HasDescription() bool {
	return len(b.desc) > 0
}

func (b *Builder) Description() string {
	return b.desc
}

func (b *Builder) WithCallback(fn argo.FlagCallback) argo.FlagBuilder {
	b.onHit = fn
	return b
}

func (b *Builder) HasCallback() bool {
	return b.onHit != nil
}

func (b *Builder) Callback() argo.FlagCallback {
	return b.onHit
}

func (b *Builder) WithArgument(arg argo.ArgumentBuilder) argo.FlagBuilder {
	b.arg = arg
	return b
}

func (b *Builder) HasArgument() bool {
	return b.arg != nil
}

func (b *Builder) Argument() argo.ArgumentBuilder {
	return b.arg
}

func (b *Builder) Require() argo.FlagBuilder {
	b.req = true
	return b
}

func (b *Builder) IsRequired() bool {
	return b.req
}

func (b *Builder) WithBinding(pointer any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *Builder) WithBindingAndDefault(pointer, def any, required bool) argo.FlagBuilder {
	b.arg = argument.NewBuilder().WithBinding(pointer).WithDefault(def)

	if required {
		b.arg.Require()
	}

	return b
}

func (b *Builder) MarkAsHelpFlag() argo.FlagBuilder {
	b.isHelp = true
	return b
}

func (b *Builder) IsHelpFlag() bool {
	return b.isHelp
}
