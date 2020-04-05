package impl

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	"reflect"
)

type Argument struct {
	defVal interface{}
	bind   interface{}
	hint   string
	desc   string
	raw    string
	name   string

	// Flags
	isReq   bool
	hasDef  bool
	hasBind bool

	index uint8
}

func (a *Argument) Name() string {
	return a.name
}

func (a *Argument) HasName() bool {
	return len(a.name) > 0
}

func (a *Argument) RawValue() string {
	return a.raw
}

func (a *Argument) Hint() string {
	return a.hint
}

func (a *Argument) HasHint() bool {
	return len(a.hint) > 0
}

func (a *Argument) Default() interface{} {
	return a.defVal
}

func (a *Argument) HasDefault() bool {
	return a.hasDef
}

func (a *Argument) DefaultType() reflect.Type {
	if !a.HasDefault() {
		return nil
	}
	return reflect.TypeOf(a.defVal)
}

func (a *Argument) Description() string {
	return a.desc
}

func (a *Argument) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *Argument) Required() bool {
	return a.isReq
}

func (a *Argument) SetRawValue(val string) {
	a.raw = val
}

func (a *Argument) Binding() interface{} {
	return a.bind
}

func (a *Argument) HasBinding() bool {
	return a.hasBind
}

func (a *Argument) BindingType() reflect.Type {
	return util.GetRootValue(reflect.ValueOf(a.Binding())).Type()
}

func (a *Argument) String() string {
	var name string

	if a.HasName() {
		name = a.Name()
	} else {
		name = "arg"
	}

	if a.isReq {
		return name
	} else {
		return "[" + name + "]"
	}
}

