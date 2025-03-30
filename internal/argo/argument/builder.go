package argument

import (
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.ArgumentBuilder {
	return &argumentBuilder{
		marsh: parse.NewDefaultMagicUnmarshaler(),
	}
}

type argumentBuilder struct {
	name string
	desc string

	required   bool
	repeatable bool

	binding Binding
	defVal  Default

	marsh argo.ValueUnmarshaler

	validators []any
}

func (a *argumentBuilder) WithName(name string) argo.ArgumentBuilder {
	a.name = name
	return a
}

func (a *argumentBuilder) HasName() bool {
	return len(a.name) > 0
}

func (a *argumentBuilder) Name() string {
	return a.name
}

func (a *argumentBuilder) WithDescription(desc string) argo.ArgumentBuilder {
	a.desc = desc
	return a
}

func (a *argumentBuilder) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *argumentBuilder) Description() string {
	return a.desc
}

func (a *argumentBuilder) Require() argo.ArgumentBuilder {
	a.required = true
	return a
}

func (a *argumentBuilder) IsRequired() bool {
	return a.required
}

func (a *argumentBuilder) WithBinding(binding any) argo.ArgumentBuilder {
	a.binding = NewBinding(binding)
	return a
}

func (a *argumentBuilder) Binding() argo.ArgumentBinding {
	return &a.binding
}

func (a *argumentBuilder) HasBinding() bool {
	return a.binding.BType != argo.BindingTypeNone
}

func (a *argumentBuilder) WithDefault(def any) argo.ArgumentBuilder {
	a.defVal = NewDefault(def)
	return a
}

func (a *argumentBuilder) HasDefault() bool {
	return a.defVal.Type() != argo.DefaultTypeNone
}

func (a *argumentBuilder) Default() argo.ArgumentDefault {
	return &a.defVal
}

func (a *argumentBuilder) WithUnmarshaler(fn argo.ValueUnmarshaler) argo.ArgumentBuilder {
	a.marsh = fn
	return a
}

func (a *argumentBuilder) HasUnmarshaler() bool {
	return a.marsh != nil
}

func (a *argumentBuilder) Unmarshaler() argo.ValueUnmarshaler {
	return a.marsh
}

func (a *argumentBuilder) WithValidator(fn any) argo.ArgumentBuilder {
	a.validators = append(a.validators, fn)
	return a
}

func (a *argumentBuilder) HasValidators() bool {
	return len(a.validators) > 0
}

func (a *argumentBuilder) Validators() []any {
	return a.validators
}
