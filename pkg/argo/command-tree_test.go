package argo_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestEmptyCommandTree(t *testing.T) {
	_, err := argo.NewCommandTreeBuilder().Parse([]string{"hello"})
	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithDescription(t *testing.T) {
	cmd := argo.NewCommandTreeBuilder().
		WithDescription("hello").
		WithLeaf(argo.NewCommandLeafBuilder("goodbye")).
		MustParse([]string{"hello", "goodbye"})

	if !cmd.HasDescription() {
		t.Fail()
	}

	if cmd.Description() != "hello" {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithCallback(t *testing.T) {
	counter := 0
	argo.NewCommandTreeBuilder().
		WithCallback(func(com argo.CommandTree) { counter++ }).
		WithLeaf(argo.NewCommandLeafBuilder("goodbye")).
		MustParse([]string{"hello", "goodbye"})

	if counter != 1 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithFlag(t *testing.T) {
	counter := 0
	argo.NewCommandTreeBuilder().
		WithFlag(argo.NewFlagBuilder().WithShortForm('c').WithBinding(&counter, true)).
		WithLeaf(argo.NewCommandLeafBuilder("bell")).
		MustParse([]string{"taco", "bell", "-c=3"})

	if counter != 3 {
		t.Errorf("expected 3, got %d", counter)
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithBranch(t *testing.T) {
	a := 0
	b := 0
	c := 0
	argo.NewCommandTreeBuilder().
		WithCallback(func(com argo.CommandTree) { a++ }).
		WithBranch(argo.NewCommandBranchBuilder("foo").
			WithCallback(func(com argo.CommandBranch) { b++ }).
			WithLeaf(argo.NewCommandLeafBuilder("bar").
				WithCallback(func(leaf argo.CommandLeaf) { c++ }))).
		MustParse([]string{"say", "foo", "bar"})

	if a != 1 || b != 1 || c != 1 {
		t.Fail()
	}
}
