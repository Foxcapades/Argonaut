package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestOneOfPreParseArgumentValidator(t *testing.T) {
	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(argo.OneOfPreParseArgumentValidator([]string{"hello", "goodbye"}, "invalid value"))).
		Parse([]string{"command", "world"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
	}
}

func TestOneOfPostParseArgumentValidator(t *testing.T) {
	var bind int
	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithBinding(&bind).
			WithValidator(argo.OneOfPostParseArgumentValidator([]int{1, 2}, "invalid value"))).
		Parse([]string{"command", "3"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
		t.Log(err)
	}
}

func TestNoneOfPreParseArgumentValidator(t *testing.T) {
	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithValidator(argo.NoneOfPreParseArgumentValidator([]string{"hello", "goodbye"}, "invalid value"))).
		Parse([]string{"command", "hello"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
	}
}

func TestNoneOfPostParseArgumentValidator(t *testing.T) {
	var bind int
	_, err := cli.Command().
		WithArgument(cli.Argument().
			WithBinding(&bind).
			WithValidator(argo.NoneOfPostParseArgumentValidator([]int{1, 2}, "invalid value"))).
		Parse([]string{"command", "2"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
		t.Log(err)
	}
}
