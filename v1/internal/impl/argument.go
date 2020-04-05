package impl

import (
	"reflect"
)

type Argument struct {
	defVal interface{}
	bind   interface{}
	hint   string
	desc   string
	raw    string

	// Flags
	isReq   bool
	hasDef  bool
	hasBind bool
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

func (a *Argument) Binding() interface{} {
	return a.bind
}

func (a *Argument) HasBinding() bool {
	return a.hasBind
}
