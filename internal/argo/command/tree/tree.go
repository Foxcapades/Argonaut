package tree

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandTree struct {
	parent[argo.CommandTree]

	warnings *argo.WarningContext

	selectedLeaf argo.LeafCommand
}

func (_ *commandTree) Name() string {
	return filepath.Base(os.Args[0])
}

func (t *commandTree) Warnings() []string {
	return t.warnings.GetWarnings()
}

func (t *commandTree) AppendWarning(warning string) {
	t.warnings.AppendWarning(warning)
}

func (t *commandTree) SelectedCommand() argo.LeafCommand {
	return t.selectedLeaf
}
