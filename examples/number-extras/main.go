package main

import (
	"encoding/json"
	"os"

	cli "github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/pkg/argotype"
)

type Inputs struct {
	Hex   []argotype.Hex
	UHex  map[string]argotype.UHex
	Octal []argotype.Octal
}

func main() {
	var conf Inputs

	cli.Command().
		WithFlag(cli.Flag().
			WithLongForm("hex").
			WithShortForm('x').
			WithDescription("Hex value").
			WithBinding(&conf.Hex, true)).
		WithFlag(cli.Flag().
			WithLongForm("uhex").
			WithShortForm('u').
			WithDescription("Unsigned hex value").
			WithBinding(&conf.UHex, true)).
		WithFlag(cli.Flag().
			WithLongForm("octal").
			WithShortForm('o').
			WithDescription("Octal value").
			WithBinding(&conf.Octal, true)).
		MustParse(os.Args)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
