package flag

import (
	"errors"
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func BuildGroup(g argo.FlagGroupBuilder, o argo.Options) (argo.FlagGroup, error) {
	if g.Size() == 0 {
		return nil, fmt.Errorf("flag group '%s' contains no flags", g.Name())
	}

	errs := xerr.NewMultiError()
	flags := make([]argo.Flag, 0, len(g.Flags()))

	// Ensure the group name is not blank
	if text.IsBlank(g.Name()) {
		errs.AppendError(errors.New("flag group names must not be blank"))
	}

	builders := g.Flags()

	for i := range builders {
		if flag, err := Build(builders[i]); err != nil {
			errs.AppendError(err)
		} else {
			flags = append(flags, flag)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &Group{
		name:  g.Name(),
		desc:  g.Description(),
		flags: flags,
	}, nil

}

func BuildGroups(groups []argo.FlagGroupBuilder, opts argo.Options, errs argo.MultiError) []argo.FlagGroup {
	output := make([]argo.FlagGroup, 0, len(groups))

	UniqueFlagNames(groups, errs)

	for _, builder := range groups {
		if builder.HasFlags() {
			if group, err := BuildGroup(builder, opts); err != nil {
				errs.AppendError(err)
			} else {
				output = append(output, group)
			}
		}
	}

	return output
}
