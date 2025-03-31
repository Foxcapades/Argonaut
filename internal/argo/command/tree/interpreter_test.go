package tree_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argotype"
)

func TestInvalidSubCommand(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().WithLeaf(tree.NewLeafBuilder("leaf1"))))
	_, err := tree.Parse(com, []string{"command", "leaf2"})

	if err == nil {
		t.Error("expected error for invalid subcommand but got none")
	}
}

// expect flag, expect argument
func TestTreeInterpretLongPair01(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf1")).
		WithFlag(flag.NewBuilder().WithLongForm("foo").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--foo=bar", "leaf1"}))

	f := com.FindLongFlag("foo")

	if !f.WasHit() {
		t.Fatal("expected flag to be hit but it wasn't")
	}

	if !f.Argument().WasHit() {
		t.Fatal("expected argument to be hit but it wasn't")
	}

	if f.Argument().RawValue() != "bar" {
		t.Fatalf("expected argument value to be 'bar' but it was '%s'", f.Argument().RawValue())
	}
}

// Don't expect flag at all
func TestTreeInterpretLongPair02(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf1"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf1", "--foo=bar"}))

	sub := com.SelectedCommand()

	if !sub.HasUnmappedInputs() {
		t.Fatal("expected command to have unmapped inputs")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Fatal("expected command to have exactly 1 unmapped input")
	} else if sub.UnmappedInputs()[0] != "--foo=bar" {
		t.Fatalf("expected command unmapped input to be '--foo=bar' but was '%s'", sub.UnmappedInputs()[0])
	}
}

// Have flag, doesn't expect argument
func TestTreeInterpretLongPair03(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf1")).
		WithFlag(flag.NewBuilder().WithLongForm("foo"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--foo=bar", "leaf1"}))
	f := com.FindLongFlag("foo")

	if !f.WasHit() {
		t.Fatal("expected flag to be hit but it wasn't")
	}
}

// Unexpected solo long flag (goes to unmapped)
func TestTreeInterpretLongSolo01(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("hello"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "hello", "--hello"}))

	sub := com.SelectedCommand()

	if !sub.HasUnmappedInputs() {
		t.Fatal("expected command to have unmapped inputs")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Fatal("expected command to have exactly 1 unmapped input")
	} else if sub.UnmappedInputs()[0] != "--hello" {
		t.Fatalf("expected command unmapped input to be '--hello' but was '%s'", sub.UnmappedInputs()[0])
	}
}

// Solo flag requires argument but is followed by boundary
func TestTreeInterpretLongSolo02(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder().Require())))))

	_, err := tree.Parse(com, []string{"command", "leaf", "--flag", "--"})

	if err == nil {
		t.Errorf("expected a parse error, but none was returned")
	}
}

// Solo flag gets the required argument it craves so badly
func TestTreeInterpretLongSolo03(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder().Require())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "argument"}))

	f := com.FindLongFlagRecursive("flag")

	if !f.WasHit() {
		t.Fatal("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Fatal("expected argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "argument" {
		t.Fatalf("expected argument value to be 'argument' but it was '%s'", f.Argument().RawValue())
	}
}

// Solo flag gets a plain argument that it optionally accepts
func TestTreeInterpretLongSolo04(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--flag", "value", "leaf"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "value" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// Solo flag gets an argument that resembles a long flag, but isn't
// Solo flag gets a plain argument that it optionally accepts
func TestTreeInterpretLongSolo05(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--flag", "--not-a-flag", "leaf"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "--not-a-flag" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

func TestTreeInterpretLongSolo06(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	// TODO: flag is eating the next subcommand as its optional argument.  If the
	//       argument value was invalid, it would be used as the subcommand, but
	//       since the argument here accepts any value, the subcommand name is
	//       being eaten.
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--flag", "leaf"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("didn't expect argument to be hit, but it was")
	}
}

// Solo flag expects an optional argument but gets boundary
func TestTreeInterpretLongSolo07(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "--", "hoopla"}))

	sub := com.SelectedCommand()

	f := sub.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("didn't expect argument to be hit, but it was")
	}
}

