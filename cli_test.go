package cli_test

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func init() {
	utils.Exit = func(int) {}
}

// region Unit Tests

// ╔════════════════════════════════════════════════════════════════════════╗ //
// ║                                                                        ║ //
// ║    Unit Tests                                                          ║ //
// ║                                                                        ║ //
// ╚════════════════════════════════════════════════════════════════════════╝ //

// ╔════════════════════════════════════╗ //
// ║    Support Types                   ║ //
// ╚════════════════════════════════════╝ //

type nmrshlr struct {
	value string
}

func (n *nmrshlr) Unmarshal(raw string) error {
	n.value = raw
	return nil
}

// ╔════════════════════════════════════╗ //
// ║    Test Bodies                     ║ //
// ╚════════════════════════════════════╝ //

func TestFlag_withSliceOfUnmarshalable(t *testing.T) {
	var values []*nmrshlr

	os.Args = []string{"command", "-f", "goodbye", "-fcruel", "-f=world"}
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.ShortFlag('f').
			WithBinding(&values, false)))

	if len(values) != 3 {
		t.Errorf("expected values slice to have a length of 3 but was %d instead", len(values))
	}

	if values[0].value != "goodbye" {
		t.Errorf("expected value 1 to be 'goodbye' but was '%s'", values[0].value)
	}

	if values[1].value != "cruel" {
		t.Errorf("expected value 2 to be 'cruel' but was '%s'", values[1].value)
	}

	if values[2].value != "world" {
		t.Errorf("expected value 2 to be 'world' but was '%s'", values[2].value)
	}
}

func TestFlag_withMapOfUnmarshalable(t *testing.T) {
	var values map[string]*nmrshlr

	os.Args = []string{"command", "-v", "foo:bar", "-vfizz:buzz", "-v=happy:sad"}
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.ShortFlag('v').
			WithBinding(&values, true)))

	if len(values) != 3 {
		t.Errorf("expected values map to have a length of 3 but was %d", len(values))
	}

	if val, ok := values["foo"]; !ok {
		t.Error("expected map key was not present")
	} else if val.value != "bar" {
		t.Error("expected map value to match input")
	}

	if val, ok := values["fizz"]; !ok {
		t.Error("expected map key was not present")
	} else if val.value != "buzz" {
		t.Error("expected map value to match input")
	}

	if val, ok := values["happy"]; !ok {
		t.Error("expected map key was not present")
	} else if val.value != "sad" {
		t.Error("expected map value to match input")
	}
}

func TestFlag_withMapOfSliceOfUnmarshalable(t *testing.T) {
	var values map[string][]*nmrshlr

	os.Args = []string{"command", "-v", "foo:bar", "-vfoo:fizz", "-v=foo:buzz"}
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.ShortFlag('v').
			WithBinding(&values, true)))

	vals, ok := values["foo"]

	if !ok {
		t.Error("expected map to contain input map key")
	} else {
		if len(vals) != 3 {
			t.Errorf("expected values map to have a length of 3 but was %d", len(vals))
		}

		if val := vals[0]; val.value != "bar" {
			t.Error("expected map value to match input")
		}

		if val := vals[1]; val.value != "fizz" {
			t.Error("expected map value to match input")
		}

		if val := vals[2]; val.value != "buzz" {
			t.Error("expected map value to match input")
		}
	}
}

func TestTree_handleIncompleteViaCallback(t *testing.T) {
	os.Args = []string{"root"}

	hit := false

	_, _, _ = cli.ParseTree(cli.Tree().
		WithIncompleteHandler(func(argo.TreeCommand) { hit = true }).
		WithLeaf(cli.Leaf("leaf")))

	if !hit {
		t.Error("tree root callback was not hit")
	}
}

func TestArgument_default01(t *testing.T) {
	val := 0
	foo := func() int { return 3 }

	os.Args = []string{"command"}
	com, _ := cli.MustParseCommand(cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)))

	arg := com.Arguments()[0]

	if !arg.WasHit() {
		t.Error("expected argument to have been hit but it wasn't")
	} else if val != 3 {
		t.Error("expected bind value to equal the given function return value but it didn't")
	}
}

