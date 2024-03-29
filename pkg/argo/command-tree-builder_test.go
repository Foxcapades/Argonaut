package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
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
	com := argo.NewCommandTreeBuilder().
		WithFlag(argo.NewFlagBuilder().WithShortForm('c').WithBinding(&counter, true)).
		WithLeaf(argo.NewCommandLeafBuilder("bell")).
		MustParse([]string{"taco", "bell", "-c=3", "banana", "--", "pickle"})

	if counter != 3 {
		t.Errorf("expected 3, got %d", counter)
		t.Fail()
	}

	if com.SelectedCommand().FindShortFlag('c') == nil {
		t.Fail()
	}

	if !com.SelectedCommand().HasUnmappedInputs() {
		t.Fail()
	}

	if len(com.SelectedCommand().UnmappedInputs()) != 1 {
		t.Fail()
	}

	if com.SelectedCommand().UnmappedInputs()[0] != "banana" {
		t.Fail()
	}

	if !com.SelectedCommand().HasPassthroughInputs() {
		t.Fail()
	}

	if len(com.SelectedCommand().PassthroughInputs()) != 1 {
		t.Fail()
	}

	if com.SelectedCommand().PassthroughInputs()[0] != "pickle" {
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

func TestCommandTreeBuilder_WithFlagGroup(t *testing.T) {
	value := 0
	argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("no-thanks")).
		WithFlagGroup(argo.NewFlagGroupBuilder("derpy").
			WithFlag(argo.NewFlagBuilder().
				WithShortForm('b').
				WithArgument(argo.NewArgumentBuilder().
					WithDefault(3).
					WithBinding(&value)))).
		MustParse([]string{"hoopla", "no-thanks"})

	if value != 3 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithLeaf(t *testing.T) {
	_, err := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("ass").
			WithFlag(argo.NewFlagBuilder().
				WithLongForm("butt").
				WithArgument(argo.NewArgumentBuilder().Require()))).
		Parse([]string{"my", "ass", "--butt"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_Parse(t *testing.T) {
	builder := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("leaf").
			WithFlag(argo.NewFlagBuilder().
				WithShortForm('f').
				Require().
				WithArgument(argo.NewArgumentBuilder().Require())))

	_, err := builder.Parse([]string{"tree", "leaf"})

	if err == nil {
		t.Fail()
	}

	_, err = builder.Parse([]string{"tree", "leaf", "-f"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_InvalidArgumentBindingAndDefault(t *testing.T) {
	value := "hello"
	_, err := argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("bar").
			WithArgument(argo.NewArgumentBuilder().
				WithBinding(&value).
				WithDefault(3))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("bar").
			WithArgument(argo.NewArgumentBuilder().
				WithBinding(&value).
				WithDefault(func() int { return 3 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("bar").
			WithArgument(argo.NewArgumentBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int) { return "hello", 3 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = argo.NewCommandTreeBuilder().
		WithLeaf(argo.NewCommandLeafBuilder("bar").
			WithArgument(argo.NewArgumentBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int, int) { return "hello", 3, 5 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}
}

// Unrecognized short solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortSoloWarning(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-a"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized short pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortPairWarning(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('b').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-ab=1"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongSoloWarning(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "--apple"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongPairWarning(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "--apple=1"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Conflicting branch and leaf names in a single command group
func TestCommandTreeBuilder_Build01(t *testing.T) {
	_, err := cli.Tree().
		WithBranch(cli.Branch("something").WithLeaf(cli.Leaf("something-else"))).
		WithLeaf(cli.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting branch and leaf names across command groups
func TestCommandTreeBuilder_Build02(t *testing.T) {
	_, err := cli.Tree().
		WithCommandGroup(cli.CommandGroup("foo").
			WithBranch(cli.Branch("something").WithLeaf(cli.Leaf("something-else")))).
		WithCommandGroup(cli.CommandGroup("bar").
			WithLeaf(cli.Leaf("something"))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names in a single flag group
func TestCommandTreeBuilder_Build03(t *testing.T) {
	_, err := cli.Tree().
		WithFlag(cli.LongFlag("hello")).
		WithFlag(cli.LongFlag("hello")).
		WithLeaf(cli.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting short flag names in a single flag group
func TestCommandTreeBuilder_Build04(t *testing.T) {
	_, err := cli.Tree().
		WithFlag(cli.ShortFlag('f')).
		WithFlag(cli.ShortFlag('f')).
		WithLeaf(cli.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names across flag groups
func TestCommandTreeBuilder_Build05(t *testing.T) {
	_, err := cli.Tree().
		WithFlagGroup(cli.FlagGroup("hoopla").
			WithFlag(cli.LongFlag("hello"))).
		WithFlagGroup(cli.FlagGroup("wednesday").
			WithFlag(cli.LongFlag("hello"))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}

// Conflicting short flag names across flag groups
func TestCommandTreeBuilder_Build06(t *testing.T) {
	_, err := cli.Tree().
		WithFlagGroup(cli.FlagGroup("hoopla").
			WithFlag(cli.ShortFlag('g'))).
		WithFlagGroup(cli.FlagGroup("wednesday").
			WithFlag(cli.ShortFlag('g'))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}
