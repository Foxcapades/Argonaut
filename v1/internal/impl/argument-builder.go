package impl

import (
	"github.com/Foxcapades/Argonaut/v1/internal/util"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
	R "reflect"
)

func NewArgBuilder() argo.ArgumentBuilder {
	return new(ArgumentBuilder)
}

type bindValidator func(ArgumentBuilder) error

type ArgumentBuilder struct {
	bindValidators []bindValidator

	required bool
	hasDef   bool
	hasBind  bool
	error    error
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
	a.hasDef = true
	a.defVal = val
	return a
}

func (a *ArgumentBuilder) Bind(ptr interface{}) argo.ArgumentBuilder {
	a.hasBind = true
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
	if a.hasBind {
		// Binding is not usable
		if !util.IsUnmarshalable(a.binding) {
			b := R.TypeOf(a.binding)
			if a.hasDef {
				d := R.TypeOf(a.defVal)
				return nil, argo.NewInvalidArgError(argo.InvalidArgBindingError, &b, &d)
			} else {
				return nil, argo.NewInvalidArgError(argo.InvalidArgBindingError, &b, nil)
			}
		}

		// Binding and Default val are incompatible
		if a.hasDef && !util.Compatible(a.binding, a.defVal) {
			b := R.TypeOf(a.binding)
			d := R.TypeOf(a.defVal)
			return nil, argo.NewInvalidArgError(argo.InvalidArgDefaultError, &b, &d)
		}
	}

	return &Argument{
		defVal:  a.defVal,
		bind:    a.binding,
		hint:    a.hintTxt,
		desc:    a.descTxt,
		isReq:   a.required,
		hasDef:  a.hasDef,
		hasBind: a.hasBind,
	}, nil
}

func (a *ArgumentBuilder) MustBuild() argo.Argument {
	if out, err := a.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
