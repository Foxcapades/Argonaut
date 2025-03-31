package argument

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.ArgumentBuilder {
	return new(Builder)
}

type Builder struct {
	name string
	desc string

	required   bool
	repeatable bool

	binding Binding
	defVal  Default

	marsh argo.ValueUnmarshaler

	validators []any
}

func (a *Builder) WithName(name string) argo.ArgumentBuilder {
	a.name = name
	return a
}

func (a *Builder) HasName() bool {
	return len(a.name) > 0
}

func (a *Builder) Name() string {
	return a.name
}

func (a *Builder) WithDescription(desc string) argo.ArgumentBuilder {
	a.desc = desc
	return a
}

func (a *Builder) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *Builder) Description() string {
	return a.desc
}

func (a *Builder) Require() argo.ArgumentBuilder {
	a.required = true
	return a
}

func (a *Builder) IsRequired() bool {
	return a.required
}

func (a *Builder) WithBinding(binding any) argo.ArgumentBuilder {
	a.binding = NewBinding(binding)
	return a
}

func (a *Builder) Binding() argo.ArgumentBinding {
	return &a.binding
}

func (a *Builder) HasBinding() bool {
	return a.binding.BType != argo.BindingTypeNone
}

func (a *Builder) WithDefault(def any) argo.ArgumentBuilder {
	a.defVal = NewDefault(def)
	return a
}

func (a *Builder) HasDefault() bool {
	return a.defVal.Type() != argo.DefaultTypeNone
}

func (a *Builder) Default() argo.ArgumentDefault {
	return &a.defVal
}

func (a *Builder) WithUnmarshaler(fn argo.ValueUnmarshaler) argo.ArgumentBuilder {
	a.marsh = fn
	return a
}

func (a *Builder) HasUnmarshaler() bool {
	return a.marsh != nil
}

func (a *Builder) Unmarshaler() argo.ValueUnmarshaler {
	return a.marsh
}

func (a *Builder) WithValidator(fn any) argo.ArgumentBuilder {
	a.validators = append(a.validators, fn)
	return a
}

func (a *Builder) HasValidators() bool {
	return len(a.validators) > 0
}

func (a *Builder) Validators() []any {
	return a.validators
}
