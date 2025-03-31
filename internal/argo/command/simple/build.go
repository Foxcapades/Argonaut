package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
	"github.com/foxcapades/argonaut/v3/pkg/argoutil"
)

func Build(builder argo.CommandBuilder) (argo.Command, error) {
	options := FixOptions(builder.Options())
	comOpts := WrapOptions(&options)

	errs := xerr.NewMultiError()
	com := &Command{
		helpDisabled:  builder.IsHelpDisabled(),
		description:   builder.Description(),
		unmappedLabel: builder.UnmappedInputLabel(),
		options:       options,
		callback:      builder.Callback(),
		arguments:     argument.BuildMulti(builder.Arguments(), errs), // filled below

		// flagGroups:    nil, // filled below

		// unmapped:      nil, // filled on parse
	}

	if !builder.IsHelpDisabled() {
		common.TryAddHelpFlags(builder, MakeRenderHelpCallback(com, comOpts), comOpts)
	}

	com.flagGroups = flag.BuildGroups(builder.FlagGroups(true), errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return com, nil
}

func FixOptions(opts argo.CommandOptions) argo.CommandOptions {
	consoleWidth, _ := argoutil.GetConsoleWidth()

	if opts.HelpTextMaxWidth < 60 {
		opts.HelpTextMaxWidth = max(min(consoleWidth, 120), 60)
	}

	if len(opts.MetaFlagGroupName) == 0 {
		opts.MetaFlagGroupName = "General Flags"
	}

	return opts
}
