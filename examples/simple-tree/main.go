package main

import (
	"fmt"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func main() {
	cli.Tree().
		WithCallback(func(tree argo.CommandTree) { fmt.Print("hello") }).
		WithBranch(cli.Branch("foo").
			WithCallback(func(branch argo.CommandBranch) { fmt.Print(" to") }).
			WithLeaf(cli.Leaf("bar").
				WithCallback(func(leaf argo.CommandLeaf) { fmt.Println(" you") }))).
		WithBranch(cli.Branch("fizz").
			WithCallback(func(branch argo.CommandBranch) { fmt.Print(" you") }).
			WithLeaf(cli.Leaf("buzz").
				WithCallback(func(leaf argo.CommandLeaf) { fmt.Println(" goober") }))).
		MustParse(os.Args)
}
