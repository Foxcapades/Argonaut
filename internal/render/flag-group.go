package render

import (
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const (
	fgDefaultName = "General Flags"
)

func FlagGroups(groups []argo.FlagGroup, padding uint8, out *strings.Builder, forceHeaders bool) {
	for i, group := range groups {
		if i > 0 {
			out.WriteByte(chars.CharLF)
		}

		FlagGroup(group, padding, out, forceHeaders || len(groups) > 1)
	}
}

func FlagGroup(
	group argo.FlagGroup,
	padding uint8,
	out *strings.Builder,
	multiple bool,
) {
	out.WriteString(headerPadding[padding])
	if group.Name() == consts.DefaultGroupName {
		if multiple || group.HasDescription() {
			out.WriteString(fgDefaultName)
		}
	} else {
		out.WriteString(group.Name())
	}

	// If the group has a description, print it out.
	if group.HasDescription() {
		out.WriteByte(chars.CharLF)
		BreakFmt(group.Description(), subLinePadding[padding], maxWidth, out)
		out.WriteByte(chars.CharLF)
	}

	out.WriteByte(chars.CharLF)

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			out.WriteString(paragraphBreak)
		}

		Flag(flag, padding+1, out)
	}
}
