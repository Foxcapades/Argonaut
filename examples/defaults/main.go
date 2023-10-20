package main

import (
	"encoding/json"
	"os"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type Demo struct {
	SliceVal []string
	IntVal   int
	FloatVal float32
	HexVal   argo.Hex64
}

func main() {
	var demo Demo

	cli.Command().
		// Normal default
		WithFlag(cli.Flag().
			WithLongForm("int").
			WithArgument(cli.Argument().
				WithBinding(&demo.IntVal).
				Require().
				WithDefault(4))).
		// Default from errorless provider
		WithFlag(cli.Flag().
			WithLongForm("float").
			WithArgument(cli.Argument().
				WithBinding(&demo.FloatVal).
				Require().
				WithDefault(func() float32 { return 7.3 }))).
		// Default from provider with error
		WithFlag(cli.Flag().
			WithLongForm("hex").
			WithArgument(cli.Argument().
				WithBinding(&demo.HexVal).
				Require().
				WithDefault(func() (argo.Hex64, error) { return argo.Hex64(17), nil }))).
		// Default from string (behaves like CLI)
		WithFlag(cli.Flag().
			WithLongForm("slice").
			WithArgument(cli.Argument().
				WithBinding(&demo.SliceVal).
				WithDefault("hello"))).
		MustParse(os.Args)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(demo)
}
