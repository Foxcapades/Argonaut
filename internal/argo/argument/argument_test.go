package argument_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
)

func TestArgument_Name(t *testing.T) {
	arg, err := argument.Build(argument.NewBuilder().WithName("name"))

	if err != nil {
		t.Error(err)
	}

	if !arg.HasName() {
		t.Error("expected argument to have name but it didn't")
	}

	if arg.Name() != "name" {
		t.Error("expected argument name to match configured name but it didn't")
	}
}

func TestArgument_Description(t *testing.T) {
	arg, err := argument.Build(argument.NewBuilder().WithDescription("description"))

	if err != nil {
		t.Error(err)
	}

	if !arg.HasDescription() {
		t.Error("expected argument to have description but it didn't")
	}

	if arg.Description() != "description" {
		t.Error("expected argument description to match configured description but it didn't")
	}
}

func TestArgument_Default(t *testing.T) {
	def := map[string]int8{"hello": -128}
	_, err := argument.Build(argument.NewBuilder().WithDefault(def))

	if err == nil {
		t.Error(err)
	}
}

// TestArgumentBuilder_Build01 ensures that the argument build will fail if the
// builder is passed an invalid 2 arg argument validator function due to that
// function not matching the type of the binding.
func TestArgumentBuilder_Build01(t *testing.T) {
	var bind int
	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&bind).
		WithValidator(func(foo float32, bar string) error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

// TestArgumentBuilder_Build02 ensures that the argument build will fail if the
// builder is passed an invalid 2 arg argument validator function due to that
// function not matching the type of the binding pointer.
func TestArgumentBuilder_Build02(t *testing.T) {
	foo := 0
	bar := &foo

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&bar).
		WithValidator(func(foo float32, bar string) error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

// TestArgumentBuilder_Build03 ensures that the argument build will fail if the
// builder is passed an invalid 2 arg argument validator function due to that
// function not matching the second param type "string".
func TestArgumentBuilder_Build03(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(foo int, bar int) error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

// TestArgumentBuilder_Build04 ensures that the argument build will fail if the
// builder is passed an invalid 2 arg argument validator function due to that
// function not returning a value of type `error`.
func TestArgumentBuilder_Build04(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(foo int, bar string) string { return "" }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build05(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(foo int, bar string) {}))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build06(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func() error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build07(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(int, string, int) error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build08(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(int) error { return nil }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build09(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(string) string { return "" }))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build10(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(string) {}))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build11(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(int, string) error { return nil }))

	if err != nil {
		t.Error("expected error to be nil but it wasn't")
	}
}

func TestArgumentBuilder_Build12(t *testing.T) {
	foo := 0

	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&foo).
		WithValidator(func(string) error { return nil }))

	if err != nil {
		t.Error("expected error to be nil but it wasn't")
	}
}

func TestArgumentBuilder_Build13(t *testing.T) {
	_, err := argument.Build(argument.NewBuilder().
		WithBinding(3).
		WithValidator(3))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}
}

func TestArgumentBuilder_Build14(t *testing.T) {
	var bind map[string]*[]string
	_, err := argument.Build(argument.NewBuilder().
		WithBinding(&bind))

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}
}
