package render

import (
	"bufio"
	"github.com/Foxcapades/Argonaut/v3/internal/chars"
	cli_tree "github.com/Foxcapades/Argonaut/v3/pkg/argo/command/tree"
	cli_flag "github.com/Foxcapades/Argonaut/v3/pkg/argo/flag"
)

const (
	argReqPrefix = '<'
	argReqSuffix = '>'
	argOptPrefix = '['
	argOptSuffix = ']'
)

type renderBase struct{}
