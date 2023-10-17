package argo

import (
	"strings"
)

func renderFlag(flag Flag, padding uint8, sb *strings.Builder) {
	sb.WriteString(headerPadding[padding])

	if flag.HasLongForm() {
		if flag.HasShortForm() {
			sb.WriteByte(charDash)
			sb.WriteByte(flag.ShortForm())

			if flag.HasArgument() {
				sb.WriteByte(charSpace)
				renderArgumentName(flag.Argument(), sb)
			}

			sb.WriteString(flagDivider)
		}

		sb.WriteString(strDoubleDash)
		sb.WriteString(flag.LongForm())

		if flag.HasArgument() {
			sb.WriteByte(charEquals)
			renderArgumentName(flag.Argument(), sb)
		}
	} else {
		sb.WriteByte(charDash)
		sb.WriteByte(flag.ShortForm())

		if flag.HasArgument() {
			sb.WriteByte(charSpace)
			renderArgumentName(flag.Argument(), sb)
		}
	}

	if flag.HasDescription() {
		sb.WriteByte(charLF)
		breakFmt(flag.Description(), descriptionPadding[padding], helpTextMaxWidth, sb)
		sb.WriteByte(charLF)
	}

	if flag.HasArgument() {
		renderArgument(flag.Argument(), padding+1, sb)
	}
}

func renderShortestFlagLine(flag Flag, sb *strings.Builder) {
	if flag.HasShortForm() {
		sb.WriteByte(charDash)
		sb.WriteByte(flag.ShortForm())

		if flag.HasArgument() {
			sb.WriteByte(charSpace)
			renderArgumentName(flag.Argument(), sb)
		}
	} else {
		sb.WriteString(strDoubleDash)
		sb.WriteString(flag.LongForm())

		if flag.HasArgument() {
			sb.WriteByte(charEquals)
			renderArgumentName(flag.Argument(), sb)
		}
	}
}

const (
	fgDefaultName = "General Flags"
	fgSingleName  = "Flags"
)

func renderFlagGroups(groups []FlagGroup, padding uint8, out *strings.Builder) {
	for i, group := range groups {
		if i > 0 {
			out.WriteByte(charLF)
		}

		renderFlagGroup(group, padding, out, len(groups) > 1)
	}
}

func renderFlagGroup(
	group FlagGroup,
	padding uint8,
	out *strings.Builder,
	multiple bool,
) {
	out.WriteString(headerPadding[padding])

	if group.Name() == defaultGroupName {
		if multiple {
			out.WriteString(fgDefaultName)
		} else {
			out.WriteString(fgSingleName)
		}
	} else {
		out.WriteString(group.Name())
	}

	out.WriteByte(charLF)

	// If the group has a description, print it out.
	if group.HasDescription() {
		breakFmt(group.Description(), subLinePadding[padding], helpTextMaxWidth, out)
		out.WriteByte(charLF)
	}

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			out.WriteByte(charLF)
		}

		renderFlag(flag, padding+1, out)
	}
}
