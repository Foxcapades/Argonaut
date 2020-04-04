package argo

import (
	"reflect"
)

func NewHitCounterArg(ptr interface{}) ArgumentBuilder {
	//if !isPointer(ptr) {
	//	ret
	//}
	//
	return nil
}

type hitCounter struct {
	value uint
	bind  interface{}
}

func (h *hitCounter) Hint() string {
	return ""
}

func (h *hitCounter) HasHint() bool {
	return false
}

func (h *hitCounter) Default() interface{} {
	return 0
}

func (h *hitCounter) HasDefault() bool {
	panic("implement me")
}

func (h *hitCounter) DefaultType() reflect.Type {
	panic("implement me")
}

func (h *hitCounter) Description() string {
	panic("implement me")
}

func (h *hitCounter) HasDescription() bool {
	panic("implement me")
}

func (h *hitCounter) RawValue() string {
	panic("implement me")
}

func (h *hitCounter) RequiredBool() {
	panic("implement me")
}

func (h *hitCounter) binding() interface{} {
	panic("implement me")
}

func (h *hitCounter) hasBinding() bool {
	panic("implement me")
}

func (h *hitCounter) parse(string) error {
	h.value++
	return nil
}
