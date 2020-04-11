package argument

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	R "reflect"

	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

func NewBuilder(A.Provider) A.ArgumentBuilder {
	return new(Builder)
}

type Builder struct {
	ParentElement interface{}

	IsArgRequired bool

	Error error

	IsDefaultSet bool
	DefaultValue interface{}
	RootDefault  R.Value

	IsBindingSet bool
	BindValue    interface{}
	RootBinding  R.Value

	DescriptionValue trait.Described

	NameValue trait.Named
}

func (a *Builder) Name(name string) A.ArgumentBuilder { a.NameValue.NameTxt = name; return a }

func (a *Builder) Default(val interface{}) A.ArgumentBuilder {
	a.IsDefaultSet = true
	a.DefaultValue = val
	return a
}

func (a *Builder) HasDefaultProvider() bool {
	return a.IsDefaultSet && util.GetRootValue(R.ValueOf(a.DefaultValue)).Kind() == R.Func
}

func (a *Builder) Bind(ptr interface{}) A.ArgumentBuilder {
	a.IsBindingSet = true
	a.BindValue = ptr
	return a
}

func (a *Builder) Description(desc string) A.ArgumentBuilder {
	a.DescriptionValue.DescTxt = desc
	return a
}

func (a *Builder) Require() A.ArgumentBuilder {
	a.IsArgRequired = true
	return a
}

func (a *Builder) Required(req bool) A.ArgumentBuilder {
	a.IsArgRequired = req
	return a
}

func (a *Builder) IsRequired() bool {
	return a.IsArgRequired
}

func (a *Builder) Parent(par interface{}) A.ArgumentBuilder {
	a.ParentElement = par
	return a
}

func (a *Builder) GetName() string         { return a.NameValue.NameTxt }
func (a *Builder) HasName() bool           { return a.NameValue.HasName() }
func (a *Builder) GetDefault() interface{} { return a.DefaultValue }
func (a *Builder) HasDefault() bool        { return a.IsDefaultSet }
func (a *Builder) GetBinding() interface{} { return a.BindValue }
func (a *Builder) HasBinding() bool        { return a.IsBindingSet }

func (a *Builder) Build() (A.Argument, error) {
	if err := a.ValidateBinding(); err != nil {
		return nil, err
	}

	if err := a.ValidateDefault(); err != nil {
		return nil, err
	}

	if err := a.ValidateParent(); err != nil {
		return nil, err
	}

	return &Argument{
		Named:     a.NameValue,
		Described: a.DescriptionValue,

		DefaultValue: a.DefaultValue,
		RootDefault:  a.RootDefault,
		IsDefaultSet: a.IsDefaultSet,
		IsBindingSet: a.IsBindingSet,

		BindValue:     a.BindValue,
		RootBinding:   a.RootBinding,
		IsRequired:    a.IsArgRequired,
		ParentElement: a.ParentElement,
	}, nil
}

func (a *Builder) MustBuild() A.Argument {
	if out, err := a.Build(); err != nil {
		panic(err)
	} else {
		return out
	}
}
