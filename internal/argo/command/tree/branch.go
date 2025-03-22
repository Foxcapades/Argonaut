package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandBranch struct {
	parent[argo.BranchCommand]
	child[argo.BranchCommand]
}
