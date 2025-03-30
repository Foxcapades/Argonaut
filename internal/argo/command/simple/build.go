package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.CommandBuilder, options argo.Options) (argo.Command, error) {
	opts.FixOptions(&options)

	errs := xerr.NewMultiError()
	com := new(Command)

	if !builder.IsHelpDisabled() {
		common.TryAddHelpFlags(builder, MakeRenderHelpCallback(com, options), options)
	}

	com.flagGroups = flag.BuildGroups(builder.FlagGroups(true), options, errs)

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
