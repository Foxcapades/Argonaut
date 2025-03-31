package command

import (
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderHelp(command argo.Command, opts opts.Options, writer io.Writer) error {
	buf := utils.NewBatchWriter(writer)
	renderCommand(command, opts, buf)
	buf.Flush()
	return buf.Error
}

func MakeRenderHelpCallback(command argo.Command, opts opts.Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderHelp(command, opts, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommand(com argo.Command, opts opts.Options, out *utils.BatchWriter) {
	renderCommandUsageBlock(com, out)
	out.WriteByte(text.LineFeedByte)
	common.RenderCommandBackHalf(com, opts, out)
}

func renderCommandUsageBlock(com argo.Command, out *utils.BatchWriter) {
	out.WriteString(common.CommandRenderPrefix)
	out.WriteString(render.SubLinePadding[0])
	out.WriteString(com.Name())
	// If the command has flag groups
	if com.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, f := range flag.Stream(com) {
			if f.IsRequired() {
				out.WriteByte(text.SpaceByte)
				flag.RenderShortestForUsage(f, out)
			} else {
				hasOptionalFlags = true
			}
		}

		// If there are any optional flags append the general "[OPTIONS]" text.
		if hasOptionalFlags {
			out.WriteString(common.CommandRenderOptionalFlags)
		}
	}

	// After all the flag groups have been rendered, append the argument names.
	argument.RenderForUsageLine(com.Arguments(), out)

	if com.HasUnmappedInputLabel() {
		out.WriteByte(text.SpaceByte)
		out.WriteByte(argument.OptPrefix)
		out.WriteString(com.UnmappedInputLabel())
		out.WriteByte(argument.OptSuffix)
	}
}