func TestArgument_default02(t *testing.T) {
	val := 0
	foo := func() (int, error) { return 3, nil }

	os.Args = []string{"command"}
	com, _ := cli.MustParseCommand(cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)))

	arg := com.Arguments()[0]

	if !arg.WasHit() {
		t.Error("expected argument to have been hit but it wasn't")
	} else if val != 3 {
		t.Error("expected bind value to equal the given function return value but it didn't")
	}
}

func TestArgument_default03(t *testing.T) {
	val := 0
	foo := func() (int, error) { return 0, errors.New("butt") }

	os.Args = []string{"command"}
	_, _, err := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().WithBinding(&val).WithDefault(foo)))

	if err == nil {
		t.Error("expected parsing to error but it didn't")
	}
}

func TestArgument_default04(t *testing.T) {
	con := argo.UnmarshalerFunc(func(val string) error { return nil })
	foo := func() (int, error) { return 3, nil }

	os.Args = []string{"command", "poo"}
	_, res, _ := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().WithBinding(con).WithDefault(foo)))

	if res.Error == nil {
		t.Error("expected a configuration error but err was nil")
	} else if res.ErrorType != argo.ConfigurationError {
		t.Errorf("expected a configuration error but err type was %s", res.ErrorType)
	}
}

type argUnmarshaler struct{ val int }

func (a *argUnmarshaler) Unmarshal(raw string) (err error) { a.val, err = strconv.Atoi(raw); return }

// Expect a failure because the unmarshaler type is incompatible with the output
// of the default value provider.
func TestArgument_default05(t *testing.T) {
	con := argUnmarshaler{}
	foo := func() (int, error) { return 3, nil }

	os.Args = []string{"command", "poo"}
	_, res, _ := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().WithBinding(con).WithDefault(foo)))

	if res.Error == nil {
		t.Error("expected a configuration error but err was nil")
	} else if res.ErrorType != argo.ConfigurationError {
		t.Errorf("expected a configuration error but err type was %s", res.ErrorType)
	}
}

// Expect an OK because the unmarshaler type is compatible with the output of
// the default value provider.
// TODO: What should this test actually be?  It doesn't make sense to have the
//       provider return an unmarshaler instance.
// func TestArgumentDefault06(t *testing.T) {
// 	con := argUnmarshaler{}
// 	foo := func() (argUnmarshaler, error) { return argUnmarshaler{3}, nil }
// 	_ = cli.Command().
// 		WithArgument(cli.Argument().WithBinding(&con).WithDefault(foo)).
// 		MustParse([]string{"command", "3"})
//
// 	if con.val != 3 {
// 		t.Error("expected unmarshaler value to have been replaced but it wasn't")
// 	}
// }

// Expect an OK because the unmarshaler type is compatible with the output of
// the default value provider.
// TODO: What should this test actually be?  It doesn't make sense to have the
//       provider return an unmarshaler instance.
// func TestArgumentDefault07(t *testing.T) {
// 	con := argUnmarshaler{}
// 	foo := func() (argUnmarshaler, error) { return argUnmarshaler{3}, nil }
// 	_ = cli.Command().
// 		WithArgument(cli.Argument().WithBinding(&con).WithDefault(foo)).
// 		MustParse([]string{"command"})
//
// 	if con.val != 3 {
// 		t.Error("expected unmarshaler value to have been replaced but it wasn't")
// 	}
// }

func TestArgument_PreParseValidator01(t *testing.T) {
	var binding int

	os.Args = []string{"command", "32"}
	_, _, err := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(string) error { return errors.New("dummy error") }).
			WithBinding(&binding)))

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else if err.Error() != "dummy error" {
		t.Error("expected err to match validator output but it didn't")
	}
	t.Log(err)
}

func TestArgument_PreParseValidator02(t *testing.T) {
	var binding int

	os.Args = []string{"command", "32"}
	_, _, err := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(string) error { return nil }).
			WithBinding(&binding)))

	if err != nil {
		t.Error("expected err to be nil but it wasn't")
	}
	t.Log(err)
}

