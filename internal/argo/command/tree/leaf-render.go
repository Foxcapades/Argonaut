package tree

import (
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderLeafHelp(leaf argo.LeafCommand, options Options, writer io.Writer) error {
	if leaf.IsHelpDisabled() {
		return nil
	}

	buf := utils.NewBatchWriter(writer)
	renderCommandLeaf(leaf, options, buf)
	buf.Flush()
	return buf.Error
}

func MakeRenderLeafHelpCallback(leaf argo.LeafCommand, options Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderLeafHelp(leaf, options, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommandLeaf(leaf argo.LeafCommand, options Options, out *utils.BatchWriter) {
	renderCommandLeafUsage(leaf, out)
	out.WriteByte(text.LineFeedByte)
	TryRenderAliases(leaf, out)
	renderCommandLeafBackHalf(leaf, options, out)
}

func renderCommandLeafUsage(leaf argo.LeafCommand, out *utils.BatchWriter) {
	out.WriteString(common.CommandRenderPrefix)
	RenderSubCommandPath(leaf, out)

	// Possibly render "[options]" if there are any optional flags
	for _, f := range flag.Stream(leaf) {
		if !f.IsRequired() {
			out.WriteString(common.CommandRenderOptionalFlags)
			break
		}
	}

	argument.RenderForUsageLine(leaf.Arguments(), out)

	if leaf.HasUnmappedInputLabel() {
		out.WriteByte(argument.OptPrefix)
		out.WriteString(leaf.UnmappedInputLabel())
		out.WriteByte(argument.OptSuffix)
	}
}

func renderCommandLeafBackHalf(com argo.LeafCommand, options Options, out *utils.BatchWriter) {

	// If the command has a description, append it.
	TryRenderDescription(com, options, out, com.HasAliases())

	// Figure out if we have any printable arguments.
	//
	// Showing the arguments is conditional based on whether any of the arguments
	// have a description value.  If none do, then there is no value in rendering
	// them as they already appear in the usage line.
	//
	// This is calculated ahead of time as it informs whether the flag group
	// headers should be printed.
	writeArgs := false
	if com.HasArguments() {
		for _, arg := range com.Arguments() {
			if arg.HasDescription() {
				writeArgs = true
			}
		}
	}

	TryRenderFlags(com, options, out)

	TryRenderInheritedFlags(com, options, out)

	if writeArgs {
		out.WriteString(render.ParagraphBreak)
		out.WriteString(render.HeaderPadding[0])
		out.WriteString(common.CommandRenderArgs)

		multiArgs := len(com.Arguments()) > 1

		for i, arg := range com.Arguments() {
			if i > 0 {
				out.WriteByte(text.LineFeedByte)
			}
			out.WriteByte(text.LineFeedByte)
			if multiArgs {
				argument.Render(arg, options, 1, out, i+1)
			} else {
				argument.Render(arg, options, 1, out, 0)
			}
		}
	}

	out.WriteByte(text.LineFeedByte)
}
