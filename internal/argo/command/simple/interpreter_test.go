package command_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/simple"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
	"github.com/foxcapades/argonaut/v3/pkg/argotype"
)

// Unknown short solo flag.
func TestCommandInterpreterShortSolo01(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder(), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-f"}))

	if !com.HasUnmappedInputs() {
		t.Error("expected Command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected Command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-f" {
		t.Error("expected Command unmapped input to match input but it didn't")
	}
}

// Known short solo flag.
func TestCommandInterpreterShortSolo02(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().WithFlag(flag.NewBuilder().WithShortForm('f')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-f"}))

	if !com.FindShortFlag('f').WasHit() {
		t.Error("expected flag to have been hit, but it wasn't")
	}
}

// Unknown short solo at the start of a block
func TestCommandInterpreterShortSolo03(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().WithFlag(flag.NewBuilder().WithShortForm('b')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-ab"}))

	if !com.HasUnmappedInputs() {
		t.Error("expected Command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected Command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-a" {
		t.Error("expected unmapped input to match input value but it didn't")
	}

	if !com.FindShortFlag('b').WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	}
}

// Unknown short solo in the middle of a block
func TestCommandInterpreterShortSolo04(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().WithFlag(flag.NewBuilder().WithShortForm('a')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-ab"}))

	if !com.HasUnmappedInputs() {
		t.Error("expected Command to have unmapped inputs but it didn't")
	} else if len(com.UnmappedInputs()) != 1 {
		t.Error("expected Command to have exactly 1 unmapped input but it didn't")
	} else if com.UnmappedInputs()[0] != "-b" {
		t.Error("expected unmapped input to match input value but it didn't")
	}

	if !com.FindShortFlag('a').WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	}
}

// Short solo requires arg but hits eof
func TestCommandInterpreterShortSolo05(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder().Require())), opt))
	_, err := command.Parse(com, []string{"Command", "-a"})

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

// Short solo requires arg but hits boundary
func TestCommandInterpreterShortSolo06(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder().Require())), opt))
	_, err := command.Parse(com, []string{"Command", "-a", "--"})

	if err == nil {
		t.Error("expected error not to be nil, but it was")
	}
}

// Short solo requires arg and hits any value
func TestCommandInterpreterShortSolo07(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder().Require())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "-b"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "-b" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

// Short solo requires arg and clobbers block
func TestCommandInterpreterShortSolo08(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder().Require())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-ab"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "b" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo09(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithShortForm('b')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-ab"}))

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
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-ab"}))

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
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

func TestCommandInterpreterShortSolo12(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if f.Argument().WasHit() {
		t.Error("expected flag argument not to have been hit but it was")
	}
}

func TestCommandInterpreterShortSolo13(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "lamp"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "lamp" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo14(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "-l"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "-l" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo15(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--paul"}))

	f := com.FindShortFlag('a')

	if !f.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f.Argument().RawValue() != "--paul" {
		t.Error("expected flag argument to equal input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo16(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithShortForm('l')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "-l"}))

	f1 := com.FindShortFlag('a')
	f2 := com.FindShortFlag('l')

	if !f1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if f1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !f2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo17(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithLongForm("atom")), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--atom"}))

	f1 := com.FindShortFlag('a')
	f2 := com.FindLongFlag("atom")

	if !f1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if f1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !f2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo18(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithShortForm('l')), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "-l=1"}))

	f1 := com.FindShortFlag('a')
	f2 := com.FindShortFlag('l')

	if !f1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if f1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !f2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo19(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())).
		WithFlag(flag.NewBuilder().WithLongForm("atom")), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--atom=1"}))

	f1 := com.FindShortFlag('a')
	f2 := com.FindLongFlag("atom")

	if !f1.WasHit() {
		t.Error("expected flag 1 to have been hit but it wasn't")
	} else if f1.Argument().WasHit() {
		t.Error("expected flag 1 argument to not have been hit but it was")
	}

	if !f2.WasHit() {
		t.Error("expected flag 2 to have been hit but it wasn't")
	}
}

func TestCommandInterpreterShortSolo20(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "-l=1"}))

	f1 := com.FindShortFlag('a')

	if !f1.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f1.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f1.Argument().RawValue() != "-l=1" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

func TestCommandInterpreterShortSolo21(t *testing.T) {
	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('a').
			WithArgument(argument.NewBuilder())), opt))
	_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--atom=1"}))

	f1 := com.FindShortFlag('a')

	if !f1.WasHit() {
		t.Error("expected flag to have been hit but it wasn't")
	} else if !f1.Argument().WasHit() {
		t.Error("expected flag argument to have been hit but it wasn't")
	} else if f1.Argument().RawValue() != "--atom=1" {
		t.Error("expected flag argument value to match input value but it didn't")
	}
}

// https://github.com/Foxcapades/Argonaut/issues/18
func TestRegression18Command(t *testing.T) {
	opt := argo.Options{}
	{
		bind := false
		com := utils.MustReturn(command.Build(command.NewBuilder().
			WithFlag(flag.NewBuilder().WithShortForm('a').WithBinding(&bind, false)), opt))
		_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a"}))

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
	{
		bind := false
		com := utils.MustReturn(command.Build(command.NewBuilder().
			WithFlag(flag.NewBuilder().WithShortForm('a').WithBinding(&bind, false)), opt))
		_ = utils.MustReturn(command.Parse(com, []string{"Command", "-a", "--", "flumps"}))

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

		if !com.HasUnmappedInputs() {
			t.Error("expected flag to have passthrough arguments but it didn't")
		} else if len(com.UnmappedInputs()) != 1 {
			t.Error("expected flag to have exactly 1 passthrough argument but it didn't")
		} else if com.UnmappedInputs()[0] != "flumps" {
			t.Error("expected flag passthrough argument to match input value but it didn't")
		}
	}
}

// https://github.com/Foxcapades/Argonaut/issues/58
func TestRegression58Command(t *testing.T) {
	var removeNAValues bool
	var inputsAreSorted bool
	var outputFormat uint8
	var printHeaders bool
	var inputFile string

	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('r').WithLongForm("rm-na").
			WithBinding(&removeNAValues, false)).
		WithFlag(flag.NewBuilder().WithShortForm('s').WithLongForm("sorted-inputs").
			WithBindingAndDefault(&inputsAreSorted, false, false)).
		WithFlag(flag.NewBuilder().WithShortForm('f').WithLongForm("format").
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
		WithFlag(flag.NewBuilder().WithShortForm('t').WithLongForm("headers").
			WithBinding(&printHeaders, false)).
		WithArgument(argument.NewBuilder().
			WithName("file").
			WithBinding(func(path []string) (err error) {
				inputFile = path[0]
				return
			})), opt))
	_, err := command.Parse(com, []string{"build/linux/find-bin-width", "-s", "-f", "tsv", "some-file"})

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
func TestRegression62Command(t *testing.T) {
	var value argotype.Hex8

	opt := argo.Options{}
	com := utils.MustReturn(command.Build(command.NewBuilder().
		WithFlag(flag.NewBuilder().WithShortForm('i').WithLongForm("interactive").
			WithDescription("Interactive mode: auto (0), none (1), minimal (2), full (3).  Defaults to auto").
			WithBindingAndDefault(&value, argotype.Hex8(23), true)), opt))
	_, err := command.Parse(com, []string{"something", "gen-meta"})

	if err != nil {
		t.Error("expected err to be nil but was", err)
	}

	if value != 23 {
		t.Error("expected value to be 23 but was", value)
	}
}
