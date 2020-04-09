package render

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

func Flag(f A.Flag) string {
	var out strings.Builder
	flag(f, &out)
	return out.String()
}

func flag(f A.Flag, out *strings.Builder) {

}