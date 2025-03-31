package command_test

import (
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestCommandBuilder_Parse(t *testing.T) {
	com := utils.MustReturn(command.Build(command.NewBuilder()))
	utils.MustReturn(command.Parse(com, []string{"hello", "--foo", "bar"}))

	if !com.HasUnmappedInputs() {
		t.Error("command has no unmapped inputs")
	}

	if len(com.UnmappedInputs()) != 2 {
		t.Error("expected command to have 2 unmapped inputs but had", len(com.UnmappedInputs()))
	}

	if com.UnmappedInputs()[0] != "--foo" {
		t.Error("expected command unmapped input 1 to be '--foo' but was", com.UnmappedInputs()[0])
	}

	if com.UnmappedInputs()[1] != "bar" {
		t.Error("expected command unmapped input 2 to be 'var' but was", com.UnmappedInputs()[1])
	}
}

func TestCommandBuilder_WithArgument(t *testing.T) {
	var foo map[string]string

	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithArgument(cli.Argument().
			WithBinding(&foo))))
	utils.MustReturn(command.Parse(com, []string{"hello", "goober=banana"}))

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

	utils.MustReturn(command.Parse(
		utils.MustReturn(command.Build(cli.Command().
			WithUnmappedInputLabel("DUCKS...").
			WithFlag(cli.Flag().WithLongForm("value").WithBinding(&foo, true)))),
		[]string{
			"hello",
			"goodbye",
			"--value=flumps",
			"--value",
			"teddy",
		}))

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

func TestCommandBuilder_ConflictingLongFlags(t *testing.T) {
	_, err := command.Build(command.NewBuilder().
		WithFlag(cli.Flag().WithLongForm("hello")).
		WithFlagGroup(cli.FlagGroup("nope").
			WithFlag(cli.Flag().WithLongForm("hello"))))

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

func TestCommandBuilder_ConflictingShortFlags(t *testing.T) {
	_, err := command.Build(command.NewBuilder().
		WithFlag(cli.Flag().WithShortForm('a')).
		WithFlag(cli.Flag().WithShortForm('a')))

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

func TestCommandBuilder_ParseUnhitRequiredFlag(t *testing.T) {
	_, err := command.Parse(utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(cli.Flag().WithLongForm("apple").Require()).
		WithFlag(cli.Flag().WithShortForm('x')))), []string{"hello", "-x=banana", "--banana=orange"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandBuilder_OptionalArgumentBeforeRequiredArgument(t *testing.T) {
	_, err := command.Build(command.NewBuilder().
		WithArgument(cli.Argument()).
		WithArgument(cli.Argument().Require()))

	if !strings.Contains(err.(argo.MultiError).Errors()[0].Error(), "argument 1") {
		t.Error("expected a build error for argument 1")
	}
}
