package argo_test

import (
	"fmt"
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestCommandBuilder_Parse(t *testing.T) {
	com := argo.NewCommandBuilder().
		MustParse([]string{"hello", "--foo", "bar"})

	if !com.HasUnmappedInputs() {
		t.Fail()
	}

	if len(com.UnmappedInputs()) != 2 {
		t.Fail()
	}

	if com.UnmappedInputs()[0] != "--foo" {
		t.Fail()
	}

	if com.UnmappedInputs()[1] != "bar" {
		t.Fail()
	}
}

func TestCommandBuilder_WithArgument(t *testing.T) {
	var foo map[string]string

	argo.NewCommandBuilder().
		WithArgument(argo.NewArgumentBuilder().
			WithBinding(&foo)).
		MustParse([]string{"hello", "goober=banana"})

	if len(foo) != 1 {
		t.Fail()
	}

	if value, ok := foo["goober"]; !ok {
		t.Fail()
	} else if value != "banana" {
		t.Fail()
	}
}

func TestCommandBuilder_WithUnmappedLabel(t *testing.T) {
	var foo []string
	argo.NewCommandBuilder().
		WithUnmappedLabel("DUCKS...").
		WithFlag(argo.NewFlagBuilder().WithLongForm("value").WithBinding(&foo, true)).
		MustParse([]string{
			"hello",
			"goodbye",
			"--value=flumps",
			"--value",
			"teddy",
		})

	if len(foo) != 2 {
		t.Fail()
	}

	if foo[0] != "flumps" {
		t.Fail()
	}

	if foo[1] != "teddy" {
		t.Fail()
	}

}

func TestCommandBuilder_ConflictingFlags(t *testing.T) {
	com, err := argo.NewCommandBuilder().
		WithFlag(argo.NewFlagBuilder().WithLongForm("hello")).
		WithFlagGroup(argo.NewFlagGroupBuilder("nope").
			WithFlag(argo.NewFlagBuilder().WithLongForm("hello"))).
		Parse([]string{"something"})

	if com != nil {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}

	fmt.Println(err)
}
