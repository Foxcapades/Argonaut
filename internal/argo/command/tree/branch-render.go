package tree

import (
	"fmt"
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderBranchHelp(branch argo.BranchCommand, options Options, writer io.Writer) error {
	if branch.HasSelectedChild() {
		switch child := branch.SelectedChild().(type) {
		case argo.LeafCommand:
			return RenderLeafHelp(child, options, writer)
		case argo.BranchCommand:
			return RenderBranchHelp(child, options, writer)
		default:
			panic(fmt.Sprintf("invalid state: unrecognized argo.ChildNode type %T", branch.SelectedChild()))
		}
	}

	if branch.IsHelpDisabled() {
		return nil
	}

	buf := utils.NewBatchWriter(writer)
	renderCommandBranch(branch, options, buf)
	buf.Flush()
	return buf.Error
}

func MakeRenderBranchHelpCallback(branch argo.BranchCommand, options Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderBranchHelp(branch, options, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommandBranch(branch argo.BranchCommand, options Options, out *utils.BatchWriter) {
	renderCommandBranchUsage(branch, out)

	out.WriteByte(text.LineFeedByte)

	TryRenderAliases(branch, out)

	TryRenderDescription(branch, options, out, branch.HasAliases())

	RenderCommandGroups(branch.CommandGroups(), options, 0, out)

	TryRenderFlags(branch, options, out)

	TryRenderInheritedFlags(branch, options, out)

	out.WriteByte(text.LineFeedByte)
}

func renderCommandBranchUsage(node argo.BranchCommand, out *utils.BatchWriter) {
	out.WriteString(common.CommandRenderPrefix)
	RenderSubCommandPath(node, out)

	// Possibly render "[options]" if there are any optional flags
	for _, f := range flag.Stream(node) {
		if !f.IsRequired() {
			out.WriteString(common.CommandRenderOptionalFlags)
			break
		}
	}

	out.WriteString(subcommandPlaceholder)
}