func TestArgument_PostParseValidator01(t *testing.T) {
	var binding int

	os.Args = []string{"command", "32"}
	_, _, err := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(int, string) error { return errors.New("dummy error") }).
			WithBinding(&binding)))

	if err == nil {
		t.Error("expected err not to be nil but it was")
	} else if err.Error() != "dummy error" {
		t.Errorf("expected err to match validator output but it was: %s", err)
	}
	t.Log(err)
}

func TestArgument_PostParseValidator02(t *testing.T) {
	var binding int

	os.Args = []string{"command", "32"}
	_, _, err := cli.ParseCommand(cli.Command().
		WithArgument(cli.Argument().
			WithValidator(func(int, string) error { return nil }).
			WithBinding(&binding)))

	if err != nil {
		t.Error("expected err to be nil but it wasn't")
	}
	t.Log(err)
}

// endregion Unit Tests

// region Examples

// ╔════════════════════════════════════════════════════════════════════════╗ //
// ║                                                                        ║ //
// ║    Examples                                                            ║ //
// ║                                                                        ║ //
// ╚════════════════════════════════════════════════════════════════════════╝ //

// ╔════════════════════════════════════╗ //
// ║    Package Examples                ║ //
// ╚════════════════════════════════════╝ //

// Builds a command that takes the boolean flag option -0 (or --nil-delim) and
// binds that flag to a struct field that will be set on CLI parse.
func Example_boolBinding() {
	// Simulate user call
	os.Args = []string{"command", "-0"}

	// Variable to bind
	useNilDelim := false

	// Build & parse command with variable bound to the -0/--nil-delim flag.
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.ComboFlag('0', "nil-delim").
			WithDescription("End output with a null byte instead of a newline.").
			WithBinding(&useNilDelim, false)))

	fmt.Printf("%t", useNilDelim)
	// Output: true
}

// Builds a command that takes string inputs to the flag option -s (or --string)
// and uses those string inputs to populate a bound slice.
func Example_sliceBinding() {
	// Simulate user call
	os.Args = []string{"command", "-s", "hello", "--string", "goodbye"}

	// Variable to populate
	var mySlice []string

	// Build & parse command with variable bound to the -s/--string flag.
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.ComboFlag('s', "string").
			WithBinding(&mySlice, true)))

	fmt.Print(mySlice)
	// Output: [hello goodbye]
}

// ╔════════════════════════════════════╗ //
// ║    cli.Command()                   ║ //
// ╚════════════════════════════════════╝ //

// Example using the UnmappedInputs method to retrieve CLI call inputs that were
// not mapped to a positional argument.
func ExampleCommand_extraInputs() {
	// Simulate user call.
	os.Args = []string{"command", "foo", "bar", "fizz", "buzz"}

	cli.MustParseCommand(cli.Command().
		WithCallback(func(command argo.Command) {
			fmt.Println(command.UnmappedInputs())
		}))
	// Output: [foo bar fizz buzz]
}

// ╔════════════════════════════════════╗ //
// ║    cli.Tree()                      ║ //
// ╚════════════════════════════════════╝ //

func ExampleTree() {
	os.Args = []string{"command", "foo", "--", "bar"}
	cli.MustParseTree(cli.Tree().
		WithLeaf(cli.Leaf("foo").
			WithCallback(func(leaf argo.LeafCommand) {
				fmt.Println(leaf.UnmappedInputs())
			})).
		WithLeaf(cli.Leaf("bar").
			WithCallback(func(leaf argo.LeafCommand) {
				panic(leaf)
			})))

	// Output: [bar]
}

// Demonstration of tree node callback execution order.
func ExampleTree_callbackExecutionOrder() {
	// Simulate user input
	os.Args = []string{"command", "foo", "bar"}

	// Build tree with callbacks
	cli.MustParseTree(cli.Tree().
		WithCallback(func(argo.TreeCommand) {
			fmt.Print("Hello")
		}).
		WithBranch(cli.Branch("foo").
			WithCallback(func(branch argo.BranchCommand) {
				fmt.Print(" to")
			}).
			WithLeaf(cli.Leaf("bar").
				WithCallback(func(leaf argo.LeafCommand) {
					fmt.Println(" you!")
				}))))
	// Output: Hello to you!
}

