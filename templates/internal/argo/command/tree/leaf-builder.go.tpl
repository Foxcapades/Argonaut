{{ $vars := (types "argo.LeafCommandBuilder" "LeafCommandBuilder" "argo.LeafCommand") -}}
package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewLeafBuilder(name string) argo.LeafCommandBuilder {
	return &LeafCommandBuilder{
		name: name,
	}
}

type LeafCommandBuilder struct {
	{{ template "ChildCommandBuilderProps" $vars }}

	{{ template "CommandBuilderBaseProps" $vars }}

	{{ template "CommandEndBuilderBaseProps" $vars }}
}

{{ template "ChildCommandBuilderFuncs" $vars }}

{{ template "CommandBuilderBaseFuncs" $vars }}

{{ template "CommandEndBuilderBaseFuncs" $vars}}
