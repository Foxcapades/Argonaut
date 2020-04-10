package argument

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	R "reflect"

	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewBuilder(A.Provider) A.ArgumentBuilder {
	return new(Builder)
}

type Builder struct {
	parent interface{}

	required bool
	hasDef   bool
	hasBind  bool
	error    error
	defVal   interface{}
	binding  interface{}
	desc     trait.Described
	name     trait.Named
}

func (a *Builder) Name(name string) A.ArgumentBuilder { a.name.NameValue = name; return a }

func (a *Builder) Default(val interface{}) A.ArgumentBuilder {
	a.hasDef = true
	a.defVal = val
	return a
}

func (a *Builder) HasDefaultProvider() bool {
	return a.hasDef && R.TypeOf(a.defVal).Elem().Kind() == R.Func
}

func (a *Builder) Bind(ptr interface{}) A.ArgumentBuilder {
	a.hasBind = true
	a.binding = ptr
	return a
}

func (a *Builder) Description(desc string) A.ArgumentBuilder {
	a.desc.DescriptionValue = desc
	return a
}

func (a *Builder) Require() A.ArgumentBuilder {
	a.required = true
	return a
}

func (a *Builder) Required(req bool) A.ArgumentBuilder {
	a.required = req
	return a
}

func (a *Builder) Parent(par interface{}) A.ArgumentBuilder {
	a.parent = par
	return a
}

func (a *Builder) GetName() string         { return a.name.NameValue }
func (a *Builder) HasName() bool           { return a.name.HasName() }
func (a *Builder) GetDefault() interface{} { return a.defVal }
func (a *Builder) HasDefault() bool        { return a.hasDef }
func (a *Builder) GetBinding() interface{} { return a.binding }
func (a *Builder) HasBinding() bool        { return a.hasBind }

func (a *Builder) Build() (A.Argument, error) {
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
		Named:     a.name,
		defVal:    a.defVal,
		bind:      a.binding,
		Described: a.desc,
		isReq:     a.required,
		hasDef:    a.hasDef,
		hasBind:   a.hasBind,
		parent:    a.parent,
	}, nil
}

func (a *Builder) MustBuild() A.Argument {
	if out, err := a.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
