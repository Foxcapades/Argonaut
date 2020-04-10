package render

import (
	"strings"
)

const (
	comPrefix = "Usage:\n  "
	comOpts   = " [OPTIONS]"
)

func Command(com ac) string {
	var out strings.Builder
	command(com, &out)
	return out.String()
}

func command(com ac, out *strings.Builder) {
	out.WriteString(comPrefix)
	out.WriteString(com.Name())
	appendComFlags(com, out)
	appendComArgs(com, out)
	appendComDesc(com, out)

	appendComFlagGroups(com, out)
}

func appendComFlags(com ac, out *strings.Builder) {
	if len(com.FlagGroups()) == 0 {
		return
	}

	out.WriteString(comOpts)
}

func appendComDesc(com ac, out *strings.Builder) {
	if com.HasDescription() {
		out.WriteString(dblLineBreak)
		out.WriteString(com.Description())
	}
}

func appendComArgs(com ac, out *strings.Builder) {
	if len(com.Arguments()) == 0 {
		return
	}

	for _, arg := range com.Arguments() {
		if arg.Required() {
			out.WriteByte(' ')
			FormattedArgName(arg, out)
		} else {
			FormattedArgName(arg, out)
		}
	}
}

func appendComFlagGroups(com ac, out *strings.Builder) {
	for _, fg := range com.FlagGroups() {
		out.WriteString(dblLineBreak)
		flagGroup(fg, out)
	}
}
