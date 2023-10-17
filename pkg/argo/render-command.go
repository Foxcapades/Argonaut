package argo

import (
	"slices"
	"strings"
)

const (
	comPrefix = "Usage:\n  "
	comOpts   = " [options]"
	comArgs   = "Positional Arguments"
)

func renderCommand(com Command) string {
	out := strings.Builder{}
	renderCommandUsageBlock(com, &out)
	renderCommandBackHalf(com, &out)
	return out.String()
}

func renderCommandUsageBlock(com Command, out *strings.Builder) {
	out.WriteString(comPrefix)
	out.WriteString(com.Name())
	renderCommandUsageBackHalf(com, out)
}

func renderCommandLeaf(leaf CommandLeaf) string {
	out := strings.Builder{}
	renderCommandLeafUsage(leaf, &out)
	renderCommandBackHalf(leaf, &out)
	return out.String()
}

func renderCommandBackHalf(com Command, out *strings.Builder) {
	// If the command has a description, append it.
	if com.HasDescription() {
		out.WriteByte(charLF)
		breakFmt(com.Description(), subLinePadding[0], helpTextMaxWidth, out)
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
		out.WriteString(paragraphBreak)
		renderFlagGroups(com.FlagGroups(), 0, out)
	}

	if writeArgs {
		out.WriteString(paragraphBreak)
		out.WriteString(headerPadding[0])
		out.WriteString(comArgs)

		for _, arg := range com.Arguments() {
			out.WriteString(paragraphBreak)
			renderArgument(arg, 1, out)
		}
	}

	out.WriteByte(charLF)
}

func renderCommandLeafUsage(leaf CommandLeaf, out *strings.Builder) {
	out.WriteString(comPrefix)
	renderSubCommandPath(leaf, out)
	renderCommandUsageBackHalf(leaf, out)
}

func renderCommandUsageBackHalf(com Command, out *strings.Builder) {
	// If the command has flag groups
	if com.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range com.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(charSpace)
					renderShortestFlagLine(flag, out)
				} else {
					hasOptionalFlags = true
				}
			}
		}

		// If there are any optional flags append the general "[OPTIONS]" text.
		if hasOptionalFlags {
			out.WriteString(comOpts)
		}
	}

	// After all the flag groups have been rendered, append the argument names.
	if com.HasArguments() {
		for _, arg := range com.Arguments() {
			out.WriteByte(charSpace)
			renderArgumentName(arg, out)
		}
	}

	if com.HasUnmappedLabel() {
		out.WriteByte(charSpace)
		out.WriteString(com.GetUnmappedLabel())
	}

}

func renderCommandBranch(branch CommandBranch) string {
	out := strings.Builder{}

	renderCommandBranchUsage(branch, &out)

	hd := branch.HasDescription()
	hf := branch.HasFlagGroups()

	if branch.HasAliases() {
		out.WriteString(" (")

		aliases := branch.Aliases()
		slices.Sort(aliases)

		for i, alias := range aliases {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(alias)
		}

		out.WriteByte(')')
	}

	if hd {
		out.WriteByte(charLF)
		breakFmt(branch.Description(), descriptionPadding[0], helpTextMaxWidth, &out)
	}

	// render flags
	if hf {
		out.WriteString(paragraphBreak)
		renderFlagGroups(branch.FlagGroups(), 0, &out)
	}

	out.WriteString(paragraphBreak)
	renderCommandGroups(branch.CommandGroups(), 0, &out, hf || hd)

	out.WriteByte(charLF)

	return out.String()
}

func renderCommandBranchUsage(node CommandBranch, out *strings.Builder) {
	out.WriteString(comPrefix)
	renderSubCommandPath(node, out)

	if node.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range node.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(charSpace)
					renderShortestFlagLine(flag, out)
				} else {
					hasOptionalFlags = true
				}
			}
		}

		// If there are any optional flags append the general "[OPTIONS]" text.
		if hasOptionalFlags {
			out.WriteString(comOpts)
		}
	}
}

func renderSubCommandPath(node CommandNode, out *strings.Builder) {
	path := make([]string, 0, 4)

	current := node
	for current != nil {
		path = append(path, current.Name())
		current = current.Parent()
	}

	slices.Reverse(path)

	out.WriteString(subLinePadding[0])
	for i, segment := range path {
		if i > 0 {
			out.WriteByte(charSpace)
		}
		out.WriteString(segment)
	}
}

const (
	defaultComGroupName = "Commands"
)

func renderCommandGroups(groups []CommandGroup, padding uint8, sb *strings.Builder, forceHeaders bool) {
	for i, group := range groups {
		if i > 0 {
			sb.WriteString(paragraphBreak)
		}

		renderCommandGroup(group, padding, sb, forceHeaders || len(groups) > 1)
	}
}

func renderCommandGroup(
	group CommandGroup,
	padding uint8,
	sb *strings.Builder,
	forceHeader bool,
) {
	sb.WriteString(headerPadding[padding])

	if group.Name() == defaultGroupName {
		sb.WriteString(defaultComGroupName)
	} else {
		sb.WriteString(group.Name())
	}

	if group.HasDescription() {
		sb.WriteByte(charLF)
		breakFmt(group.Description(), descriptionPadding[padding], helpTextMaxWidth, sb)
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

	for _, name := range ordered {
		sb.WriteByte(charLF)

		if group.HasDescription() {
			sb.WriteByte(charLF)
		}

		sb.WriteString(headerPadding[padding+1])
		sb.WriteString(name)

		node := lookup[name]

		if node.HasDescription() {
			sb.WriteByte(charLF)
			breakFmt(node.Description(), descriptionPadding[padding+1], helpTextMaxWidth, sb)
		}
	}
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
func renderCommandTree(tree CommandTree) string {
	out := strings.Builder{}

	renderCommandTreeUsageBlock(tree, &out)

	hd := tree.HasDescription()
	hf := tree.HasFlagGroups()

	if hd {
		out.WriteString(paragraphBreak)
		breakFmt(tree.Description(), subLinePadding[0], helpTextMaxWidth, &out)
	}

	if hf {
		out.WriteString(paragraphBreak)
		renderFlagGroups(tree.FlagGroups(), 0, &out)
	}

	out.WriteString(paragraphBreak)
	renderCommandGroups(tree.CommandGroups(), 0, &out, hd || hf)

	out.WriteByte(charLF)

	return out.String()
}

func renderCommandTreeUsageBlock(tree CommandTree, out *strings.Builder) {
	out.WriteString(comPrefix)
	out.WriteString(tree.Name())

	if tree.HasFlagGroups() {
		hasOptionalFlags := false

		for _, group := range tree.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(charSpace)
					renderShortestFlagLine(flag, out)
				} else {
					hasOptionalFlags = true
				}
			}
		}

		if hasOptionalFlags {
			out.WriteString(comOpts)
		}
	}

	out.WriteString(subcommandPlaceholder)
}
