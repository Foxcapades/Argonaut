package argument

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"reflect"
)

type Argument struct {
	trait.Named
	trait.Described

	ParentElement interface{}

	RawInput string

	IsRequired bool

	BindValue    interface{}
	RootBinding  reflect.Value
	IsBindingSet bool

	DefaultValue interface{}
	RootDefault  reflect.Value
	IsDefaultSet bool
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Interface Implementation                                            ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (a *Argument) Required() bool         { return a.IsRequired }
func (a *Argument) Parent() interface{}    { return a.ParentElement }

func (a *Argument) RawValue() string       { return a.RawInput }
func (a *Argument) SetRawValue(val string) { a.RawInput = val }

func (a *Argument) Default() interface{}   { return a.DefaultValue }
func (a *Argument) RootDefaultValue() reflect.Value { return a.RootDefault }
func (a *Argument) HasDefault() bool       { return a.IsDefaultSet }

func (a *Argument) Binding() interface{}   { return a.BindValue }
func (a *Argument) RootBindValue() reflect.Value { return a.RootBinding }
func (a *Argument) HasBinding() bool       { return a.IsBindingSet }

func (a *Argument) IsFlagArg() bool {
	if _, ok := a.ParentElement.(A.Flag); ok {
		return true
	}
	return false
}

func (a *Argument) IsPositionalArg() bool {
	if _, ok := a.ParentElement.(A.Command); ok {
		return true
	}
	return false
}

func (a *Argument) BindingType() reflect.Type {
	if !a.HasBinding() {
		return nil
	}

	return a.RootBinding.Type()
}

func (a *Argument) DefaultType() reflect.Type {
	if !a.HasDefault() {
		return nil
	}

	return a.RootDefault.Type()
}

func (a *Argument) String() string {
	if a.HasName() {
		return a.Name()
	} else {
		return "arg"
	}
}
