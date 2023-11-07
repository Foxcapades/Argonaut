package argo

import (
	"bufio"
	"io"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

// CommandTreeHelpRenderer returns a HelpRenderer instance that is suited to
// rendering help text for CommandTree instances.
func CommandTreeHelpRenderer() HelpRenderer[CommandTree] {
	return comTreeRenderer{}
}

type comTreeRenderer struct{ renderBase }

func (r comTreeRenderer) RenderHelp(tree CommandTree, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return r.renderCommandTree(tree, buf)
	} else {
		buf := bufio.NewWriter(writer)
		err := r.renderCommandTree(tree, buf)
		_ = buf.Flush()
		return err
	}
}

func (r comTreeRenderer) renderCommandTree(tree CommandTree, out *bufio.Writer) error {
	if err := r.renderCommandTreeUsageBlock(tree, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	hd := tree.HasDescription()
	hf := tree.HasFlagGroups()

	if hd {
		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[0], chars.HelpTextMaxWidth, out)
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
		if err := renderFlagGroups(tree.FlagGroups(), 0, out); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}
	if err := renderCommandGroups(tree.CommandGroups(), 0, out); err != nil {
		return err
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	return nil
}

func (r comTreeRenderer) renderCommandTreeUsageBlock(tree CommandTree, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
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
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					if err := out.WriteByte(chars.CharSpace); err != nil {
						return err
					}
					if err := r.renderShortestFlagLine(flag, out); err != nil {
						return err
					}
				} else {
					hasOptionalFlags = true
				}
			}
		}

		if hasOptionalFlags {
			if _, err := out.WriteString(comOpts); err != nil {
				return err
			}
		}
	}

	if _, err := out.WriteString(subcommandPlaceholder); err != nil {
		return err
	}

	return nil
}
