package argument_test

import (
	"reflect"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestDetermineBindKind01(t *testing.T) {
	kind, err := argument.DetermineBindType(0)

	if kind != argo.BindingTypeInvalid {
		t.Error("expected binding kind to be invalid but it wasn't")
	}

	if err == nil {
		t.Error("expected err to not be nil but it was")
	} else {
		t.Log(err)
	}
}

func TestDetermineBindKind02(t *testing.T) {
	kind, err := argument.DetermineBindType(nil)

	if kind != argo.BindingTypeInvalid {
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
	kind, err := argument.DetermineBindType(&foo)

	if kind != argo.BindingTypeInvalid {
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
	kind, err := argument.DetermineBindType(foo)

	if kind != argo.BindingTypeErrorFunc {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind05(t *testing.T) {
	var foo int
	kind, err := argument.DetermineBindType(&foo)

	if kind != argo.BindingTypePointer {
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
	kind, err := argument.DetermineBindType(&foo)

	if kind != argo.BindingTypeUnmarshaler {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind07(t *testing.T) {
	var foo = func(string) {}
	kind, err := argument.DetermineBindType(foo)

	if kind != argo.BindingTypeSimpleFunc {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind08(t *testing.T) {
	var foo = func(string) error { return nil }
	kind, err := argument.DetermineBindType(foo)

	if kind != argo.BindingTypeErrorFunc {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestDetermineBindKind09(t *testing.T) {
	var foo = func() error { return nil }
	kind, err := argument.DetermineBindType(foo)

	if kind != argo.BindingTypeInvalid {
		t.Error("expected binding kind to be valid but it wasn't")
	}

	if err != nil {
		t.Log(err)
	}
}

func TestResemblesProviderFunction01(t *testing.T) {
	if argument.ResemblesProviderFunction(reflect.TypeOf(func(string) {})) {
		t.Fail()
	}
}
