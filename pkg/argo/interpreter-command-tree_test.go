package argo_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TestInvalidSubCommand(t *testing.T) {
	_, err := cli.Tree().
		WithLeaf(cli.Leaf("leaf1")).
		Parse([]string{"command", "leaf2"})

	if err == nil {
		t.Error(err)
	}
}

// expect flag, expect argument
func TestTreeInterpretLongPair01(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf1")).
		WithFlag(cli.Flag().WithLongForm("foo").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf1", "--foo=bar"})

	flag := com.FindLongFlag("foo")

	if !flag.WasHit() {
		t.Fail()
	}

	if !flag.Argument().WasHit() {
		t.Fail()
	}

	if flag.Argument().RawValue() != "bar" {
		t.Fail()
	}
}

// Don't expect flag at all
func TestTreeInterpretLongPair02(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf1")).
		MustParse([]string{"command", "leaf1", "--foo=bar"}).
		SelectedCommand()

	if !com.HasUnmappedInputs() {
		t.Fail()
	} else if len(com.UnmappedInputs()) != 1 {
		t.Fail()
	} else if com.UnmappedInputs()[0] != "--foo=bar" {
		t.Fail()
	}
}

// Have flag, doesn't expect argument
func TestTreeInterpretLongPair03(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf1")).
		WithFlag(cli.Flag().WithLongForm("foo")).
		MustParse([]string{"command", "leaf1", "--foo=bar"})
	flag := com.FindLongFlag("foo")

	if !flag.WasHit() {
		t.Fail()
	}
}

// Unexpected solo long flag (goes to unmapped)
func TestTreeInterpretLongSolo01(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("hello")).
		MustParse([]string{"command", "hello", "--hello"}).
		SelectedCommand()

	if !com.HasUnmappedInputs() {
		t.Fail()
	} else if len(com.UnmappedInputs()) != 1 {
		t.Fail()
	} else if com.UnmappedInputs()[0] != "--hello" {
		t.Fail()
	}
}

// Solo flag requires argument but is followed by boundary
func TestTreeInterpretLongSolo02(t *testing.T) {
	_, err := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument().Require())).
		Parse([]string{"command", "leaf", "--flag", "--"})

	if err == nil {
		t.Fail()
	}
}

// Solo flag gets the required argument it craves so badly
func TestTreeInterpretLongSolo03(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument().Require())).
		MustParse([]string{"command", "leaf", "--flag", "argument"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Fail()
	} else if !flag.Argument().WasHit() {
		t.Fail()
	} else if flag.Argument().RawValue() != "argument" {
		t.Fail()
	}
}

// Solo flag gets a plain argument that it optionally accepts
func TestTreeInterpretLongSolo04(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "value"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "value" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// Solo flag gets an argument that resembles a long flag, but isn't
// Solo flag gets a plain argument that it optionally accepts
func TestTreeInterpretLongSolo05(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "--not-a-flag"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "--not-a-flag" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// Solo flag expects an optional argument but gets end of input
func TestTreeInterpretLongSolo06(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("didn't expect argument to be hit, but it was")
	}
}

// Solo flag expects an optional argument but gets boundary
func TestTreeInterpretLongSolo07(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "--", "hoopla"}).
		SelectedCommand()

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("didn't expect argument to be hit, but it was")
	}

	if !com.HasPassthroughInputs() {
		t.Error("expected command to have passthrough inputs but it didn't")
	} else if len(com.PassthroughInputs()) != 1 {
		t.Error("expected command to have 1 passthrough input but it didn't")
	} else if com.PassthroughInputs()[0] != "hoopla" {
		t.Error("expected command passthrough value to match input but it didn't")
	}
}

// Solo flag expects an optional argument but gets a long flag pair that isn't
// registered.
func TestTreeInterpretLongSolo08(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "--teddy=bear"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "--teddy=bear" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// Solo flag expects an optional argument but gets a short flag that isn't
// registered.
func TestTreeInterpretLongSolo09(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "-g"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "-g" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// Solo flag expects an optional argument but gets a short flag pair that isn't
// registered.
func TestTreeInterpretLongSolo10(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "-p=eriod"})

	flag := com.FindLongFlag("flag")

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "-p=eriod" {
		t.Errorf("expected argument value to match input, but it didn't")
	}
}

