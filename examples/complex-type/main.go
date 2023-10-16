package main

import (
	"encoding/json"
	"os"

	cli "github.com/Foxcapades/Argonaut"
)

type Inputs struct {
	Strings       []string
	IntToBool     map[int]bool
	StringToBytes map[string]*[]byte
}

func main() {
	var conf Inputs

	cli.Command().
		WithFlag(cli.Flag().
			WithLongForm("string-slice").
			WithShortForm('s').
			WithBinding(&conf.Strings, true)).
		WithFlag(cli.Flag().
			WithLongForm("int-bool-map").
			WithShortForm('i').
			WithBinding(&conf.IntToBool, true)).
		WithFlag(cli.Flag().
			WithLongForm("string-bytes").
			WithShortForm('b').
			WithBinding(&conf.StringToBytes, true)).
		MustParse(os.Args)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