// Demonstrates handling incomplete commands with a provided callback function.
//
// Typically, the incomplete callback function should exit when finished.  If it
// doesn't, the error value returned from the ParseTree call will contain an
// error for the tree command CLI call being incomplete.
func ExampleTree_handleIncompleteViaCallback() {
	// Simulate incomplete CLI call
	os.Args = []string{"root"}

	_, _, _ = cli.ParseTree(cli.Tree().
		WithIncompleteHandler(func(argo.TreeCommand) {
			fmt.Println("incomplete command!")
			// os.Exit(1)  // commented out for example execution
		}).
		WithLeaf(cli.Leaf("leaf")))
	// Output: incomplete command!
}

// ╔════════════════════════════════════╗ //
// ║    cli.Branch()                    ║ //
// ╚════════════════════════════════════╝ //

// Defines a branch subcommand with 2 leaf commands under it.
func ExampleBranch_simple() {
	cli.Branch("branch-name").
		WithLeaves(
			cli.Leaf("leaf-1"),
			cli.Leaf("leaf-2"),
		)
}

// Explanation of common branch builder methods.
func ExampleBranch_expanded() {
	cli.Branch("level-1").
		WithAliases("lvl1", "l1").                          // alternate names for the subcommand
		WithDescription("first level command tree node").   // description for help text
		WithFlag(cli.Flag()).                               // an ungrouped flag
		WithFlagGroup(cli.FlagGroup("flag group")).         // a group of flags
		WithIncompleteHandler(func(argo.BranchCommand) {}). // a custom handler for incomplete CLI calls
		WithCallback(func(argo.BranchCommand) {}).          // a callback for successful CLI calls
		WithLeaf(cli.Leaf("level-2")).                      // an ungrouped subcommand
		WithCommandGroup(cli.CommandGroup("command group")) // a group of subcommands
}

// func ExampleBranch() {
// 	// Simulate user input
// 	os.Args = []string{"command", "foo", "bar"}
//
// 	// Spec out the branch
// 	branch := cli.Branch("foo").
// 		WithCallback(func(branch argo.BranchCommand) {
// 			fmt.Print("hello from ")
// 		}).
// 		WithLeaf(cli.Leaf("bar").
// 			WithCallback(func(leaf argo.LeafCommand) {
// 				fmt.Println("a branch!")
// 			}))
//
// 	cli.MustParseCommand(cli.Tree().WithBranch(branch))
// 	// Output: hello from a branch!
// }

// ╔════════════════════════════════════╗ //
// ║    cli.Leaf()                      ║ //
// ╚════════════════════════════════════╝ //

func ExampleLeaf() {
	var zone string

	os.Args = []string{"command", "time", "UTC"}
	cli.MustParseTree(cli.Tree().
		WithLeaf(cli.Leaf("time").
			WithArgument(cli.Argument().
				WithName("zone").
				WithBinding(&zone))))

	fmt.Println(zone)
	// Output: UTC
}

// ╔════════════════════════════════════╗ //
// ║    cli.CommandGroup()              ║ //
// ╚════════════════════════════════════╝ //

func ExampleCommandGroup() {
	os.Args = []string{"command", "foo"}
	com, _ := cli.MustParseTree(cli.Tree().
		WithCommandGroup(cli.CommandGroup("my commands").
			WithDescription("a group of commands for me").
			WithLeaf(cli.Leaf("foo")).
			WithLeaf(cli.Leaf("bar"))))

	fmt.Println(com.SelectedCommand().Name())
	// Output: foo
}

// ╔════════════════════════════════════╗ //
// ║    cli.Flag()                      ║ //
// ╚════════════════════════════════════╝ //

// Demonstrates that repeated usage of a flag results in a single execution of
// that flag's defined callback.
func ExampleFlag_repeatedUsageWithCallback() {
	os.Args = []string{"command", "-ssss", "--selection", "--selection"}
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.Flag().
			WithShortForm('s').
			WithLongForm("selection").
			WithCallback(func(flag argo.Flag) { fmt.Println(flag.HitCount()) })))
	// Output: 6
}

