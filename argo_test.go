package cli_test

import (
	"fmt"
	"testing"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
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

	cli.Command().
		WithFlag(cli.ShortFlag('f').
			WithBinding(&values, true)).
		MustParse([]string{"command", "-f", "goodbye", "-fcruel", "-f=world"})

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

func ExampleCommand() {
	cli.Command().
		WithCallback(func(command argo.Command) {
			fmt.Println(command.UnmappedInputs())
		}).
		MustParse([]string{"command", "foo", "bar", "fizz", "buzz"})

	// Output: [foo bar fizz buzz]
}

func ExampleArgument() {
	var file string
	var count uint

	cli.Command().
		WithArgument(cli.Argument().
			WithName("file").
			WithBinding(&file)).
		WithArgument(cli.Argument().
			WithName("count").
			WithBinding(&count)).
		MustParse([]string{"command", "foo.txt", "36"})

	fmt.Println(file, count)

	// Output: foo.txt 36
}

func ExampleFlag() {
	cli.Command().
		WithFlag(cli.Flag().
			WithShortForm('s').
			WithLongForm("selection").
			WithCallback(func(flag argo.Flag) {
				fmt.Println(flag.HitCount())
			})).
		MustParse([]string{"command", "-ssss", "--selection", "--selection"})

	// Output: 6
}

func ExampleCommand_complex() {
	var config = struct {
		NilDelim bool
	}{}

	cli.Command().
		WithFlagGroup(cli.FlagGroup("Output Control").
			WithFlag(cli.Flag().
				WithShortForm('0').
				WithLongForm("nil-delim").
				WithDescription("End output with a null byte instead of a newline.").
				WithBinding(&config.NilDelim, true))).
		MustParse([]string{"command", "-0"})

	fmt.Println(config.NilDelim)

	// Output: true
}

func ExampleTree() {
	cli.Tree().
		WithLeaf(cli.Leaf("foo").
			WithCallback(func(leaf argo.CommandLeaf) {
				fmt.Println(leaf.PassthroughInputs())
			})).
		MustParse([]string{"command", "foo", "--", "bar"})

	// Output: [bar]
}

func ExampleBranch() {
	cli.Tree().
		WithBranch(cli.Branch("foo").
			WithCallback(func(branch argo.CommandBranch) {
				fmt.Print("hello from ")
			}).
			WithLeaf(cli.Leaf("bar").
				WithCallback(func(leaf argo.CommandLeaf) {
					fmt.Println("a branch!")
				}))).
		MustParse([]string{"command", "foo", "bar"})

	// Output: hello from a branch!
}

func ExampleLeaf() {
	var zone string

	cli.Tree().
		WithLeaf(cli.Leaf("time").
			WithArgument(cli.Argument().
				WithName("zone").
				WithBinding(&zone))).
		MustParse([]string{"command", "time", "UTC"})

	fmt.Println(zone)
	// Output: UTC
}

func ExampleCommandGroup() {
	com := cli.Tree().
		WithCommandGroup(cli.CommandGroup("my commands").
			WithDescription("a group of commands for me").
			WithLeaf(cli.Leaf("foo")).
			WithLeaf(cli.Leaf("bar"))).
		MustParse([]string{"command", "foo"})

	fmt.Println(com.SelectedCommand().Name())
	// Output: foo
}

func ExampleFlagGroup() {
	cli.Command().
		WithFlagGroup(cli.FlagGroup("my flags").
			WithFlag(cli.ShortFlag('c').
				WithCallback(func(flag argo.Flag) { fmt.Print("hello ") }))).
		WithFlagGroup(cli.FlagGroup("your flags").
			WithFlag(cli.LongFlag("clutch").
				WithCallback(func(flag argo.Flag) { fmt.Println("world") }))).
		MustParse([]string{"command", "-c", "--clutch"})

	// Output: hello world
}

func ExampleLongFlag() {
	cli.Command().
		WithFlag(cli.LongFlag("hello").
			WithCallback(func(flag argo.Flag) {
				fmt.Println(flag.WasHit())
			})).
		MustParse([]string{"command", "--hello"})

	// Output: true
}

func ExampleShortFlag() {
	cli.Command().
		WithFlag(cli.ShortFlag('a').
			WithCallback(func(flag argo.Flag) { fmt.Println(flag.HitCount()) })).
		MustParse([]string{"command", "-aaa", "-a", "-a"})

	// Output: 5
}
