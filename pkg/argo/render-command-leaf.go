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

	return renderCommandLeafBackHalf(leaf, out)
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

func renderCommandLeafBackHalf(com CommandLeaf, out *bufio.Writer) error {

	// If the command has a description, append it.
	if com.HasDescription() {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[0], chars.HelpTextMaxWidth, out)
		if err := formatter.Format(com.Description()); err != nil {
			return err
		}
	}

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

	if com.HasFlagGroups() {
		if com.HasDescription() {
			if err := out.WriteByte(chars.CharLF); err != nil {
				return err
			}
		}

		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := renderFlagGroups(com.FlagGroups(), 0, out); err != nil {
			return err
		}
	}

	inherited := flattenFlagInheritance(com)
	if len(inherited) > 0 {
		if com.HasFlagGroups() {
			if err := out.WriteByte(chars.CharLF); err != nil {
				return err
			}
		}

		if _, err := out.WriteString("\nInherited Flags"); err != nil {
			return err
		}

		for i := range inherited {
			if err := out.WriteByte(chars.CharLF); err != nil {
				return err
			}
			if err := renderInheritedFlag(&inherited[i], 1, out); err != nil {
				return err
			}
		}
	}

	if writeArgs {
		if _, err := out.WriteString(chars.ParagraphBreak); err != nil {
			return err
		}
		if _, err := out.WriteString(chars.HeaderPadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString(comArgs); err != nil {
			return err
		}

		multiArgs := len(com.Arguments()) > 1

		for i, arg := range com.Arguments() {
			if i > 0 {
				if err := out.WriteByte(chars.CharLF); err != nil {
					return err
				}
			}
			if err := out.WriteByte(chars.CharLF); err != nil {
				return err
			}
			if multiArgs {
				if err := renderArgument(arg, 1, out, i+1); err != nil {
					return err
				}
			} else {
				if err := renderArgument(arg, 1, out, 0); err != nil {
					return err
				}
			}
		}
	}

	return out.WriteByte(chars.CharLF)
}
