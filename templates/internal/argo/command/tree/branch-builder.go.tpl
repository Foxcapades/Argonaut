{{ $vars := (types "argo.BranchCommandBuilder" "CommandBranchBuilder" "argo.BranchCommand") -}}
package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBranchBuilder(name string) argo.BranchCommandBuilder {
	return &CommandBranchBuilder{
		name: name,
	}
}

type CommandBranchBuilder struct {
	{{ template "ParentCommandBuilderProps" $vars }}
	{{ template "ChildCommandBuilderProps" $vars }}
	{{ template "CommandBuilderBaseProps" $vars }}
}

{{ template "ParentCommandBuilderFuncs" $vars }}
{{ template "ChildCommandBuilderFuncs" $vars }}
{{ template "CommandBuilderBaseFuncs" $vars }}
