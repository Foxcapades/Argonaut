package argo_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestArgument_Name(t *testing.T) {
	arg, err := cli.Argument().WithName("name").Build()

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
	arg, err := cli.Argument().WithDescription("description").Build()

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
	_, err := cli.Argument().WithDefault(def).Build()

	if err == nil {
		t.Error(err)
	}
}

func TestArgumentDefault01(t *testing.T) {
	val := 0
	foo := func() int { return 3 }
	com := cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)).
		MustParse([]string{"command"})

	arg := com.Arguments()[0]

	if !arg.WasHit() {
		t.Error("expected argument to have been hit but it wasn't")
	} else if val != 3 {
		t.Error("expected bind value to equal the given function return value but it didn't")
	}
}

func TestArgumentDefault02(t *testing.T) {
	val := 0
	foo := func() (int, error) { return 3, nil }
	com := cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)).
		MustParse([]string{"command"})

	arg := com.Arguments()[0]

	if !arg.WasHit() {
		t.Error("expected argument to have been hit but it wasn't")
	} else if val != 3 {
		t.Error("expected bind value to equal the given function return value but it didn't")
	}
}

func TestArgumentDefault03(t *testing.T) {
	val := 0
	foo := func() (int, error) { return 0, errors.New("butt") }
	_, err := cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)).
		Parse([]string{"command"})

	if err == nil {
		t.Error("expected parsing to error but it didn't")
	}
}

func TestArgumentDefault04(t *testing.T) {
	con := argo.UnmarshalerFunc(func(val string) error { return nil })
	foo := func() (int, error) { return 3, nil }
	_, err := cli.Command().
		WithArgument(cli.Argument().WithBinding(con).WithDefault(foo)).
		Parse([]string{"command", "poo"})

	t.Log(err)
	if err == nil {
		t.Fail()
	}
}

type argUnmarshaler struct{ val int }

func (a *argUnmarshaler) Unmarshal(raw string) (err error) { a.val, err = strconv.Atoi(raw); return }

// Expect a failure because the unmarshaler type is incompatible with the output
// of the default value provider.
func TestArgumentDefault05(t *testing.T) {
	con := argUnmarshaler{}
	foo := func() (int, error) { return 3, nil }
	_, err := cli.Command().
		WithArgument(cli.Argument().WithBinding(&con).WithDefault(foo)).
		Parse([]string{"command", "poo"})

	t.Log(err)
	if err == nil {
		t.Fail()
	}
}

// Expect an OK because the unmarshaler type is compatible with the output of
// the default value provider.
// TODO: What should this test actually be?  It doesn't make sense to have the
//       provider return an unmarshaler instance.
// func TestArgumentDefault06(t *testing.T) {
// 	con := argUnmarshaler{}
// 	foo := func() (argUnmarshaler, error) { return argUnmarshaler{3}, nil }
// 	_ = cli.Command().
// 		WithArgument(cli.Argument().WithBinding(&con).WithDefault(foo)).
// 		MustParse([]string{"command", "3"})
//
// 	if con.val != 3 {
// 		t.Error("expected unmarshaler value to have been replaced but it wasn't")
// 	}
// }

// Expect an OK because the unmarshaler type is compatible with the output of
// the default value provider.
// TODO: What should this test actually be?  It doesn't make sense to have the
//       provider return an unmarshaler instance.
// func TestArgumentDefault07(t *testing.T) {
// 	con := argUnmarshaler{}
// 	foo := func() (argUnmarshaler, error) { return argUnmarshaler{3}, nil }
// 	_ = cli.Command().
// 		WithArgument(cli.Argument().WithBinding(&con).WithDefault(foo)).
// 		MustParse([]string{"command"})
//
// 	if con.val != 3 {
// 		t.Error("expected unmarshaler value to have been replaced but it wasn't")
// 	}
// }

func TestArgument_PreParseValidator01(t *testing.T) {
	var binding int

	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(string) error { return errors.New("dummy error") }).
			WithBinding(&binding)).
		Parse([]string{"command", "32"})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else if err.Error() != "dummy error" {
		t.Error("expected err to match validator output but it didn't")
	}
	t.Log(err)
}

func TestArgument_PreParseValidator02(t *testing.T) {
	var binding int

	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(string) error { return nil }).
			WithBinding(&binding)).
		Parse([]string{"command", "32"})

	if err != nil {
		t.Error("expected err to be nil but it wasn't")
	}
	t.Log(err)
}

func TestArgument_PostParseValidator01(t *testing.T) {
	var binding int

	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(int, string) error { return errors.New("dummy error") }).
			WithBinding(&binding)).
		Parse([]string{"command", "32"})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else if err.Error() != "dummy error" {
		t.Error("expected err to match validator output but it didn't")
	}
	t.Log(err)
}

func TestArgument_PostParseValidator02(t *testing.T) {
	var binding int

	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(int, string) error { return nil }).
			WithBinding(&binding)).
		Parse([]string{"command", "32"})

	if err != nil {
		t.Error("expected err to be nil but it wasn't")
	}
	t.Log(err)
}

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
