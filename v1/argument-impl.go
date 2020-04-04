package argo

import (
	"reflect"
)

type argumentProp = uint8

const (
	argIsReq argumentProp = 1 << iota
)

type argument struct {
	props  argumentProp
	defVal interface{}
	bind   interface{}
	hint   string
	desc   string
}

func (a *argument) Hint() string {
	return a.hint
}

func (a *argument) HasHint() bool {
	return len(a.hint) > 0
}

func (a *argument) Default() interface{} {
	return a.defVal
}

func (a *argument) HasDefault() bool {
	return a.defVal != nil
}

func (a *argument) DefaultType() reflect.Type {
	if !a.HasDefault() {
		return nil
	}
	return reflect.TypeOf(a.defVal)
}

func (a *argument) Description() string {
	return a.desc
}

func (a *argument) HasDescription() bool {
	return len(a.desc) > 0
}

func (a *argument) Required() bool {
	return argIsReq == argIsReq & a.props
}

func (a *argument) binding() interface{} {
	return a.bind
}

func (a *argument) hasBinding() bool {
	return a.bind != nil
}