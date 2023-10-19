package argo

import (
	"bufio"
	"slices"
)

const (
	subcommandPlaceholder = " <command>"
)

const (
	comPrefix = "Usage:\n"
	comOpts   = " [options]"
	comArgs   = "Arguments"
)

func renderSubCommandPath(node CommandNode, out *bufio.Writer) error {
	path := make([]string, 0, 4)

	current := node
	for current != nil {
		path = append(path, current.Name())
		current = current.Parent()
	}

	slices.Reverse(path)

	if _, err := out.WriteString(subLinePadding[0]); err != nil {
		return err
	}

	for i, segment := range path {
		if i > 0 {
			if err := out.WriteByte(charSpace); err != nil {
				return err
			}
		}
		if _, err := out.WriteString(segment); err != nil {
			return err
		}
	}

	return nil
}

type renderCommandBase struct{ renderBase }

func (r renderCommandBase) renderCommandBackHalf(com Command, out *bufio.Writer) error {

	// If the command has a description, append it.
	if com.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(com.Description(), descriptionPadding[0], helpTextMaxWidth, out); err != nil {
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
			if err := out.WriteByte(charLF); err != nil {
				return err
			}
		}

		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := r.renderFlagGroups(com.FlagGroups(), 0, out); err != nil {
			return err
		}
	}

	if writeArgs {
		if _, err := out.WriteString(paragraphBreak); err != nil {
			return err
		}
		if _, err := out.WriteString(headerPadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString(comArgs); err != nil {
			return err
		}

		for i, arg := range com.Arguments() {
			if i > 0 {
				if err := out.WriteByte(charLF); err != nil {
					return err
				}
			}
			if err := out.WriteByte(charLF); err != nil {
				return err
			}
			if err := r.renderArgument(arg, 1, out); err != nil {
				return err
			}
		}
	}

	return out.WriteByte(charLF)
}

func (r renderCommandBase) renderCommandUsageBackHalf(com Command, out *bufio.Writer) error {
	// If the command has flag groups
	if com.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range com.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					if err := out.WriteByte(charSpace); err != nil {
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

	// After all the flag groups have been rendered, append the argument names.
	if com.HasArguments() {
		for _, arg := range com.Arguments() {
			if err := out.WriteByte(charSpace); err != nil {
				return err
			}
			if err := r.renderArgumentName(arg, out); err != nil {
				return err
			}
		}
	}

	if com.HasUnmappedLabel() {
		if err := out.WriteByte(charSpace); err != nil {
			return err
		}
		if err := out.WriteByte(argOptPrefix); err != nil {
			return err
		}
		if _, err := out.WriteString(com.GetUnmappedLabel()); err != nil {
			return err
		}
		if err := out.WriteByte(argOptSuffix); err != nil {
			return err
		}
	}

	return nil
}
