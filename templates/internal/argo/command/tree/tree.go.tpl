{{ $vars := (map "ImplType" "Tree" "OutputType" "argo.CommandTree") -}}
package tree

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Tree struct {
	common.CommandBase[argo.CommandTree]
	{{ template "tree-parent-impl-props" $vars }}
	selectedLeaf argo.LeafCommand
}

func (_ *Tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (i *Tree) SelectedCommand() argo.LeafCommand {
	return i.selectedLeaf
}

{{ template "tree-parent-impl-funcs" $vars }}
