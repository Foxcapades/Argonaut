package argo

import (
	"bufio"
	"io"
	"slices"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

// CommandLeafHelpRenderer returns a HelpRenderer instance that is suited to
// rendering help text for CommandLeaf instances.
func CommandLeafHelpRenderer() HelpRenderer[CommandLeaf] {
	return comLeafRenderer{}
}

type comLeafRenderer struct{ renderCommandBase }

func (r comLeafRenderer) RenderHelp(leaf CommandLeaf, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return r.renderCommandLeaf(leaf, buf)
	} else {
		buf := bufio.NewWriter(writer)
		err := r.renderCommandLeaf(leaf, buf)
		_ = buf.Flush()
		return err
	}
}

func (r comLeafRenderer) renderCommandLeaf(leaf CommandLeaf, out *bufio.Writer) error {
	if err := r.renderCommandLeafUsage(leaf, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	if leaf.HasAliases() {
		if _, err := out.WriteString(chars.SubLinePadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString("Aliases: "); err != nil {
			return err
		}

		aliases := leaf.Aliases()
		slices.Sort(aliases)

		for i, alias := range aliases {
			if i > 0 {
				if _, err := out.WriteString(", "); err != nil {
					return err
				}
			}
			if _, err := out.WriteString(alias); err != nil {
				return err
			}
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	return r.renderCommandBackHalf(leaf, out)
}

func (r comLeafRenderer) renderCommandLeafUsage(leaf CommandLeaf, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if err := renderSubCommandPath(leaf, out); err != nil {
		return err
	}
	return r.renderCommandUsageBackHalf(leaf, out)
}
