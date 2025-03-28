{{ $vars := (types "argo.TreeCommandBuilder" "TreeCommandBuilder" "argo.TreeCommand") -}}
package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.TreeCommandBuilder {
	return new(TreeCommandBuilder)
}

type TreeCommandBuilder struct {
	{{ template "ParentCommandBuilderProps" $vars }}
	{{ template "CommandBuilderBaseProps" $vars }}
}

{{ template "ParentCommandBuilderFuncs" $vars }}
{{ template "CommandBuilderBaseFuncs" $vars }}
