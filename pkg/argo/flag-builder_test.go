package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
)

// no flag forms set
func TestFlagBuilder_Build01(t *testing.T) {
	_, err := cli.Flag().Build(nil)

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid short flag character
func TestFlagBuilder_Build02(t *testing.T) {
	_, err := cli.Flag().WithShortForm('-').Build(nil)

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid long flag name 1
func TestFlagBuilder_Build03(t *testing.T) {
	_, err := cli.Flag().WithLongForm("@@@@").Build(nil)

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid long flag name 1
func TestFlagBuilder_Build04(t *testing.T) {
	_, err := cli.Flag().WithLongForm("a@@@").Build(nil)

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// busted-ass argument
func TestFlagBuilder_Build05(t *testing.T) {
	_, err := cli.Flag().
		WithLongForm("test").
		WithBindingAndDefault(3, 4, true).
		Build(nil)

	if err == nil {
		t.Error("expected err to not have been nil, but it was")
	}
	t.Log(err)
}
