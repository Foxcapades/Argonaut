package main

import (
	"fmt"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func main() {
	cli.Tree().
		WithDescription("This is a simple command tree example.").
		WithCallback(func(tree argo.CommandTree) { fmt.Print("hello") }).
		WithBranch(cli.Branch("foo").
			WithDescription("this is the description for the foo branch").
			WithAliases("fo").
			WithCallback(func(branch argo.CommandBranch) { fmt.Print(" to") }).
			WithLeaf(cli.Leaf("bar").
				WithDescription("this is the description for the bar leaf").
				WithAliases("b").
				WithCallback(func(leaf argo.CommandLeaf) { fmt.Println(" you") }))).
		WithBranch(cli.Branch("fizz").
			WithDescription("This is the description for the fizz branch.").
			WithAliases("fi").
			WithCallback(func(branch argo.CommandBranch) { fmt.Print(" you") }).
			WithLeaf(cli.Leaf("buzz").
				WithDescription("This is the description for the buzz leaf.").
				WithAliases("b").
				WithCallback(func(leaf argo.CommandLeaf) { fmt.Println(" goober") }).
				WithFlag(cli.Flag().
					WithShortForm('c').
					WithDescription("some very long description line that will need to be broken down into multiple lines to test that the help rendering functionality is correctly breaking down lines as expected.  This should be broken into at least 3 lines or possibly more depending on how much i end up writing here.  Just for good measure, lets do this:\n\nHello fools.")))).
		MustParse(os.Args)
}
