package argo

import (
	"bufio"
	"slices"
)

const (
	comPrefix = "Usage:\n"
	comOpts   = " [options]"
	comArgs   = "Positional Arguments"
)

func renderCommand(com Command, out *bufio.Writer) error {
	if err := renderCommandUsageBlock(com, out); err != nil {
		return err
	}
	if err := out.WriteByte(charLF); err != nil {
		return err
	}
	return renderCommandBackHalf(com, out)
}

func renderCommandUsageBlock(com Command, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(subLinePadding[0]); err != nil {
		return err
	}
	if _, err := out.WriteString(com.Name()); err != nil {
		return err
	}
	return renderCommandUsageBackHalf(com, out)
}

func renderCommandLeaf(leaf CommandLeaf, out *bufio.Writer) error {
	if err := renderCommandLeafUsage(leaf, out); err != nil {
		return err
	}
	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	if leaf.HasAliases() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if _, err := out.WriteString(subLinePadding[0]); err != nil {
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
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	return renderCommandBackHalf(leaf, out)
}

func renderCommandBackHalf(com Command, out *bufio.Writer) error {

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
		if err := renderFlagGroups(com.FlagGroups(), 0, out); err != nil {
			return err
		}
	}

	if writeArgs {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if _, err := out.WriteString(headerPadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString(comArgs); err != nil {
			return err
		}

		for _, arg := range com.Arguments() {
			if _, err := out.WriteString(paragraphBreak); err != nil {
				return err
			}
			if err := renderArgument(arg, 1, out); err != nil {
				return err
			}
		}
	}

	return out.WriteByte(charLF)
}

func renderCommandLeafUsage(leaf CommandLeaf, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if err := renderSubCommandPath(leaf, out); err != nil {
		return err
	}
	return renderCommandUsageBackHalf(leaf, out)
}

func renderCommandUsageBackHalf(com Command, out *bufio.Writer) error {
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
					if err := renderShortestFlagLine(flag, out); err != nil {
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
			if err := renderArgumentName(arg, out); err != nil {
				return err
			}
		}
	}

	if com.HasUnmappedLabel() {
		if err := out.WriteByte(charSpace); err != nil {
			return err
		}
		if _, err := out.WriteString(com.GetUnmappedLabel()); err != nil {
			return err
		}
	}

	return nil
}

func renderCommandBranch(branch CommandBranch, out *bufio.Writer) error {
	if err := renderCommandBranchUsage(branch, out); err != nil {
		return err
	}
	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	hd := branch.HasDescription()
	hf := branch.HasFlagGroups()

	if branch.HasAliases() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if _, err := out.WriteString(subLinePadding[0]); err != nil {
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
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	if hd {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(branch.Description(), descriptionPadding[0], helpTextMaxWidth, out); err != nil {
			return err
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	// render flags
	if hf {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := renderFlagGroups(branch.FlagGroups(), 0, out); err != nil {
			return err
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	if err := renderCommandGroups(branch.CommandGroups(), 0, out); err != nil {
		return err
	}

	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	return nil
}

func renderCommandBranchUsage(node CommandBranch, out *bufio.Writer) error {
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
					if err := out.WriteByte(charSpace); err != nil {
						return err
					}
					if err := renderShortestFlagLine(flag, out); err != nil {
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

	return nil
}

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

const (
	defaultComGroupName = "Commands"
)

func renderCommandGroups(groups []CommandGroup, padding uint8, sb *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := sb.WriteString(paragraphBreak); err != nil {
				return err
			}
		}

		if err := renderCommandGroup(group, padding, sb); err != nil {
			return err
		}
	}

	return nil
}

func renderCommandGroup(group CommandGroup, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(headerPadding[padding]); err != nil {
		return err
	}

	if group.Name() == defaultGroupName {
		if _, err := sb.WriteString(defaultComGroupName); err != nil {
			return err
		}
	} else {
		if _, err := sb.WriteString(group.Name()); err != nil {
			return err
		}
	}

	if err := sb.WriteByte(charLF); err != nil {
		return err
	}

	if group.HasDescription() {
		if err := breakFmt(group.Description(), descriptionPadding[padding], helpTextMaxWidth, sb); err != nil {
			return err
		}
		if err := sb.WriteByte(charLF); err != nil {
			return err
		}
	}

	ordered := make([]string, 0, len(group.Leaves())+len(group.Branches()))
	lookup := make(map[string]CommandNode, len(group.Leaves())+len(group.Branches()))

	for _, node := range group.Leaves() {
		ordered = append(ordered, node.Name())
		lookup[node.Name()] = node
	}
	for _, node := range group.Branches() {
		ordered = append(ordered, node.Name())
		lookup[node.Name()] = node
	}

	slices.Sort(ordered)

	for i, name := range ordered {
		if i > 0 {
			if _, err := sb.WriteString(paragraphBreak); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(headerPadding[padding+1]); err != nil {
			return err
		}
		if _, err := sb.WriteString(name); err != nil {
			return err
		}

		node := lookup[name]

		if node.HasDescription() {
			if err := sb.WriteByte(charLF); err != nil {
				return err
			}
			if err := breakFmt(node.Description(), descriptionPadding[padding+1], helpTextMaxWidth, sb); err != nil {
				return err
			}
		}
	}

	return nil
}

const (
	subcommandPlaceholder = " <command>"
)

// renderCommandTree renders a help page for the given argo.CommandTree
// instance.
//
//     Usage:
//       svcctl [options] <command>
//
//       Description of my app.
//
//     General Flags
//       -f <string> | --file=<string>
//          Path to the target file.
//       -r | --reverse
//          Whether the output should be reversed.
//       -v | --verbose
//          Enable verbose logging.
//
//     Meta Commands
//       -h | --help
//          Prints this help text.
//
//     Commands
//       services (svc, s)
//         Aliases: svc
//            Subcommands for operating on one or more services.
//       storage
//         Aliases: store
//
//         Storage management operations.
func renderCommandTree(tree CommandTree, out *bufio.Writer) error {
	if err := renderCommandTreeUsageBlock(tree, out); err != nil {
		return err
	}
	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	hd := tree.HasDescription()
	hf := tree.HasFlagGroups()

	if hd {
		if err := breakFmt(tree.Description(), descriptionPadding[0], helpTextMaxWidth, out); err != nil {
			return err
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	if hf {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := renderFlagGroups(tree.FlagGroups(), 0, out); err != nil {
			return err
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	if err := out.WriteByte(charLF); err != nil {
		return err
	}
	if err := renderCommandGroups(tree.CommandGroups(), 0, out); err != nil {
		return err
	}

	if err := out.WriteByte(charLF); err != nil {
		return err
	}

	return nil
}

func renderCommandTreeUsageBlock(tree CommandTree, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(subLinePadding[0]); err != nil {
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
					if err := out.WriteByte(charSpace); err != nil {
						return err
					}
					if err := renderShortestFlagLine(flag, out); err != nil {
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
