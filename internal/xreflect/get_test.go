package xreflect_test

import (
	"reflect"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/xreflect"
)

func TestRootType01(t *testing.T) {
	it := reflect.TypeOf("")
	ot := xreflect.RootType(it)

	if it != ot {
		t.Error("expected in type to equal out type but it didn't")
	}
}

func TestRootType02(t *testing.T) {
	it := reflect.TypeOf((*******string)(nil))
	ot := xreflect.RootType(it)

	if ot.String() != "string" {
		t.Error("expected out type to equal 'string' but it didn't")
	}
}

func TestGetFunctionParamType(t *testing.T) {
	pt := xreflect.GetFunctionParamType(reflect.TypeOf(TestRootType01), 0)
	if pt.String() != "*testing.T" {
		t.Error("expected out type was not matched")
	}
}
