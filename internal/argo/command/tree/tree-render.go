package tree

import (
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderHelp(tree argo.TreeCommand, opts Options, writer io.Writer) error {
	buf := utils.NewBatchWriter(writer)
	renderCommandTree(tree, opts, buf)
	buf.Flush()
	return buf.Error
}

func MakeRenderTreeHelpCallback(tree argo.TreeCommand, opts Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderHelp(tree, opts, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommandTree(tree argo.TreeCommand, opts Options, out *utils.BatchWriter) {
	renderTreeUsageBlock(tree, out)

	out.WriteByte(text.LineFeedByte)

	TryRenderDescription(tree, opts, out, false)

	RenderCommandGroups(tree.CommandGroups(), opts, 0, out)

	TryRenderFlags(tree, opts, out)

	out.WriteByte(text.LineFeedByte)
}

func renderTreeUsageBlock(tree argo.TreeCommand, out *utils.BatchWriter) {
	out.WriteString(common.CommandRenderPrefix)
	out.WriteString(render.SubLinePadding[0])
	out.WriteString(tree.Name())

	if tree.HasFlagGroups() {
		hasOptionalFlags := false

		for _, f := range flag.Stream(tree) {
			if f.IsRequired() {
				out.WriteByte(text.SpaceByte)
				flag.RenderShortestForUsage(f, out)
			} else {
				hasOptionalFlags = true
			}
		}

		if hasOptionalFlags {
			out.WriteString(common.CommandRenderOptionalFlags)
		}
	}

	out.WriteString(subcommandPlaceholder)
}
