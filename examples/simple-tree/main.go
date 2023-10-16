package main

import (
	"fmt"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func main() {
	cli.CommandTree().
		WithBranch(cli.CommandBranch("foo").
			WithLeaf(cli.CommandLeaf("bar").
				WithCallback(func(leaf argo.CommandLeaf) { fmt.Println("hello world") }))).
		MustParse(os.Args)
}
