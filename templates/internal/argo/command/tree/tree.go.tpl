{{ $vars := (map "ImplType" "Tree" "OutputType" "argo.CommandTree") -}}
package tree

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Tree struct {
	{{ template "CommandBaseProps" $vars }}
	{{ template "ParentCommonProps" $vars }}
	selectedLeaf argo.LeafCommand
}

func (_ *Tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (i *Tree) SelectedCommand() argo.LeafCommand {
	return i.selectedLeaf
}

{{ template "ParentCommonFuncs" $vars }}

{{ template "CommandBaseFuncs" $vars }}
