package cli_test

import (
	"fmt"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

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