// Solo flag expects an optional argument but gets a long flag pair that isn't
// registered.
func TestTreeInterpretLongSolo08(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "--flag", "--teddy=bear", "leaf"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "--teddy=bear" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// Solo flag expects an optional argument but gets a short flag that isn't
// registered.
func TestTreeInterpretLongSolo09(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "-g"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "-g" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// Solo flag expects an optional argument but gets a short flag pair that isn't
// registered.
func TestTreeInterpretLongSolo10(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "-p=eriod"}))

	f := com.FindLongFlag("flag")

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "-p=eriod" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// solo flag expects an optional argument but gets a short flag that _is_
// registered
func TestTreeInterpretLongSolo11(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithShortForm('g'))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "-g"}))

	longFlag := com.FindLongFlag("flag")
	shortFlag := com.FindShortFlag('g')

	if !longFlag.WasHit() {
		t.Error("expected long flag to be hit but it wasn't")
	} else if longFlag.Argument().WasHit() {
		t.Error("didn't expect long flag argument to be hit, but it was")
	}

	if !shortFlag.WasHit() {
		t.Error("expected short flag to be hit but it wasn't")
	}
}

// Solo flag expects an optional argument but gets a short flag pair that _is_
// registered
func TestTreeInterpretLongSolo12(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder())).
			WithFlag(flag.NewBuilder().WithShortForm('g').WithArgument(argument.NewBuilder())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "-g=randma"}))

	longFlag := com.FindLongFlag("flag")
	shortFlag := com.FindShortFlag('g')

	if !longFlag.WasHit() {
		t.Error("expected long flag to be hit but it wasn't")
	} else if longFlag.Argument().WasHit() {
		t.Error("didn't expect long flag argument to be hit, but it was")
	}

	if !shortFlag.WasHit() {
		t.Error("expected short flag to be hit but it wasn't")
	} else if !shortFlag.Argument().WasHit() {
		t.Error("expected short flag argument to be hit but it wasn't")
	} else if shortFlag.Argument().RawValue() != "randma" {
		t.Error("expected short flag argument value to match input, but it didn't")
	}
}

// Solo flag expects an optional argument but gets a long flag that _is_
// registered
func TestTreeInterpretLongSolo13(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder())).
			WithFlag(flag.NewBuilder().WithLongForm("other")))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "--other"}))

	flag1 := com.FindLongFlag("flag")
	flag2 := com.FindLongFlag("other")

	if !flag1.WasHit() {
		t.Error("expected long flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("didn't expect long flag 1 argument to be hit, but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected long flag 2 to be hit but it wasn't")
	}
}

// Solo flag expects an optional argument but gets a long flag pair that _is_
// registered
func TestTreeInterpretLongSolo14(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithLongForm("flag").WithArgument(argument.NewBuilder())).
			WithFlag(flag.NewBuilder().WithLongForm("other").WithArgument(argument.NewBuilder())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "--flag", "--other=thing"}))

	flag1 := com.FindLongFlag("flag")
	flag2 := com.FindLongFlag("other")

	if !flag1.WasHit() {
		t.Error("expected long flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("didn't expect long flag 1 argument to be hit, but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected long flag 2 to be hit but it wasn't")
	} else if !flag2.Argument().WasHit() {
		t.Error("expected long flag 2 argument to be hit but it wasn't")
	} else if flag2.Argument().RawValue() != "thing" {
		t.Error("expected long flag 2 argument value to match input, but it didn't")
	}
}

// Short pair takes optional argument that is plain text.
func TestTreeInterpretShortPair01(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(flag.NewBuilder().WithShortForm('c').WithArgument(argument.NewBuilder())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-c=for cookie"}))

	f := com.FindShortFlag('c')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if f.Argument().RawValue() != "for cookie" {
		t.Errorf("expected argument value to be 'for cookie' input but it was '%s'", f.Argument().RawValue())
	}
}

// Short pair is unregistered
func TestTreeInterpretShortPair02(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf"))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-c=for cookie"}))

	sub := com.SelectedCommand()

	if !sub.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if sub.UnmappedInputs()[0] != "-c=for cookie" {
		t.Error("expected unmapped value to match input but it didn't")
	}
}

// Short pair block is empty???
func TestTreeInterpretShortPair03(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf"))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-=for cookie"}))

	sub := com.SelectedCommand()

	if !sub.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if sub.UnmappedInputs()[0] != "-=for cookie" {
		t.Error("expected unmapped value to match input but it didn't")
	}
}

// multi-flag block with short pair, first flag requires argument
func TestTreeInterpretShortPair04(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('c').WithArgument(argument.NewBuilder().Require())).
		WithFlag(cli.ShortFlag('d').WithArgument(argument.NewBuilder().Require()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-cd=foo"}))

	flag1 := com.FindShortFlag('c')
	flag2 := com.FindShortFlag('d')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if !flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to be hit but it wasn't")
	} else if flag1.Argument().RawValue() != "d=foo" {
		t.Error("expected flag 1 argument to match input but it didn't")
	}

	if flag2.WasHit() {
		t.Error("expected flag 2 not to be hit but it was")
	}
}

