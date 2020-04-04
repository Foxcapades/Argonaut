package impl

import (
	"reflect"
)

func NewArg() *Argument {
	return new(Argument)
}

type argumentProp = uint8

const (
	argIsReq argumentProp = 1 << iota
)

type Argument struct {
	props  argumentProp
	defVal interface{}
	bind   interface{}
	hint   string
	desc   string
}

func (a *Argument) RawValue() string {
	panic("implement me")
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
	return a.defVal != nil
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
	return argIsReq == argIsReq&a.props
}

func (a *Argument) binding() interface{} {
	return a.bind
}

func (a *Argument) hasBinding() bool {
	return a.bind != nil
}
