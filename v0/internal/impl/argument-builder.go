package impl

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
)

func NewArgBuilder() A.ArgumentBuilder {
	return new(ArgumentBuilder)
}

type ArgumentBuilder struct {
	parent interface{}

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
	a.name = name
	return a
}

func (a *ArgumentBuilder) TypeHint(hint string) A.ArgumentBuilder {
	a.hintTxt = hint
	return a
}

func (a *ArgumentBuilder) Default(val interface{}) A.ArgumentBuilder {
	a.hasDef = true
	a.defVal = val
	return a
}

func (a *ArgumentBuilder) HasDefaultProvider() bool {
	return a.hasDef && R.TypeOf(a.defVal).Elem().Kind() == R.Func
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

func (a *ArgumentBuilder) Parent(par interface{}) A.ArgumentBuilder {
	a.parent = par
	return a
}

func (a *ArgumentBuilder) GetName() string         { return a.name }
func (a *ArgumentBuilder) HasName() bool           { return len(a.name) > 0 }
func (a *ArgumentBuilder) GetHint() string         { return a.hintTxt }
func (a *ArgumentBuilder) HasHint() bool           { return len(a.hintTxt) > 0 }
func (a *ArgumentBuilder) GetDefault() interface{} { return a.defVal }
func (a *ArgumentBuilder) HasDefault() bool        { return a.hasDef }
func (a *ArgumentBuilder) GetBinding() interface{} { return a.binding }
func (a *ArgumentBuilder) HasBinding() bool        { return a.hasBind }

func (a *ArgumentBuilder) Build() (A.Argument, error) {
	if a.hasBind {
		// Binding is not usable
		if !util.IsUnmarshalable(a.binding) {
			return nil, A.NewInvalidArgError(A.ArgErrInvalidBindingBadType, a, "")
		}

		if a.hasDef {
			if err := checkDefault(a); err != nil {
				return nil, err
			}
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