// multi-flag block with a short pair, last flag requires argument
func TestTreeInterpretShortPair05(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())).
		WithFlag(cli.ShortFlag('b').WithArgument(argument.NewBuilder().Require()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab=foo"}))

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument not to be hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to be hit but it wasn't")
	} else if !flag2.Argument().WasHit() {
		t.Error("expected flag 2 argument to be hit but it wasn't")
	} else if flag2.Argument().RawValue() != "foo" {
		t.Error("expected flag 2 argument value to match input but it didn't")
	}
}

// multi-flag block with unknown first flag
func TestTreeInterpretShortPair06(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('b').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab=value"}))
	sub := com.SelectedCommand()

	f := com.FindShortFlag('b')

	if !sub.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if sub.UnmappedInputs()[0] != "-a" {
		t.Error("expected unmapped value to match input but it didn't")
	}

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "value" {
		t.Error("expected flag argument value to match input but it didn't")
	}
}

// multi-flag block with unknown middle flag
func TestTreeInterpretShortPair07(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('a')).
			WithFlag(cli.ShortFlag('c').WithArgument(argument.NewBuilder())))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-abc=value"}))

	sub := com.SelectedCommand()

	flag1 := com.FindShortFlagRecursive('a')
	flag2 := com.FindShortFlagRecursive('c')

	if !sub.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if sub.UnmappedInputs()[0] != "-b" {
		t.Error("expected unmapped value to match input but it didn't")
	}

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	} else if !flag2.Argument().WasHit() {
		t.Error("expected flag 2 argument to have been hit but it wasn't")
	} else if flag2.Argument().RawValue() != "value" {
		t.Error("expected flag 2 argument value to match input but it didn't")
	}
}

// flag that doesn't except an argument
func TestTreeInterpretShortPair08(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('a')).
			WithFlag(cli.ShortFlag('b')))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab=value"}))

	sub := com.SelectedCommand()

	flag1 := sub.FindShortFlagRecursive('a')
	flag2 := sub.FindShortFlagRecursive('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

// first flag accepts an argument, second flag is unrecognized
func TestTreeInterpretShortPair09(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab=value"}))

	sub := com.SelectedCommand()

	f := sub.FindShortFlagRecursive('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "b=value" {
		t.Error("expected flag argument to match input value but it doesn't")
	}
}

// short solo unknown flag
func TestTreeInterpreterShortSolo01(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf"))))

	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-c"}))

	sub := com.SelectedCommand()

	if !sub.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(sub.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if sub.UnmappedInputs()[0] != "-c" {
		t.Error("expected command unmapped input to match input but it didn't")
	}
}

// short solo requires arg, hits boundary
func TestTreeInterpreterShortSolo02(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('c').WithArgument(argument.NewBuilder().Require())))))

	_, err := tree.Parse(com, []string{"command", "leaf", "-c", "--"})

	if err == nil {
		t.Error("expected parsing to error out but it didn't")
	}
}

// solo short flag requires arg, eats rest of block
func TestTreeInterpreterShortSolo03(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder().Require())).
			WithFlag(cli.ShortFlag('b')))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-ab"}))

	flag1 := com.FindShortFlagRecursive('a')
	flag2 := com.FindShortFlagRecursive('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if !flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to have been hit but it wasn't")
	} else if flag1.Argument().RawValue() != "b" {
		t.Error("expected flag 1 argument to match the input but it didn't")
	}

	if flag2.WasHit() {
		t.Error("expected flag 2 not to have been hit but it was")
	}
}

// solo short flag expects optional argument but hits eof
func TestTreeInterpreterShortSolo04(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf").
			WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a"}))

	f := com.FindShortFlagRecursive('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

// solo short flag expects optional argument but hits boundary
func TestTreeInterpreterShortSolo05(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "--"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

// solo short flag expects optional and is followed by plain text
func TestTreeInterpreterShortSolo06(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "zoids"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "zoids" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known short flag
func TestTreeInterpreterShortSolo07(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())).
		WithFlag(cli.ShortFlag('b'))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "-b"}))

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to be hit but it wasn't")
	}
}

