package main

import (
	"encoding/json"
	"os"

	cli "github.com/foxcapades/argonaut/v3"
	"github.com/foxcapades/argonaut/v3/pkg/argotype"
)

type Demo struct {
	SliceVal []string
	IntVal   int
	FloatVal float32
	HexVal   argotype.Hex64
}

func main() {
	var demo Demo

	cli.MustParse(cli.Command().
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
				WithDefault(func() (argotype.Hex64, error) { return argotype.Hex64(17), nil }))).
		// Default from string (behaves like CLI)
		WithFlag(cli.Flag().
			WithLongForm("slice").
			WithArgument(cli.Argument().
				WithBinding(&demo.SliceVal).
				WithDefault("hello"))))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(demo)
}