// ╔════════════════════════════════════╗ //
// ║    cli.FlagGroup()                 ║ //
// ╚════════════════════════════════════╝ //

func ExampleFlagGroup() {
	// Simulate CLI call
	os.Args = []string{"command", "-c", "--clutch"}

	// Define the flag groups
	mine := cli.FlagGroup("my flags").
		WithFlag(cli.ShortFlag('c').
			WithCallback(func(flag argo.Flag) { fmt.Print("hello ") }))
	yours := cli.FlagGroup("your flags").
		WithFlag(cli.LongFlag("clutch").
			WithCallback(func(flag argo.Flag) { fmt.Println("world") }))

	// Parse CLI call
	cli.MustParseCommand(cli.Command().WithFlagGroups(mine, yours))
	// Output: hello world
}

func ExampleFlagGroup_helpRendering() {
	// Simulate CLI call
	os.Args = []string{"command", "-h"}

	// Define the flag groups
	mine := cli.FlagGroup("my flags").
		WithFlag(cli.ShortFlag('c').
			WithCallback(func(flag argo.Flag) { fmt.Print("hello ") }))

	yours := cli.FlagGroup("your flags").
		WithFlag(cli.LongFlag("clutch").
			WithCallback(func(flag argo.Flag) { fmt.Println("world") }))

	// Parse CLI call
	cli.MustParseCommand(cli.Command().WithFlagGroups(mine, yours))
	// Output:
	// Usage:
	//   command [options]
	//
	// my flags
	//   -c
	//
	// your flags
	//   --clutch
	//
	// General Flags
	//   -h | --help
	//       Prints this help text.
}

// ╔════════════════════════════════════╗ //
// ║    cli.LongFlag()                  ║ //
// ╚════════════════════════════════════╝ //

func ExampleLongFlag() {
	os.Args = []string{"command", "--hello"}
	cli.MustParseCommand(cli.Command().
		WithFlag(cli.LongFlag("hello").
			WithCallback(func(flag argo.Flag) {
				fmt.Println(flag.WasHit())
			})))

	// Output: true
}

// ╔════════════════════════════════════╗ //
// ║    cli.ShortFlag()                 ║ //
// ╚════════════════════════════════════╝ //

// // Demonstrates that repeated usage of a flag results in a single execution of
// // that flag's defined callback.
// func ExampleShortFlag_repeatedUsageWithCallback() {
// 	// Simulate CLI call
// 	os.Args = []string{"command", "-aaa", "-a", "-a"}
//
// 	// Define short flag
// 	flag := cli.ShortFlag('a').
// 		WithCallback(func(flag argo.Flag) { fmt.Println(flag.HitCount()) })
//
// 	// Parse CLI call
// 	cli.MustParseCommand(cli.Command().WithFlag(flag))
// 	// Output: 5
// }

// ╔════════════════════════════════════╗ //
// ║    cli.Argument()                  ║ //
// ╚════════════════════════════════════╝ //

func ExampleArgument_positionalCommandArguments() {
	// Simulate CLI call
	os.Args = []string{"command", "foo.txt", "36"}

	// Target variables
	var file string
	var count uint

	// Define arguments
	fileArg := cli.Argument().WithBinding(&file)
	countArg := cli.Argument().WithBinding(&count)

	// Parse CLI call
	cli.MustParseCommand(cli.Command().WithArguments(fileArg, countArg))

	fmt.Println(file, count)
	// Output: foo.txt 36
}

func ExampleArgument_defaultValues() {
	// Simulate CLI call
	os.Args = []string{"command"}

	// Target variables
	var width int
	var height int

	// Define arguments
	widthArg := cli.Argument().
		WithName("width").
		WithBinding(&width).
		WithDefault(800)
	heightArg := cli.Argument().
		WithName("height").
		WithBinding(&height).
		WithDefault(600)

	// Parse CLI call
	cli.MustParseCommand(cli.Command().WithArguments(widthArg, heightArg))

	fmt.Printf("%dx%d", width, height)
	// Output: 800x600
}
