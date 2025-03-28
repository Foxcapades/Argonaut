package main

import (
	"encoding/json"
	"os"

	"github.com/foxcapades/argonaut/v3"
)

type Inputs struct {
	Strings       []string
	IntToBool     map[int]bool
	StringToBytes map[string][]byte
}

func main() {
	var conf Inputs

	com := cli.MustBuildCommand(cli.Command().
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
			WithBinding(&conf.StringToBytes, true)))

	_ = cli.MustParse(com)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(conf)
}
