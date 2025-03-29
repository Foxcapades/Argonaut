package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type nmrshlr struct {
	value string
}

func (n *nmrshlr) Unmarshal(raw string) error {
	n.value = raw
	return nil
}

func TestFlag_withSliceOfUnmarshalable(t *testing.T) {
	var values []*nmrshlr

	os.Args = []string{"command", "-f", "goodbye", "-fcruel", "-f=world"}
	cli.MustParse(cli.Command().
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
	cli.MustParse(cli.Command().
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
	cli.MustParse(cli.Command().
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

func ExampleCommand() {
	os.Args = []string{"command", "foo", "bar", "fizz", "buzz"}
	cli.MustParse(cli.Command().
		WithCallback(func(command argo.Command) {
			fmt.Println(command.UnmappedInputs())
		}))

	// Output: [foo bar fizz buzz]
}

func ExampleArgument() {
	var file string
	var count uint

	os.Args = []string{"command", "foo.txt", "36"}
	cli.MustParse(cli.Command().
		WithArgument(cli.Argument().
			WithName("file").
			WithBinding(&file)).
		WithArgument(cli.Argument().
			WithName("count").
			WithBinding(&count)))

	fmt.Println(file, count)

	// Output: foo.txt 36
}

func ExampleFlag() {
	os.Args = []string{"command", "-ssss", "--selection", "--selection"}
	cli.MustParse(cli.Command().
		WithFlag(cli.Flag().
			WithShortForm('s').
			WithLongForm("selection").
			WithCallback(func(flag argo.Flag) { fmt.Println(flag.HitCount()) })))

	// Output: 6
}

func ExampleCommand_complex() {
	var config = struct {
		NilDelim bool
	}{}

	os.Args = []string{"command", "-0"}
	cli.MustParse(cli.Command().
		WithFlagGroup(cli.FlagGroup("Output Control").
			WithFlag(cli.Flag().
				WithShortForm('0').
				WithLongForm("nil-delim").
				WithDescription("End output with a null byte instead of a newline.").
				WithBinding(&config.NilDelim, false))))

	fmt.Println(config.NilDelim)

	// Output: true
}

func ExampleTree() {
	os.Args = []string{"command", "foo", "--", "bar"}
	cli.MustParse(cli.Tree().
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

func ExampleBranch() {
	os.Args = []string{"command", "foo", "bar"}
	cli.MustParse(cli.Tree().
		WithBranch(cli.Branch("foo").
			WithCallback(func(branch argo.BranchCommand) {
				fmt.Print("hello from ")
			}).
			WithLeaf(cli.Leaf("bar").
				WithCallback(func(leaf argo.LeafCommand) {
					fmt.Println("a branch!")
				}))))

	// Output: hello from a branch!
}

func ExampleLeaf() {
	var zone string

	os.Args = []string{"command", "time", "UTC"}
	cli.MustParse(cli.Tree().
		WithLeaf(cli.Leaf("time").
			WithArgument(cli.Argument().
				WithName("zone").
				WithBinding(&zone))))

	fmt.Println(zone)

	// Output: UTC
}

func ExampleCommandGroup() {
	os.Args = []string{"command", "foo"}
	com := cli.MustBuildTree(cli.Tree().
		WithCommandGroup(cli.CommandGroup("my commands").
			WithDescription("a group of commands for me").
			WithLeaf(cli.Leaf("foo")).
			WithLeaf(cli.Leaf("bar"))))
	cli.MustParse(com)

	fmt.Println(com.SelectedCommand().Name())

	// Output: foo
}

func ExampleFlagGroup() {
	os.Args = []string{"command", "-c", "--clutch"}
	cli.MustParse(cli.Command().
		WithFlagGroup(cli.FlagGroup("my flags").
			WithFlag(cli.ShortFlag('c').
				WithCallback(func(flag argo.Flag) { fmt.Print("hello ") }))).
		WithFlagGroup(cli.FlagGroup("your flags").
			WithFlag(cli.LongFlag("clutch").
				WithCallback(func(flag argo.Flag) { fmt.Println("world") }))))

	// Output: hello world
}

func ExampleLongFlag() {
	os.Args = []string{"command", "--hello"}
	cli.MustParse(cli.Command().
		WithFlag(cli.LongFlag("hello").
			WithCallback(func(flag argo.Flag) {
				fmt.Println(flag.WasHit())
			})))

	// Output: true
}

func ExampleShortFlag() {
	os.Args = []string{"command", "-aaa", "-a", "-a"}
	cli.MustParse(cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithCallback(func(flag argo.Flag) { fmt.Println(flag.HitCount()) })))

	// Output: 5
}
