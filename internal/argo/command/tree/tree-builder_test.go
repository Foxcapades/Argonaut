package tree_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TestEmptyCommandTree(t *testing.T) {
	opt := argo.Options{}
	_, err := tree.Build(tree.NewBuilder(), opt)
	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithDescription(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithDescription("hello").
		WithLeaf(tree.NewLeafBuilder("goodbye")), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"hello", "goodbye"}))

	if !com.HasDescription() {
		t.Fail()
	}

	if com.Description() != "hello" {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithCallback(t *testing.T) {
	counter := 0
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithCallback(func(com argo.TreeCommand) { counter++ }).
		WithLeaf(tree.NewLeafBuilder("goodbye")), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"hello", "goodbye"}))

	if counter != 1 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithFlag(t *testing.T) {
	counter := 0
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('c').WithBinding(&counter, true)).
		WithLeaf(tree.NewLeafBuilder("bell")), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"taco", "bell", "-c=3", "banana", "--", "pickle"}))

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

	if len(com.SelectedCommand().UnmappedInputs()) != 2 {
		t.Fail()
	}

	if com.SelectedCommand().UnmappedInputs()[0] != "banana" {
		t.Fail()
	}

	if com.SelectedCommand().UnmappedInputs()[1] != "pickle" {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithBranch(t *testing.T) {
	a := 0
	b := 0
	c := 0
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithCallback(func(com argo.TreeCommand) { a++ }).
		WithBranch(tree.NewBranchBuilder("foo").
			WithCallback(func(com argo.BranchCommand) { b++ }).
			WithLeaf(tree.NewLeafBuilder("bar").
				WithCallback(func(leaf argo.LeafCommand) { c++ }))), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"say", "foo", "bar"}))

	if a != 1 || b != 1 || c != 1 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithFlagGroup(t *testing.T) {
	value := 0
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("no-thanks")).
		WithFlagGroup(flag.NewGroupBuilder("derpy").
			WithFlag(flag.NewBuilder().
				WithShortForm('b').
				WithArgument(argument.NewBuilder().
					WithDefault(3).
					WithBinding(&value)))), opt))
	_ = utils.MustReturn(tree.Parse(com, []string{"hoopla", "no-thanks"}))

	if value != 3 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithLeaf(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("ass").
			WithFlag(flag.NewBuilder().
				WithLongForm("butt").
				WithArgument(argument.NewBuilder().Require()))), opt))

	_, err := tree.Parse(com, []string{"my", "ass", "--butt"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_Parse(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().
				WithShortForm('f').
				Require().
				WithArgument(argument.NewBuilder().Require()))), opt))

	_, err := tree.Parse(com, []string{"tree", "leaf"})

	if err == nil {
		t.Fail()
	}

	_, err = tree.Parse(com, []string{"tree", "leaf", "-f"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_InvalidArgumentBindingAndDefault(t *testing.T) {
	value := "hello"
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(3))), opt))

	_, err := tree.Parse(com, []string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	com = utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() int { return 3 }))), opt))

	_, err = tree.Parse(com, []string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	com = utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int) { return "hello", 3 }))), opt))

	_, err = tree.Parse(com, []string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	com = utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int, int) { return "hello", 3, 5 }))), opt))

	_, err = tree.Parse(com, []string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}
}

// Unrecognized short solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortSoloWarning(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")), opt))
	res := utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a"}))

	if len(res.Warnings) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if res.Warnings[0].Message != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized short pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortPairWarning(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithShortForm('b').WithArgument(argument.NewBuilder())), opt))

	res := utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab=1"}))

	if len(res.Warnings) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if res.Warnings[0].Message != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongSoloWarning(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")), argo.Options{}))

	res := utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--apple"}))

	if len(res.Warnings) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if res.Warnings[0].Message != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongPairWarning(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")), argo.Options{}))

	res := utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--apple=1"}))

	if len(res.Warnings) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if res.Warnings[0].Message != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Conflicting branch and leaf names in a single command group
func TestCommandTreeBuilder_Build01(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithBranch(tree.NewBranchBuilder("something").WithLeaf(tree.NewLeafBuilder("something-else"))).
		WithLeaf(tree.NewLeafBuilder("something")), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting branch and leaf names across command groups
func TestCommandTreeBuilder_Build02(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithCommandGroup(tree.NewGroupBuilder("foo").
			WithBranch(tree.NewBranchBuilder("something").WithLeaf(tree.NewLeafBuilder("something-else")))).
		WithCommandGroup(tree.NewGroupBuilder("bar").
			WithLeaf(tree.NewLeafBuilder("something"))), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names in a single flag group
func TestCommandTreeBuilder_Build03(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithFlag(flag.NewBuilder().WithLongForm("hello")).
		WithFlag(flag.NewBuilder().WithLongForm("hello")).
		WithLeaf(tree.NewLeafBuilder("something")), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting short flag names in a single flag group
func TestCommandTreeBuilder_Build04(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('f')).
		WithFlag(flag.NewBuilder().WithShortForm('f')).
		WithLeaf(tree.NewLeafBuilder("something")), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names across flag groups
func TestCommandTreeBuilder_Build05(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithFlagGroup(flag.NewGroupBuilder("hoopla").
			WithFlag(flag.NewBuilder().WithLongForm("hello"))).
		WithFlagGroup(flag.NewGroupBuilder("wednesday").
			WithFlag(flag.NewBuilder().WithLongForm("hello"))), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}

// Conflicting short flag names across flag groups
func TestCommandTreeBuilder_Build06(t *testing.T) {
	_, err := tree.Build(tree.NewBuilder().
		WithFlagGroup(flag.NewGroupBuilder("hoopla").
			WithFlag(flag.NewBuilder().WithShortForm('g'))).
		WithFlagGroup(flag.NewGroupBuilder("wednesday").
			WithFlag(flag.NewBuilder().WithShortForm('g'))), argo.Options{})

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}
