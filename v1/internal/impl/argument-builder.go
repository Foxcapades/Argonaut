package impl

import "github.com/Foxcapades/Argonaut/v1/pkg/argo"

func NewArgBuilder() *ArgumentBuilder {
	return new(ArgumentBuilder)
}

type bindValidator func(ArgumentBuilder) error

type ArgumentBuilder struct {
	bindValidators []bindValidator

	required bool
	defVal   interface{}
	binding  interface{}
	hintTxt  string
	descTxt  string
}

func (a *ArgumentBuilder) Hint(hint string) argo.ArgumentBuilder {
	a.hintTxt = hint
	return a
}

func (a *ArgumentBuilder) Default(val interface{}) argo.ArgumentBuilder {
	a.defVal = val
	return a
}

func (a *ArgumentBuilder) Bind(ptr interface{}) argo.ArgumentBuilder {
	a.binding = ptr
	return a
}

func (a *ArgumentBuilder) Description(desc string) argo.ArgumentBuilder {
	a.descTxt = desc
	return a
}

func (a *ArgumentBuilder) Require() argo.ArgumentBuilder {
	a.required = true
	return a
}

func (a *ArgumentBuilder) Required(req bool) argo.ArgumentBuilder {
	a.required = req
	return a
}

func (a *ArgumentBuilder) Build() (argo.Argument, error) {
	panic("implement me")
}

func (a *ArgumentBuilder) MustBuild() argo.Argument {
	panic("implement me")
}
