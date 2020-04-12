package main

import (
	"encoding/json"
	"os"

	"github.com/Foxcapades/Argonaut/v0"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type Inputs struct {
	Hex   []argo.Hex
	UHex  map[string]argo.UHex
	Octal []argo.Octal
}

func main() {
	var conf Inputs

	cli.NewCommand().
		Flag(cli.NewFlag().
			Long("hex").
			Short('x').
			Description("Hex value").
			Bind(&conf.Hex, true)).
		Flag(cli.NewFlag().
			Long("uhex").
			Short('u').
			Description("Unsigned hex value").
			Bind(&conf.UHex, true)).
		Flag(cli.NewFlag().
			Long("octal").
			Short('o').
			Description("Octal value").
			Bind(&conf.Octal, true)).
		MustParse()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
