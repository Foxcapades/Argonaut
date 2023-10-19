package argo_test

import (
	"testing"

	cli "github.com/Foxcapades/Argonaut"
)

// Unknown short solo flag.
func TestCommandInterpreterShortSolo01(t *testing.T) {
	com := cli.Command().MustParse([]string{"command", "-f"})

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-f" {
		t.Error("expected command unmapped input to match input but it didn't")
	}
}

// Known short solo flag.
func TestCommandInterpreterShortSolo02(t *testing.T) {
	com := cli.Command().WithFlag(cli.ShortFlag('f')).MustParse([]string{"command", "-f"})
	flag := com.FindShortFlag('f')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	}
}

// Unknown short solo at the start of a block
func TestCommandInterpreterShortSolo03(t *testing.T) {
	com := cli.Command().WithFlag(cli.ShortFlag('b')).MustParse([]string{"command", "-ab"})
	flag := com.FindShortFlag('b')

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-a" {
		t.Error("expected unmapped input to match input value but it didn't")
	}

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	}
}

// Unknown short solo in the middle of a block
func TestCommandInterpreterShortSolo04(t *testing.T) {
	com := cli.Command().WithFlag(cli.ShortFlag('a')).MustParse([]string{"command", "-ab"})
	flag := com.FindShortFlag('a')

	if !com.HasUnmappedInputs() {
		t.Error("expected command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-b" {
		t.Error("expected unmapped input to match input value but it didn't")
	}

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	}
}

// Short solo requires arg but hits eof
func TestCommandInterpreterShortSolo05(t *testing.T) {
	_, err := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument().Require())).
		Parse([]string{"command", "-a"})

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

// Short solo requires arg but hits boundary
func TestCommandInterpreterShortSolo06(t *testing.T) {
	_, err := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument().Require())).
		Parse([]string{"command", "-a", "--"})

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

// Short solo requires arg and hits any value
func TestCommandInterpreterShortSolo07(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument().Require())).
		MustParse([]string{"command", "-a", "-b"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "-b" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

// Short solo requires arg and clobbers block
func TestCommandInterpreterShortSolo08(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument().Require())).
		MustParse([]string{"command", "-ab"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "b" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}
