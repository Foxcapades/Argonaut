package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func BuildGroup(group argo.CommandGroupBuilder, options Options, parent argo.ParentNode) (argo.CommandGroup, error) {
	errs := xerr.NewMultiError()
	out := CommandGroup{}

	// Ensure that the name is not blank
	if text.IsBlank(group.Name()) {
		errs.AppendError(errors.New("command group names must not be blank"))
	}

	out.name = group.Name()
	out.description = group.Description()

	if group.HasBranches() {
		branches := group.Branches()
		out.branches = make([]argo.BranchCommand, len(branches))

		for i, builder := range branches {
			if branch, err := BuildBranch(builder, options, parent); err != nil {
				errs.AppendError(err)
			} else {
				out.branches[i] = branch
			}
		}
	}

	if group.HasLeaves() {
		leaves := group.Leaves()
		out.leaves = make([]argo.LeafCommand, len(leaves))

		for i, builder := range leaves {
			if leaf, err := BuildLeaf(builder, options, parent); err != nil {
				errs.AppendError(err)
			} else {
				out.leaves[i] = leaf
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return out, nil
}
