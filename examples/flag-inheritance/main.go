package main

import (
	"os"

	cli "github.com/Foxcapades/Argonaut"
)

func main() {
	cli.Tree().
		WithFlag(cli.ComboFlag('a', "apple").
			WithDescription("The apple flag.")).
		WithFlag(cli.ComboFlag('b', "banana").
			WithDescription("The first banana flag.")).
		WithFlag(cli.ComboFlag('c', "cantaloupe").
			WithDescription("The cantaloupe flag.")).
		WithFlag(cli.ComboFlag('d', "durian").
			WithDescription("The first durian flag.")).
		WithBranch(cli.Branch("foo").
			WithFlag(cli.ComboFlag('a', "acorn").
				WithDescription("The acorn flag.")).
			WithFlag(cli.ComboFlag('n', "banana").
				WithDescription("The second banana flag.")).
			WithLeaf(cli.Leaf("bar").
				WithFlag(cli.ComboFlag('c', "cabbage").
					WithDescription("The cabbage flag.")).
				WithFlag(cli.ComboFlag('u', "durian").
					WithDescription("The second durian flag.")))).
		MustParse(os.Args)
}