// solo short flag expects optional and is followed by an unknown short flag
func TestTreeInterpreterShortSolo08(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "-b"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "-b" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known short pair
func TestTreeInterpreterShortSolo09(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())).
		WithFlag(cli.ShortFlag('b'))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "-b=1"}))

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to be hit but it wasn't")
	}
}

// solo short flag expects optional and is followed by an unknown short flag
func TestTreeInterpreterShortSolo10(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "-b=1"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "-b=1" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known long flag
func TestTreeInterpreterShortSolo11(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())).
		WithFlag(cli.LongFlag("bacon"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "--bacon"}))

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindLongFlag("bacon")

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to be hit but it wasn't")
	}
}

// solo short flag expects optional and is followed by an unknown long flag
func TestTreeInterpreterShortSolo12(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "--beans"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "--beans" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known long pair
func TestTreeInterpreterShortSolo13(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder())).
		WithFlag(cli.LongFlag("bees"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "--bees=1"}))

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindLongFlag("bees")

	if !flag1.WasHit() {
		t.Error("expected flag 1 to be hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to be hit but it wasn't")
	}
}

// solo short flag expects optional and is followed by an unknown long pair
func TestTreeInterpreterShortSolo14(t *testing.T) {
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(argument.NewBuilder()))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a", "--bang=1"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "--bang=1" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/18
func TestRegression18Tree(t *testing.T) {
	bind := false
	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithFlag(cli.ShortFlag('a').WithBinding(&bind, false)).
		WithLeaf(tree.NewLeafBuilder("leaf"))))
	_ = utils.MustReturn(tree.Parse(com, []string{"command", "leaf", "-a"}))

	f := com.FindShortFlag('a')

	if !bind {
		t.Error("expected bind to be true, but it wasn't")
	}

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "true" {
		t.Error("expected flag argument value to be \"true\" but it wasn't")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/58
func TestRegression58Tree(t *testing.T) {
	var removeNAValues bool
	var inputsAreSorted bool
	var outputFormat uint8
	var printHeaders bool
	var inputFile string

	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("foo").
			WithFlag(cli.ComboFlag('r', "rm-na").
				WithBinding(&removeNAValues, false)).
			WithFlag(cli.ComboFlag('s', "sorted-inputs").
				WithBindingAndDefault(&inputsAreSorted, false, false)).
			WithFlag(cli.ComboFlag('f', "format").
				WithBindingAndDefault(func(val string) (err error) {
					switch strings.ToLower(val) {
					case "tsv":
						outputFormat = 1
					case "csv":
						outputFormat = 2
					case "json":
						outputFormat = 3
					case "jsonl":
						outputFormat = 4
					default:
						err = fmt.Errorf("unrecognized output format \"%s\"", val)
					}

					return
				}, "tsv", true)).
			WithFlag(cli.ComboFlag('t', "headers").
				WithBinding(&printHeaders, false)).
			WithArgument(argument.NewBuilder().
				WithName("file").
				WithBinding(func(path []string) (err error) {
					inputFile = path[0]
					return
				})))))
	_, err := tree.Parse(com, []string{"build/linux/find-bin-width", "foo", "-s", "-f", "tsv", "some-file"})

	if err != nil {
		t.Error("expected error to be nil, but was " + err.Error())
	}

	if removeNAValues {
		t.Error("expected removeNaValues to be false")
	}

	if !inputsAreSorted {
		t.Error("expected inputsAreSorted to be true")
	}

	if outputFormat != 1 {
		t.Error("expected outputFormat to be 1 but was " + strconv.Itoa(int(outputFormat)))
	}

	if printHeaders {
		t.Error("expected printHeaders to be false")
	}

	if inputFile != "some-file" {
		t.Error("expected input file to be some-file, but was '" + inputFile + "'")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/62
func TestRegression62CommandTree(t *testing.T) {
	var value argotype.Hex8

	com := utils.MustReturn(tree.Build(tree.NewBuilder().
		WithLeaf(tree.NewLeafBuilder("gen-meta").
			WithFlag(cli.ComboFlag('i', "interactive").
				WithDescription("Interactive mode: auto (0), none (1), minimal (2), full (3).  Defaults to auto").
				WithBindingAndDefault(&value, argotype.Hex8(23), true)))))

	_, err := tree.Parse(com, []string{"something", "gen-meta"})

	if err != nil {
		t.Error("expected err to be nil but was", err)
	}

	if value != 23 {
		t.Error("expected value to be 23 but was", value)
	}
}
