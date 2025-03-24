package command

import (
	"os"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Parse(com argo.Command) (argo.ParseResult, error) {
	return NewInterpreter(os.Args, com).Run()
}
