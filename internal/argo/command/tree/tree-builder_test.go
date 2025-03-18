package tree_test

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
	cli_tree "github.com/foxcapades/argonaut/v3/pkg/argo/command/tree"
	"testing"
)

func TestEmptyCommandTree(t *testing.T) {
	_, err := tree.NewTreeBuilder().Parse([]string{"hello"})
	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithDescription(t *testing.T) {
	cmd := tree.NewTreeBuilder().
		WithDescription("hello").
		WithLeaf(tree.NewLeafBuilder("goodbye")).
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
	tree.NewTreeBuilder().
		WithCallback(func(com cli_tree.CommandTree) { counter++ }).
		WithLeaf(tree.NewLeafBuilder("goodbye")).
		MustParse([]string{"hello", "goodbye"})

	if counter != 1 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithFlag(t *testing.T) {
	counter := 0
	com := tree.NewTreeBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('c').WithBinding(&counter, true)).
		WithLeaf(tree.NewLeafBuilder("bell")).
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
	tree.NewTreeBuilder().
		WithCallback(func(com cli_tree.CommandTree) { a++ }).
		WithBranch(tree.NewBranchBuilder("foo").
			WithCallback(func(com cli_tree.CommandBranch) { b++ }).
			WithLeaf(tree.NewLeafBuilder("bar").
				WithCallback(func(leaf cli_tree.CommandLeaf) { c++ }))).
		MustParse([]string{"say", "foo", "bar"})

	if a != 1 || b != 1 || c != 1 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithFlagGroup(t *testing.T) {
	value := 0
	tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("no-thanks")).
		WithFlagGroup(flag.NewGroupBuilder("derpy").
			WithFlag(flag.NewBuilder().
				WithShortForm('b').
				WithArgument(argument.NewBuilder().
					WithDefault(3).
					WithBinding(&value)))).
		MustParse([]string{"hoopla", "no-thanks"})

	if value != 3 {
		t.Fail()
	}
}

func TestCommandTreeBuilder_WithLeaf(t *testing.T) {
	_, err := tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("ass").
			WithFlag(flag.NewBuilder().
				WithLongForm("butt").
				WithArgument(argument.NewBuilder().Require()))).
		Parse([]string{"my", "ass", "--butt"})

	if err == nil {
		t.Fail()
	}
}

func TestCommandTreeBuilder_Parse(t *testing.T) {
	builder := tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().
				WithShortForm('f').
				Require().
				WithArgument(argument.NewBuilder().Require())))

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
	_, err := tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(3))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() int { return 3 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int) { return "hello", 3 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}

	_, err = tree.NewTreeBuilder().
		WithLeaf(tree.NewLeafBuilder("bar").
			WithArgument(argument.NewBuilder().
				WithBinding(&value).
				WithDefault(func() (string, int, int) { return "hello", 3, 5 }))).
		Parse([]string{"foo", "bar"})

	if err == nil {
		t.Fail()
	}
}

// Unrecognized short solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortSoloWarning(t *testing.T) {
	com := argo.Tree().
		WithLeaf(argo.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-a"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized short pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownShortPairWarning(t *testing.T) {
	com := argo.Tree().
		WithLeaf(argo.Leaf("leaf")).
		WithFlag(argo.ShortFlag('b').WithArgument(argo.Argument())).
		MustParse([]string{"command", "leaf", "-ab=1"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized short flag -a" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long solo flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongSoloWarning(t *testing.T) {
	com := argo.Tree().
		WithLeaf(argo.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "--apple"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Unrecognized long pair flag becomes a warning.
func TestCommandTreeBuilder_UnknownLongPairWarning(t *testing.T) {
	com := argo.Tree().
		WithLeaf(argo.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "--apple=1"})

	if len(com.Warnings()) != 1 {
		t.Error("expected command tree to have exactly 1 error but it didn't")
	} else if com.Warnings()[0] != "unrecognized long flag --apple" {
		t.Error("expected command tree warning to match expected warning but it didn't")
	}
}

// Conflicting branch and leaf names in a single command group
func TestCommandTreeBuilder_Build01(t *testing.T) {
	_, err := argo.Tree().
		WithBranch(argo.Branch("something").WithLeaf(argo.Leaf("something-else"))).
		WithLeaf(argo.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting branch and leaf names across command groups
func TestCommandTreeBuilder_Build02(t *testing.T) {
	_, err := argo.Tree().
		WithCommandGroup(argo.CommandGroup("foo").
			WithBranch(argo.Branch("something").WithLeaf(argo.Leaf("something-else")))).
		WithCommandGroup(argo.CommandGroup("bar").
			WithLeaf(argo.Leaf("something"))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names in a single flag group
func TestCommandTreeBuilder_Build03(t *testing.T) {
	_, err := argo.Tree().
		WithFlag(argo.LongFlag("hello")).
		WithFlag(argo.LongFlag("hello")).
		WithLeaf(argo.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting short flag names in a single flag group
func TestCommandTreeBuilder_Build04(t *testing.T) {
	_, err := argo.Tree().
		WithFlag(argo.ShortFlag('f')).
		WithFlag(argo.ShortFlag('f')).
		WithLeaf(argo.Leaf("something")).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	}
}

// Conflicting long flag names across flag groups
func TestCommandTreeBuilder_Build05(t *testing.T) {
	_, err := argo.Tree().
		WithFlagGroup(argo.FlagGroup("hoopla").
			WithFlag(argo.LongFlag("hello"))).
		WithFlagGroup(argo.FlagGroup("wednesday").
			WithFlag(argo.LongFlag("hello"))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}

// Conflicting short flag names across flag groups
func TestCommandTreeBuilder_Build06(t *testing.T) {
	_, err := argo.Tree().
		WithFlagGroup(argo.FlagGroup("hoopla").
			WithFlag(argo.ShortFlag('g'))).
		WithFlagGroup(argo.FlagGroup("wednesday").
			WithFlag(argo.ShortFlag('g'))).
		Build(nil)

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else {
		t.Log(err)
	}
}
