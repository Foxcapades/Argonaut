package tree

import (
	"bufio"
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

	if buf, ok := writer.(*bufio.Writer); ok {
		return renderCommandBranch(branch, options, buf)
	}

	buf := bufio.NewWriter(writer)
	err := renderCommandBranch(branch, options, buf)
	_ = buf.Flush()
	return err
}

func MakeRenderBranchHelpCallback(branch argo.BranchCommand, options Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderBranchHelp(branch, options, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommandBranch(branch argo.BranchCommand, options Options, out *bufio.Writer) error {
	if err := renderCommandBranchUsage(branch, out); err != nil {
		return err
	}

	if err := out.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	if err := TryRenderAliases(branch, out); err != nil {
		return err
	}

	if err := TryRenderDescription(branch, options, out, branch.HasAliases()); err != nil {
		return err
	}

	if err := RenderCommandGroups(branch.CommandGroups(), options, 0, out); err != nil {
		return err
	}

	if err := TryRenderFlags(branch, options, out); err != nil {
		return err
	}

	if err := TryRenderInheritedFlags(branch, options, out); err != nil {
		return err
	}

	if err := out.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	return nil
}

func renderCommandBranchUsage(node argo.BranchCommand, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if err := RenderSubCommandPath(node, out); err != nil {
		return err
	}

	if node.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range node.FlagGroups() {
			for _, f := range group.Flags() {
				if f.IsRequired() {
					if err := out.WriteByte(text.SpaceByte); err != nil {
						return err
					}
					if err := flag.RenderShortestForUsage(f, out); err != nil {
						return err
					}
				} else {
					hasOptionalFlags = true
				}
			}
		}

		// If there are any optional flags append the general "[OPTIONS]" text.
		if hasOptionalFlags {
			if _, err := out.WriteString(common.CommandRenderOptionalFlags); err != nil {
				return err
			}
		}
	}

	if _, err := out.WriteString(subcommandPlaceholder); err != nil {
		return err
	}

	return nil
}
