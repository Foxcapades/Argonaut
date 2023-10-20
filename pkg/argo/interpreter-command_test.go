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

func TestCommandInterpreterShortSolo09(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('b')).
		MustParse([]string{"command", "-ab"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('b')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo10(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').WithArgument(cli.Argument())).
		MustParse([]string{"command", "-ab"})

	flag1 := com.FindShortFlag('a')

	if !flag1.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag1.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag1.Argument().RawValue() != "b" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo11(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

func TestCommandInterpreterShortSolo12(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "--"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if flag.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

func TestCommandInterpreterShortSolo13(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "lamp"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "lamp" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo14(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "-l"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "-l" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo15(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "--paul"})

	flag := com.FindShortFlag('a')

	if !flag.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag.Argument().RawValue() != "--paul" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo16(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('l')).
		MustParse([]string{"command", "-a", "-l"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('l')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo17(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		WithFlag(cli.LongFlag("atom")).
		MustParse([]string{"command", "-a", "--atom"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindLongFlag("atom")

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo18(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		WithFlag(cli.ShortFlag('l')).
		MustParse([]string{"command", "-a", "-l=1"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindShortFlag('l')

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo19(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		WithFlag(cli.LongFlag("atom")).
		MustParse([]string{"command", "-a", "--atom=1"})

	flag1 := com.FindShortFlag('a')
	flag2 := com.FindLongFlag("atom")

	if !flag1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if flag1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !flag2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo20(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "-l=1"})

	flag1 := com.FindShortFlag('a')

	if !flag1.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag1.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag1.Argument().RawValue() != "-l=1" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo21(t *testing.T) {
	com := cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithArgument(cli.Argument())).
		MustParse([]string{"command", "-a", "--atom=1"})

	flag1 := com.FindShortFlag('a')

	if !flag1.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !flag1.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if flag1.Argument().RawValue() != "--atom=1" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/18
func TestRegression18Command(t *testing.T) {
	{
		bind := false
		com := cli.Command().
			WithFlag(cli.ShortFlag('a').WithBinding(&bind, false)).
			MustParse([]string{"command", "-a"})

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
	{
		bind := false
		com := cli.Command().
			WithFlag(cli.ShortFlag('a').WithBinding(&bind, false)).
			MustParse([]string{"command", "-a", "--", "flumps"})

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

		if !com.HasPassthroughInputs() {
			t.Error("expected flag to have passthrough arguments but it didn't")
		} else if len(com.PassthroughInputs()) != 1 {
			t.Error("expected flag to have exactly 1 passthrough argument but it didn't")
		} else if com.PassthroughInputs()[0] != "flumps" {
			t.Error("expected flag passthrough argument to match input value but it didn't")
		}
	}
}
