package main

import (
	"encoding/json"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type Inputs struct {
	Hex   []argo.Hex
	UHex  map[string]argo.UHex
	Octal []argo.Octal
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
