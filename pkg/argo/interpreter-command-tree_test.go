package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
)

func TestInvalidSubCommand(t *testing.T) {
	_, err := cli.Tree().
		WithLeaf(cli.Leaf("leaf1")).
		Parse([]string{"command", "leaf2"})

	if err == nil {
		t.Fail()
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