// solo flag expects an optional argument but gets a short flag that _is_
// registered
func TestTreeInterpretLongSolo11(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		WithFlag(cli.Flag().WithShortForm('g')).
		MustParse([]string{"command", "leaf", "--flag", "-g"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		WithFlag(cli.Flag().WithShortForm('g').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "-g=randma"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		WithFlag(cli.Flag().WithLongForm("other")).
		MustParse([]string{"command", "leaf", "--flag", "--other"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithLongForm("flag").WithArgument(cli.Argument())).
		WithFlag(cli.Flag().WithLongForm("other").WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "--flag", "--other=thing"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.Flag().WithShortForm('c').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-c=for cookie"})

	flag := com.FindShortFlag('c')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to be hit but it wasn't")
	} else if flag.Argument().RawValue() != "for cookie" {
		t.Log(flag.Argument().RawValue())
		t.Error("expected argument value to match input but it didn't")
	}
}

// Short pair is unregistered
func TestTreeInterpretShortPair02(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-c=for cookie"}).
		SelectedCommand()

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-c=for cookie" {
		t.Error("expected unmapped value to match input but it didn't")
	}
}

// Short pair block is empty???
func TestTreeInterpretShortPair03(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-=for cookie"}).
		SelectedCommand()

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-=for cookie" {
		t.Error("expected unmapped value to match input but it didn't")
	}
}

// multi-flag block with short pair, first flag requires argument
func TestTreeInterpretShortPair04(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('c').WithArgument(cli.Argument().Require())).
		WithFlag(cli.ShortFlag('d').WithArgument(cli.Argument().Require())).
		MustParse([]string{"command", "leaf", "-cd=foo"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('b').WithArgument(cli.Argument().Require())).
		MustParse([]string{"command", "leaf", "-ab=foo"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('b').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-ab=value"}).
		SelectedCommand()

	flag := com.FindShortFlag('b')

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-a" {
		t.Error("expected unmapped value to match input but it didn't")
	}

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "value" {
		t.Error("expected flag argument value to match input but it didn't")
	}
}

// multi-flag block with unknown middle flag
func TestTreeInterpretShortPair07(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a')).
		WithFlag(cli.ShortFlag('c').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-abc=value"}).
		SelectedCommand()

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('c')

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-b" {
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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a')).
		WithFlag(cli.ShortFlag('b')).
		MustParse([]string{"command", "leaf", "-ab=value"}).
		SelectedCommand()

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

// first flag accepts an argument, second flag is unrecognized
func TestTreeInterpretShortPair09(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-ab=value"}).
		SelectedCommand()

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "b=value" {
		t.Error("expected flag argument to match input value but it doesn't")
	}
}

// short solo unknown flag
func TestTreeInterpreterShortSolo01(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-c"}).
		SelectedCommand()

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-c" {
		t.Error("expected command unmapped input to match input but it didn't")
	}
}

// short solo requires arg, hits boundary
func TestTreeInterpreterShortSolo02(t *testing.T) {
	_, err := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('c').WithArgument(cli.Argument().Require())).
		Parse([]string{"command", "leaf", "-c", "--"})

	if err == nil {
		t.Error("expected parsing to error out but it didn't")
	}
}

// solo short flag requires arg, eats rest of block
func TestTreeInterpreterShortSolo03(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument().Require())).
		WithFlag(cli.ShortFlag('b')).
		MustParse([]string{"command", "leaf", "-ab"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

// solo short flag expects optional argument but hits boundary
func TestTreeInterpreterShortSolo05(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "--"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

// solo short flag expects optional and is followed by plain text
func TestTreeInterpreterShortSolo06(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "zoids"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "zoids" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known short flag
func TestTreeInterpreterShortSolo07(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('b')).
		MustParse([]string{"command", "leaf", "-a", "-b"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "-b"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "-b" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known short pair
func TestTreeInterpreterShortSolo09(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('b')).
		MustParse([]string{"command", "leaf", "-a", "-b=1"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "-b=1"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "-b=1" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known long flag
func TestTreeInterpreterShortSolo11(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.LongFlag("bacon")).
		MustParse([]string{"command", "leaf", "-a", "--bacon"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "--beans"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "--beans" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// solo short flag expects optional and is followed by a known long pair
func TestTreeInterpreterShortSolo13(t *testing.T) {
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.LongFlag("bees")).
		MustParse([]string{"command", "leaf", "-a", "--bees=1"})

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
	com := cli.Tree().
		WithLeaf(cli.Leaf("leaf")).
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "leaf", "-a", "--bang=1"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to be hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "--bang=1" {
		t.Error("expected flag argument to match input but it didn't")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/18
func TestRegression18Tree(t *testing.T) {
	bind := false
	com := cli.Tree().
		WithFlag(cli.ShortFlag('a').WithBinding(&bind, false)).
		WithLeaf(cli.Leaf("leaf")).
		MustParse([]string{"command", "leaf", "-a"})

	flag := com.FindShortFlag('a')

	if !bind {
		t.Error("expected bind to be true, but it wasn't")
	}

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "true" {
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

	_, err := cli.Tree().
		WithLeaf(cli.Leaf("foo").
			WithFlag(cli.ComboFlag('r', "rm-na").
				WithBinding(&removeNAValues, false)).
			WithFlag(cli.ComboFlag('s', "sorted-inputs").
				WithBindingAndDefault(&inputsAreSorted, false, true)).
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
			WithArgument(cli.Argument().
				WithName("file").
				WithBinding(func(path []string) (err error) {
					inputFile = path[0]
					return
				}))).
		Parse([]string{"build/linux/find-bin-width", "foo", "-s", "-f", "tsv", "some-file"})

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
	var value argo.Hex8

	_, err := cli.Tree().
		WithLeaf(cli.Leaf("gen-meta").
			WithFlag(cli.ComboFlag('i', "interactive").
				WithDescription("Interactive mode: auto (0), none (1), minimal (2), full (3).  Defaults to auto").
				WithBindingAndDefault(&value, argo.Hex8(23), true))).
		Parse([]string{"something", "gen-meta"})

	if err != nil {
		t.Error("expected err to be nil but was", err)
	}

	if value != 23 {
		t.Error("expected value to be 23 but was", value)
	}
}
