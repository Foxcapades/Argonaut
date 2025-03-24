package tree

import (
	"bufio"
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderHelp(tree argo.TreeCommand, options argo.Options, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return renderCommandTree(tree, options, buf)
	}

	buf := bufio.NewWriter(writer)
	err := renderCommandTree(tree, options, buf)
	_ = buf.Flush()
	return err
}

func MakeRenderTreeHelpCallback(tree argo.TreeCommand, options argo.Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderHelp(tree, options, os.Stdout))
		os.Exit(0)
	}
}

func renderCommandTree(tree argo.TreeCommand, options argo.Options, out *bufio.Writer) error {
	if err := renderTreeUsageBlock(tree, options, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	hd := tree.HasDescription()
	hf := tree.HasFlagGroups()

	if hd {
		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[0], options.HelpTextMaxWidth, out)
		if err := formatter.Format(tree.Description()); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	if hf {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := flag.RenderGroups(tree.FlagGroups(), 0, out); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}
	if err := RenderCommandGroups(tree.CommandGroups(), options, 0, out); err != nil {
		return err
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	return nil
}

func renderTreeUsageBlock(tree argo.TreeCommand, options argo.Options, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(chars.SubLinePadding[0]); err != nil {
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
					if err := out.WriteByte(chars.CharSpace); err != nil {
						return err
					}
					if err := flag.RenderShortestLine(f, out); err != nil {
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
