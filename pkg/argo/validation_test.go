package argo_test

import (
	"os"
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

	os.Args = []string{"command", "3"}
	_, res, _ := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().
			WithBinding(&bind).
			WithValidator(argo.OneOfPostParseArgumentValidator([]int{1, 2}, "invalid value"))))

	if res.Error == nil {
		t.Error("expected error to not be nil but it was")
	} else if res.ErrorType != argo.InputError {
		t.Errorf("expected error type to be InputError but was %s", res.ErrorType)
	} else if res.Error.Error() != "invalid value" {
		t.Errorf("expected error text to match configured error message but it was: %s", res.Error)
	}
}

func TestNoneOfPreParseArgumentValidator(t *testing.T) {
	com := utils.MustReturn(command.Build(
		command.NewBuilder().
			WithArgument(cli.Argument().
				WithValidator(argo.NoneOfPreParseArgumentValidator([]string{"hello", "goodbye"}, "invalid value"))),
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
	))

	_, err := command.Parse(com, []string{"command", "2"})

	if err == nil {
		t.Error("expected error to not be nil but it was")
	} else if err.Error() != "invalid value" {
		t.Error("expected error text to match configured error message but it didn't")
		t.Log(err)
	}
}
