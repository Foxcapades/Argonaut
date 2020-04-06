package impl

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
)

func NewArgBuilder() A.ArgumentBuilder {
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
	name     string
}

func (a *ArgumentBuilder) Name(name string) A.ArgumentBuilder {
	a.name = name;
	return a
}

func (a *ArgumentBuilder) GetName() string {
	return a.name
}

func (a *ArgumentBuilder) HasName() bool {
	return len(a.name) > 0
}

func (a *ArgumentBuilder) Hint(hint string) A.ArgumentBuilder {
	a.hintTxt = hint
	return a
}

func (a *ArgumentBuilder) GetHint() string {
	return a.hintTxt
}

func (a *ArgumentBuilder) HasHint() bool {
	return len(a.hintTxt) > 0
}

func (a *ArgumentBuilder) Default(val interface{}) A.ArgumentBuilder {
	a.hasDef = true
	a.defVal = val
	return a
}

func (a *ArgumentBuilder) Bind(ptr interface{}) A.ArgumentBuilder {
	a.hasBind = true
	a.binding = ptr
	return a
}

func (a *ArgumentBuilder) Description(desc string) A.ArgumentBuilder {
	a.descTxt = desc
	return a
}

func (a *ArgumentBuilder) Require() A.ArgumentBuilder {
	a.required = true
	return a
}

func (a *ArgumentBuilder) Required(req bool) A.ArgumentBuilder {
	a.required = req
	return a
}

func (a *ArgumentBuilder) Build() (A.Argument, error) {
	if a.hasBind {
		// Binding is not usable
		if !util.IsUnmarshalable(a.binding) {
			b := R.TypeOf(a.binding)
			if a.hasDef {
				d := R.TypeOf(a.defVal)
				return nil, A.NewInvalidArgError(A.InvalidArgBindingError, &b, &d)
			} else {
				return nil, A.NewInvalidArgError(A.InvalidArgBindingError, &b, nil)
			}
		}

		// Binding and Default val are incompatible
		if a.hasDef && !util.Compatible(a.binding, a.defVal) {
			b := R.TypeOf(a.binding)
			d := R.TypeOf(a.defVal)
			return nil, A.NewInvalidArgError(A.InvalidArgDefaultError, &b, &d)
		}
	}

	return &Argument{
		name:    a.name,
		defVal:  a.defVal,
		bind:    a.binding,
		hint:    a.hintTxt,
		desc:    a.descTxt,
		isReq:   a.required,
		hasDef:  a.hasDef,
		hasBind: a.hasBind,
	}, nil
}

func (a *ArgumentBuilder) MustBuild() A.Argument {
	if out, err := a.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
