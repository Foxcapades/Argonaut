package argo

import (
	"bufio"
	"io"
	"slices"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

// CommandBranchHelpRenderer returns a HelpRenderer instance that is suited to
// rendering help text for CommandBranch instances.
func CommandBranchHelpRenderer() HelpRenderer[CommandBranch] {
	return comBranchRenderer{}
}

type comBranchRenderer struct{ renderBase }

func (r comBranchRenderer) RenderHelp(branch CommandBranch, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return r.renderCommandBranch(branch, buf)
	} else {
		buf := bufio.NewWriter(writer)
		err := r.renderCommandBranch(branch, buf)
		_ = buf.Flush()
		return err
	}
}

func (r comBranchRenderer) renderCommandBranch(branch CommandBranch, out *bufio.Writer) error {
	if err := r.renderCommandBranchUsage(branch, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	ha := branch.HasAliases()
	hd := branch.HasDescription()
	hf := branch.HasFlagGroups()

	if ha {
		if _, err := out.WriteString(chars.SubLinePadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString("Aliases: "); err != nil {
			return err
		}

		aliases := branch.Aliases()
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

	if hd {
		if ha {
			if err := out.WriteByte(chars.CharLF); err != nil {
				return err
			}
		}
		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[0], chars.HelpTextMaxWidth, out)
		if err := formatter.Format(branch.Description()); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	// render flags
	if hf {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := r.renderFlagGroups(branch.FlagGroups(), 0, out); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	if err := renderCommandGroups(branch.CommandGroups(), 0, out); err != nil {
		return err
	}

	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}

	return nil
}

func (r comBranchRenderer) renderCommandBranchUsage(node CommandBranch, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if err := renderSubCommandPath(node, out); err != nil {
		return err
	}

	if node.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range node.FlagGroups() {
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

		// If there are any optional flags append the general "[OPTIONS]" text.
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
