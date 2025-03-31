package tree

import (
	"bufio"
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
	if buf, ok := writer.(*bufio.Writer); ok {
		return renderCommandTree(tree, opts, buf)
	}

	buf := bufio.NewWriter(writer)
	err := renderCommandTree(tree, opts, buf)
	_ = buf.Flush()
	return err
}

func MakeRenderTreeHelpCallback(tree argo.TreeCommand, opts Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderHelp(tree, opts, os.Stdout))
		utils.Exit(0)
	}
}

func renderCommandTree(tree argo.TreeCommand, opts Options, out *bufio.Writer) error {
	if err := renderTreeUsageBlock(tree, out); err != nil {
		return err
	}
	if err := out.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	if err := TryRenderDescription(tree, opts, out, false); err != nil {
		return err
	}

	if err := RenderCommandGroups(tree.CommandGroups(), opts, 0, out); err != nil {
		return err
	}

	if err := TryRenderFlags(tree, opts, out); err != nil {
		return err
	}

	if err := out.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	return nil
}

func renderTreeUsageBlock(tree argo.TreeCommand, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(render.SubLinePadding[0]); err != nil {
		return err
	}
	if _, err := out.WriteString(tree.Name()); err != nil {
		return err
	}

	if tree.HasFlagGroups() {
		hasOptionalFlags := false

		for _, group := range tree.FlagGroups() {
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
