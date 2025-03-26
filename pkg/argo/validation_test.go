package argo_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestOneOfPreParseArgumentValidator(t *testing.T) {
	com := utils.MustReturn(command.Build(
		command.NewBuilder().
			WithArgument(cli.Argument().
				WithValidator(argo.OneOfPreParseArgumentValidator([]string{"hello", "goodbye"}, "invalid value"))),
		argo.DefaultOptions(),
	))

	_, err := command.Parse(com, []string{"command", "world"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
	}
}

func TestOneOfPostParseArgumentValidator(t *testing.T) {
	var bind int

	com := utils.MustReturn(command.Build(
		command.NewBuilder().
			WithArgument(cli.Argument().
				WithBinding(&bind).
				WithValidator(argo.OneOfPostParseArgumentValidator([]int{1, 2}, "invalid value"))),
		argo.DefaultOptions(),
	))

	_, err := command.Parse(com, []string{"command", "3"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
		t.Log(err)
	}
}

func TestNoneOfPreParseArgumentValidator(t *testing.T) {
	com := utils.MustReturn(command.Build(
		command.NewBuilder().
			WithArgument(cli.Argument().
				WithValidator(argo.NoneOfPreParseArgumentValidator([]string{"hello", "goodbye"}, "invalid value"))),
		argo.DefaultOptions(),
	))

	_, err := command.Parse(com, []string{"command", "hello"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
	}
}

func TestNoneOfPostParseArgumentValidator(t *testing.T) {
	var bind int

	com := utils.MustReturn(command.Build(
		command.NewBuilder().
			WithArgument(cli.Argument().
				WithBinding(&bind).
				WithValidator(argo.NoneOfPostParseArgumentValidator([]int{1, 2}, "invalid value"))),
		argo.DefaultOptions(),
	))

	_, err := command.Parse(com, []string{"command", "2"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
		t.Log(err)
	}
}
