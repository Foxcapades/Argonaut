package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func BuildLeaf(l argo.LeafCommandBuilder, options argo.Options, parent argo.ParentNode) (argo.LeafCommand, error) {
	errs := xerr.NewMultiError()

	// Ensure the group name is not blank
	xerr.AppendIfPresent(ValidateNodeName(l.Name()), errs)

	// Ensure the aliases are all not blank
	for _, alias := range l.Aliases() {
		if text.IsBlank(alias) {
			errs.AppendError(errors.New("command leaf aliases must not be blank"))
		}
	}

	leaf := Leaf{
		name:          l.Name(),
		description:   l.Description(),
		aliases:       l.Aliases(),
		parent:        parent,
		callback:      l.Callback(),
		unmappedLabel: l.UnmappedInputLabel(),
		args:          argument.BuildMulti(l.Arguments(), errs),
	}

	if !l.IsHelpDisabled() {
		common.TryAddHelpFlags[argo.LeafCommandBuilder](l, MakeRenderLeafHelpCallback(&leaf, options), options)
	}
	// this needs to happen _after_ TryAddHelpFlags, which itself has to happen
	// after the leaf instance is created.
	leaf.flagGroups = flag.BuildGroups(l.FlagGroups(true), options, errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &leaf, nil
}
