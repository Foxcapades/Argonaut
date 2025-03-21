package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.CommandBuilder) (argo.Command, error) {
	errs := xerr.NewMultiError()
	com := new(command)

	var flagGroups []argo.FlagGroupBuilder
	if !builder.IsHelpDisabled() {
		flagGroups = common.ConfigureHelpFlags[argo.Command](builder, render.CommandHelpRenderer(), com)
	} else {
		flagGroups = builder.FlagGroups(true)
	}

	com.flagGroups = make([]argo.FlagGroup, 0, len(flagGroups))

	flag.UniqueFlagNames(flagGroups, errs)

	for _, builder := range flagGroups {
		if builder.HasFlags() {
			if group, err := flag.BuildGroup(builder); err != nil {
				errs.AppendError(err)
			} else {
				com.flagGroups = append(com.flagGroups, group)
			}
		}
	}

	if builder.HasArguments() {
		forceRequiredUntil := 0
		for i, argBuilder := range builder.Arguments() {
			if argBuilder.IsRequired() {
				forceRequiredUntil = i
			}
		}

		com.arguments = make([]argo.Argument, 0, len(builder.Arguments()))
		for i, builder := range builder.Arguments() {

			if i < forceRequiredUntil && !builder.IsRequired() {
				builder.Require()
				// TODO: should this be an error?
				// ctx.AppendWarning(fmt.Sprintf("argument %d was not marked as required, but preceded required argument %d", i+1, forceRequiredUntil+1))
			}

			if arg, err := argument.Build(builder); err != nil {
				errs.AppendError(err)
			} else {
				com.arguments = append(com.arguments, arg)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	com.description = builder.Description()
	com.unmappedLabel = builder.UnmappedInputLabel()
	com.callback = builder.Callback()

	return com, nil
}
