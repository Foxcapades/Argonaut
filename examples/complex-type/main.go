package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	cli "github.com/Foxcapades/Argonaut"
)

type Inputs struct {
	Strings       []string
	IntToBool     map[int]bool
	StringToBytes map[string]*[]byte
	Time          time.Time
}

func main() {
	var conf Inputs

	com := cli.Command().
		WithFlag(cli.ShortFlag('t').WithBinding(&conf.Time, true)).
		MustParse(os.Args)

	fmt.Println(com.UnmappedInputs())

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(conf)
}
