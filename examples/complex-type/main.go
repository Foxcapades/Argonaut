package main

import (
	"encoding/json"
	cli "github.com/Foxcapades/Argonaut/v0"
	"os"
)

type Inputs struct {
	Strings       []string
	IntToBool     map[int]bool
	StringToBytes map[string]*[]byte
}

func main() {
	var conf Inputs

	cli.NewCommand().
		Flag(cli.NewFlag().
			Long("string-slice").
			Short('s').
			Bind(&conf.Strings, true)).
		Flag(cli.NewFlag().
			Long("int-bool-map").
			Short('i').
			Bind(&conf.IntToBool, true)).
		Flag(cli.NewFlag().
			Long("string-bytes").
			Short('b').
			Bind(&conf.StringToBytes, true)).
		MustParse()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
