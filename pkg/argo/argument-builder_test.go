package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
)

// TestArgumentBuilder_Build01 ensures that the argument build will fail if the
// builder is passed an invalid 2 arg argument validator function due to that
// function not matching the type of the binding.
func TestArgumentBuilder_Build01(t *testing.T) {
	var bind int
	_, err := cli.Argument().
		WithBinding(&bind).
		WithValidator(func(foo float32, bar string) error { return nil }).
		Build(nil)

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

	_, err := cli.Argument().
		WithBinding(&bar).
		WithValidator(func(foo float32, bar string) error { return nil }).
		Build(nil)

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

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(foo int, bar int) error { return nil }).
		Build(nil)

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

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(foo int, bar string) string { return "" }).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build05(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(foo int, bar string) {}).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build06(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func() error { return nil }).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build07(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(int, string, int) error { return nil }).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build08(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(int) error { return nil }).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build09(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(string) string { return "" }).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build10(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(string) {}).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}

	t.Log(err)
}

func TestArgumentBuilder_Build11(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(int, string) error { return nil }).
		Build(nil)

	if err != nil {
		t.Error("expected error to be nil but it wasn't")
	}
}

func TestArgumentBuilder_Build12(t *testing.T) {
	foo := 0

	_, err := cli.Argument().
		WithBinding(&foo).
		WithValidator(func(string) error { return nil }).
		Build(nil)

	if err != nil {
		t.Error("expected error to be nil but it wasn't")
	}
}

func TestArgumentBuilder_Build13(t *testing.T) {
	_, err := cli.Argument().
		WithBinding(3).
		WithValidator(3).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}
}

func TestArgumentBuilder_Build14(t *testing.T) {
	var bind map[string]*[]string
	_, err := cli.Argument().
		WithBinding(&bind).
		Build(nil)

	if err == nil {
		t.Error("expected error not to be nil but it was")
	}
}
