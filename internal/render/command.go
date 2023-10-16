package render

import (
	"slices"
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const (
	comPrefix = "Usage:\n  "
	comOpts   = " [options]"
	comArgs   = "Positional Arguments"
)

func Command(com argo.Command) string {
	out := strings.Builder{}
	CommandUsageBlock(com, &out)
	CommandBackHalf(com, &out)
	return out.String()
}

func CommandUsageBlock(com argo.Command, out *strings.Builder) {
	out.WriteString(comPrefix)
	out.WriteString(com.Name())
	CommandUsageBackHalf(com, out)
}

func CommandLeaf(leaf argo.CommandLeaf) string {
	out := strings.Builder{}
	CommandLeafUsage(leaf, &out)
	CommandBackHalf(leaf, &out)
	return out.String()
}

func CommandBackHalf(com argo.Command, out *strings.Builder) {
	// If the command has a description, append it.
	if com.HasDescription() {
		out.WriteByte(chars.CharLF)
		BreakFmt(com.Description(), subLinePadding[0], maxWidth, out)
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
		FlagGroups(com.FlagGroups(), 1, out, com.HasDescription() || writeArgs)
	}

	if writeArgs {
		out.WriteString(paragraphBreak)
		out.WriteString(comArgs)

		for _, arg := range com.Arguments() {
			out.WriteString(paragraphBreak)
			Argument(arg, 1, out)
		}
	}
}

func CommandLeafUsage(leaf argo.CommandLeaf, out *strings.Builder) {
	out.WriteString(comPrefix)
	SubCommandPath(leaf, out)
	CommandUsageBackHalf(leaf, out)
}

func CommandUsageBackHalf(com argo.Command, out *strings.Builder) {
	// If the command has flag groups
	if com.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range com.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(chars.CharSpace)
					ShortestFlagLine(flag, out)
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
			out.WriteByte(chars.CharSpace)
			ArgumentName(arg, out)
		}
	}

	if com.HasUnmappedLabel() {
		out.WriteByte(chars.CharSpace)
		out.WriteString(com.GetUnmappedLabel())
	}

}

func CommandBranch(branch argo.CommandBranch, out *strings.Builder) {
	CommandBranchUsage(branch, out)

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
		out.WriteByte(chars.CharLF)
		BreakFmt(branch.Description(), descriptionPadding[0], maxWidth, out)
	}

	// render flags
	if hf {
		out.WriteString(paragraphBreak)
		FlagGroups(branch.FlagGroups(), 1, out, hf || hd)
	}

	out.WriteString(paragraphBreak)
	CommandGroups(branch.CommandGroups(), 1, out, hf || hd)
}

func CommandBranchUsage(node argo.CommandBranch, out *strings.Builder) {
	out.WriteString(comPrefix)
	SubCommandPath(node, out)

	if node.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range node.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(chars.CharSpace)
					ShortestFlagLine(flag, out)
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

func SubCommandPath(node argo.CommandNode, out *strings.Builder) {
	path := make([]string, 0, 4)

	current := node
	for current != nil {
		path = append(path, current.Name())
	}

	slices.Reverse(path)

	out.WriteString(subLinePadding[0])
	for i, segment := range path {
		if i > 0 {
			out.WriteByte(chars.CharSpace)
		}
		out.WriteString(segment)
	}
}

const (
	defaultComGroupName = "Default"
)

func CommandGroups(groups []argo.CommandGroup, padding uint8, sb *strings.Builder, forceHeaders bool) {
	for i, group := range groups {
		if i > 0 {
			sb.WriteString(paragraphBreak)
		}

		CommandGroup(group, padding, sb, forceHeaders || len(groups) > 1)
	}
}

func CommandGroup(
	group argo.CommandGroup,
	padding uint8,
	sb *strings.Builder,
	forceHeader bool,
) {
	if group.Name() == consts.DefaultGroupName {
		if forceHeader || group.HasDescription() {
			sb.WriteString(headerPadding[padding])
			sb.WriteString(defaultComGroupName)
		}
	} else {
		sb.WriteString(headerPadding[padding])
		sb.WriteString(group.Name())
	}

	if group.HasDescription() {
		sb.WriteByte(chars.CharLF)
		BreakFmt(group.Description(), descriptionPadding[padding], maxWidth, sb)
	}

	ordered := make([]string, 0, len(group.Leaves())+len(group.Branches()))
	lookup := make(map[string]argo.CommandNode, len(group.Leaves())+len(group.Branches()))

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
		sb.WriteByte(chars.CharLF)

		if group.HasDescription() {
			sb.WriteByte(chars.CharLF)
		}

		sb.WriteString(headerPadding[padding+1])
		sb.WriteString(name)

		node := lookup[name]

		if node.HasDescription() {
			sb.WriteByte(chars.CharLF)
			BreakFmt(node.Description(), descriptionPadding[padding+1], maxWidth, sb)
		}
	}
}

const (
	subcommandPlaceholder = " <command>"
)

// CommandTree renders a help page for the given argo.CommandTree instance.
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
func CommandTree(tree argo.CommandTree, out *strings.Builder) {
	CommandTreeUsageBlock(tree, out)

	hd := tree.HasDescription()
	hf := tree.HasFlagGroups()

	if hd {
		out.WriteString(paragraphBreak)
		BreakFmt(tree.Description(), subLinePadding[0], maxWidth, out)
	}

	if hf {
		out.WriteString(paragraphBreak)
		FlagGroups(tree.FlagGroups(), 0, out, hd)
	}

	out.WriteString(paragraphBreak)
	CommandGroups(tree.CommandGroups(), 0, out, hd || hf)
}

func CommandTreeUsageBlock(tree argo.CommandTree, out *strings.Builder) {
	out.WriteString(comPrefix)
	out.WriteString(tree.Name())

	if tree.HasFlagGroups() {
		hasOptionalFlags := false

		for _, group := range tree.FlagGroups() {
			for _, flag := range group.Flags() {
				if flag.IsRequired() {
					out.WriteByte(chars.CharSpace)
					ShortestFlagLine(flag, out)
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
