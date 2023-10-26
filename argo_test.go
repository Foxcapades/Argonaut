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
