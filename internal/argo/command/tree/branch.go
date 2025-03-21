package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandBranch struct {
	parent[argo.BranchCommand]
	child[argo.BranchCommand]

	warnings *argo.WarningContext
}

func (c commandBranch) Warnings() []string {
	return c.warnings.GetWarnings()
}

func (c commandBranch) AppendWarning(warning string) {
	c.warnings.AppendWarning(warning)
}
