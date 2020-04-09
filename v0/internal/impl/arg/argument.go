package arg

import (
	"github.com/Foxcapades/Argonaut/v0/internal/impl/trait"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"reflect"
)

type Argument struct {
	trait.Named
	trait.Described

	parent interface{}

	defVal interface{}
	bind   interface{}
	hint   string
	raw    string

	// Flags
	isReq   bool
	hasDef  bool
	hasBind bool

	index uint8
}

func (a *Argument) RawValue() string       { return a.raw }
func (a *Argument) Hint() string           { return a.hint }
func (a *Argument) HasHint() bool          { return len(a.hint) > 0 }
func (a *Argument) Default() interface{}   { return a.defVal }
func (a *Argument) HasDefault() bool       { return a.hasDef }
func (a *Argument) Required() bool         { return a.isReq }
func (a *Argument) SetRawValue(val string) { a.raw = val }
func (a *Argument) Binding() interface{}   { return a.bind }
func (a *Argument) HasBinding() bool       { return a.hasBind }
func (a *Argument) Parent() interface{}    { return a.parent }

func (a *Argument) IsFlagArg() bool {
	if _, ok := a.parent.(A.Flag); ok {
		return true
	}
	return false
}

func (a *Argument) IsPositionalArg() bool {
	if _, ok := a.parent.(A.Command); ok {
		return true
	}
	return false
}

func (a *Argument) BindingType() reflect.Type {
	return util.GetRootValue(reflect.ValueOf(a.Binding())).Type()
}

func (a *Argument) DefaultType() reflect.Type {
	if !a.HasDefault() {
		return nil
	}
	return reflect.TypeOf(a.defVal)
}

func (a *Argument) String() string {
	if a.HasName() {
		return a.Name()
	} else {
		return "arg"
	}
}
