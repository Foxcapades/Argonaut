package tree

import (
	"bufio"
	"io"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderLeafHelp(leaf argo.LeafCommand, options argo.Options, writer io.Writer) error {
	if leaf.IsHelpDisabled() {
		return nil
	}

	var buf *bufio.Writer
	defer utils.DisregardError(buf.Flush)

	if b, ok := writer.(*bufio.Writer); ok {
		buf = b
	} else {
		buf = bufio.NewWriter(writer)
	}

	return renderCommandLeaf(leaf, options, buf)
}

func renderCommandLeaf(leaf argo.LeafCommand, options argo.Options, out *bufio.Writer) error {
	if err := renderCommandLeafUsage(leaf, out); err != nil {
		return err
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	if err := tryRenderAliases(leaf, out); err != nil {
		return err
	}

	return renderCommandLeafBackHalf(leaf, options, out)
}

func renderCommandLeafUsage(leaf argo.LeafCommand, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if err := renderSubCommandPath(leaf, out); err != nil {
		return err
	}
	return common.RenderCommandUsageBackHalf(leaf, out)
}

func renderCommandLeafBackHalf(com argo.LeafCommand, options argo.Options, out *bufio.Writer) error {

	// If the command has a description, append it.
	if com.HasDescription() {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[0], options.HelpTextMaxWidth, out)
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
		if err := flag.RenderGroups(com.FlagGroups(), options, 0, out); err != nil {
			return err
		}
	}

	inherited := flag.FlattenInheritance(com)
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
			if err := flag.RenderInherited(&inherited[i], options, 1, out); err != nil {
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
		if _, err := out.WriteString(common.CommandRenderArgs); err != nil {
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
				if err := argument.Render(arg, options, 1, out, i+1); err != nil {
					return err
				}
			} else {
				if err := argument.Render(arg, options, 1, out, 0); err != nil {
					return err
				}
			}
		}
	}

	return out.WriteByte(chars.CharLF)
}
