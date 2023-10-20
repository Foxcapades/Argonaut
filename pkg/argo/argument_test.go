package argo_test

import (
	"errors"
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
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
	str := ""
	con := argo.ConsumerFunc(func(val string) error { str = val; return nil })
	foo := func() (int, error) { return 3, nil }
	com := cli.Command().
		WithArgument(cli.Argument().WithBinding(con).WithDefault(foo)).
		MustParse([]string{"command", "poo"})

	arg := com.Arguments()[0]

	if !arg.WasHit() {
		t.Error("expected argument to have been hit but it wasn't")
	} else if str != "poo" {
		t.Error("expected bind value to equal the given function return value but it didn't")
	}
}
