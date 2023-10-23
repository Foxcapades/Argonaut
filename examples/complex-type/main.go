package main

import (
	"encoding/json"
	"os"

	cli "github.com/Foxcapades/Argonaut"
)

type Inputs struct {
	Map map[string][]string
}

func main() {
	var conf Inputs

	cli.Command().
		WithFlag(cli.Flag().
			WithShortForm('s').
			WithBinding(&conf.Map, true)).
		MustParse(os.Args)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
