package xarg_test

import (
	"reflect"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/xarg"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

var unmarshalerType = reflect.TypeOf((*argo.Unmarshaler)(nil)).Elem()

func TestDetermineBindKind01(t *testing.T) {
	kind, err := xarg.DetermineBindKind(0, unmarshalerType)

	if kind != xarg.BindKindInvalid {
		t.Error("expected binding kind to be invalid but it wasn't")
	}

	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestDetermineBindKind02(t *testing.T) {
	kind, err := xarg.DetermineBindKind(nil, unmarshalerType)

	if kind != xarg.BindKindInvalid {
		t.Error("expected binding kind to be invalid but it wasn't")
	}

	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestDetermineBindKind03(t *testing.T) {
	var foo error
	kind, err := xarg.DetermineBindKind(&foo, unmarshalerType)

	if kind != xarg.BindKindInvalid {
		t.Error("expected binding kind to be invalid but it wasn't")
	}

	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestDetermineBindKind04(t *testing.T) {
	var foo = func(string) error { return nil }
	kind, err := xarg.DetermineBindKind(foo, unmarshalerType)

	if kind != xarg.BindKindFuncWithErr {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind05(t *testing.T) {
	var foo int
	kind, err := xarg.DetermineBindKind(&foo, unmarshalerType)

	if kind != xarg.BindKindPointer {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

type nmrshlr struct{}

func (*nmrshlr) Unmarshal(string) error { return nil }

func TestDetermineBindKind06(t *testing.T) {
	var foo argo.Unmarshaler = &nmrshlr{}
	kind, err := xarg.DetermineBindKind(&foo, unmarshalerType)

	if kind != xarg.BindKindUnmarshaler {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind07(t *testing.T) {
	var foo = func(string) {}
	kind, err := xarg.DetermineBindKind(foo, unmarshalerType)

	if kind != xarg.BindKindFuncPlain {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind08(t *testing.T) {
	var foo = func(string) error { return nil }
	kind, err := xarg.DetermineBindKind(foo, unmarshalerType)

	if kind != xarg.BindKindFuncWithErr {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind09(t *testing.T) {
	var foo = func() error { return nil }
	kind, err := xarg.DetermineBindKind(foo, unmarshalerType)

	if kind != xarg.BindKindInvalid {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind01(t *testing.T) {
	var bind int
	var def int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind02(t *testing.T) {
	var bind int
	var def string

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindRaw {
		t.Error("expected default kind to be raw but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind03(t *testing.T) {
	var bind int
	var def float32

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind04(t *testing.T) {
	var bind int
	def := func() int { return 0 }

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindProviderPlain {
		t.Error("expected default kind to be plain provider but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind05(t *testing.T) {
	var bind int
	def := func() (int, error) { return 0, nil }

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindProviderWithErr {
		t.Error("expected default kind to be provider with error but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind06(t *testing.T) {
	var bind int
	def := func() string { return "" }

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind07(t *testing.T) {
	var bind int
	def := func() (int, string) { return 0, "" }

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind08(t *testing.T) {
	var bind int
	def := func() {}

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind09(t *testing.T) {
	var bind []int
	var def []int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind10(t *testing.T) {
	var bind int
	var def []int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind11(t *testing.T) {
	var bind []string
	var def []int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestDetermineDefaultKind12(t *testing.T) {
	bind := func([]int) {}
	var def []int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind13(t *testing.T) {
	var bind map[string]string
	var def map[string]string

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindParsed {
		t.Error("expected default kind to be parsed but it wasn't")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDetermineDefaultKind14(t *testing.T) {
	var bind int
	var def map[int]int

	kind, err := xarg.DetermineDefaultKind(bind, def)

	if kind != xarg.DefaultKindInvalid {
		t.Error("expected default kind to be invalid but it wasn't")
	}
	if err != nil {
		t.Log(err)
	}
}

func TestResemblesProviderFunction01(t *testing.T) {
	if xarg.ResemblesProviderFunction(reflect.TypeOf(func(string) {})) {
		t.Fail()
	}
}
