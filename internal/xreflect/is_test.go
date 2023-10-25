package xreflect_test

import (
	"reflect"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/xreflect"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestIsPointer(t *testing.T) {
	if !xreflect.IsPointer(reflect.TypeOf((*string)(nil))) {
		t.Fail()
	}
}

func TestIsString(t *testing.T) {
	if !xreflect.IsString(reflect.TypeOf("")) {
		t.Fail()
	}
}

func TestIsNumeric01(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(float32(0.0))) {
		t.Fail()
	}
}

func TestIsNumeric02(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(0.0)) {
		t.Fail()
	}
}

func TestIsNumeric03(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(0)) {
		t.Fail()
	}
}

func TestIsNumeric04(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(uint(0))) {
		t.Fail()
	}
}

func TestIsNumeric05(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(int8(0))) {
		t.Fail()
	}
}

func TestIsNumeric06(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(int16(0))) {
		t.Fail()
	}
}

func TestIsNumeric07(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(int32(0))) {
		t.Fail()
	}
}

func TestIsNumeric08(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(int64(0))) {
		t.Fail()
	}
}

func TestIsNumeric09(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(uint8(0))) {
		t.Fail()
	}
}

func TestIsNumeric10(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(uint16(0))) {
		t.Fail()
	}
}

func TestIsNumeric11(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(uint32(0))) {
		t.Fail()
	}
}

func TestIsNumeric12(t *testing.T) {
	if !xreflect.IsNumeric(reflect.TypeOf(uint64(0))) {
		t.Fail()
	}
}

func TestIsUnmarshaler(t *testing.T) {
	var in nmrshlr
	if !xreflect.IsUnmarshaler(reflect.TypeOf(in), reflect.TypeOf((*argo.Unmarshaler)(nil)).Elem()) {
		t.Fail()
	}
}

type nmrshlr struct{}

func (nmrshlr) Unmarshal(string) error { return nil }

func TestIsInterface(t *testing.T) {
	var foo any
	if !xreflect.IsInterface(reflect.TypeOf(&foo).Elem()) {
		t.Fail()
	}
}

func TestIsFunction(t *testing.T) {
	if !xreflect.IsFunction(reflect.TypeOf(TestIsInterface)) {
		t.Fail()
	}
}

func TestIsSlice(t *testing.T) {
	var foo []string
	if !xreflect.IsSlice(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsBasic01(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(float32(0))) {
		t.Fail()
	}
}

func TestIsBasic02(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(float64(0))) {
		t.Fail()
	}
}

func TestIsBasic03(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(0)) {
		t.Fail()
	}
}

func TestIsBasic04(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(uint(0))) {
		t.Fail()
	}
}

func TestIsBasic05(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(int8(0))) {
		t.Fail()
	}
}

func TestIsBasic06(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(int16(0))) {
		t.Fail()
	}
}

func TestIsBasic07(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(int32(0))) {
		t.Fail()
	}
}

func TestIsBasic08(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(int64(0))) {
		t.Fail()
	}
}

func TestIsBasic09(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(uint8(0))) {
		t.Fail()
	}
}

func TestIsBasic10(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(uint16(0))) {
		t.Fail()
	}
}

func TestIsBasic11(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(uint32(0))) {
		t.Fail()
	}
}

func TestIsBasic12(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(uint64(0))) {
		t.Fail()
	}
}

func TestIsBasic13(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf(true)) {
		t.Fail()
	}
}

func TestIsBasic14(t *testing.T) {
	if !xreflect.IsBasic(reflect.TypeOf("")) {
		t.Fail()
	}
}

func TestIsByteSlice(t *testing.T) {
	if !xreflect.IsByteSlice(reflect.TypeOf([]byte{})) {
		t.Fail()
	}
}

func TestIsBasicSlice(t *testing.T) {
	if !xreflect.IsBasicSlice(reflect.TypeOf([]string{})) {
		t.Fail()
	}
}

func TestFuncHasReturn01(t *testing.T) {
	if !xreflect.FuncHasReturn(reflect.TypeOf(func() int { return 0 })) {
		t.Fail()
	}
}

func TestFuncHasReturn02(t *testing.T) {
	if xreflect.FuncHasReturn(reflect.TypeOf(func() {})) {
		t.Fail()
	}
}

func TestIsBasicPointer01(t *testing.T) {
	var foo int
	if !xreflect.IsBasicPointer(reflect.TypeOf(&foo)) {
		t.Fail()
	}
}

func TestIsBasicPointer02(t *testing.T) {
	var foo int
	if xreflect.IsBasicPointer(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsBasicMap01(t *testing.T) {
	if xreflect.IsBasicMap(reflect.TypeOf("")) {
		t.Fail()
	}
}

func TestIsBasicMap02(t *testing.T) {
	var foo map[complex64]string
	if xreflect.IsBasicMap(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsBasicMap03(t *testing.T) {
	var foo map[string]complex64
	if xreflect.IsBasicMap(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsBasicMap04(t *testing.T) {
	var foo map[string]string
	if !xreflect.IsBasicMap(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsBasicMap05(t *testing.T) {
	var foo map[string][]string
	if !xreflect.IsBasicMap(reflect.TypeOf(foo)) {
		t.Fail()
	}
}

func TestIsNil_withNilIntPointer(t *testing.T) {
	var foo *int
	var val = reflect.ValueOf(foo)
	if !xreflect.IsNil(&val) {
		t.Fail()
	}
}

func TestIsNil_withNonNilIntPointer(t *testing.T) {
	var foo int
	var val = reflect.ValueOf(&foo)
	if xreflect.IsNil(&val) {
		t.Fail()
	}
}

func TestIsNil_withInt(t *testing.T) {
	var foo int
	var val = reflect.ValueOf(foo)
	if xreflect.IsNil(&val) {
		t.Fail()
	}
}
