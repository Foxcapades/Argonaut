package argo_test

import (
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
