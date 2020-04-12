package main

import (
	"encoding/json"
	cli "github.com/Foxcapades/Argonaut/v0"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"os"
)

type Demo struct {
	SliceVal []string
	IntVal   int
	FloatVal float32
	HexVal   argo.Hex64
}

func main() {
	var demo Demo

	cli.NewCommand().
		// Normal default
		Flag(cli.NewFlag().
			Long("int").
			Bind(&demo.IntVal, true).
			Default(4)).
		// Default from errorless provider
		Flag(cli.NewFlag().
			Long("float").
			Bind(&demo.FloatVal, true).
			Default(func() float32 { return 7.3 })).
		// Default from provider with error
		Flag(cli.NewFlag().
			Long("hex").
			Bind(&demo.HexVal, true).
			Default(func() (argo.Hex64, error) { return argo.Hex64(17), nil })).
		// Default from string (behaves like CLI)
		Flag(cli.NewFlag().
			Long("slice").
			Bind(&demo.SliceVal, true).
			Default("hello")).
		MustParse()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(demo)
}
