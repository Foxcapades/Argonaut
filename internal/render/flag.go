package render

import (
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func Flag(flag argo.Flag, padding uint8, sb *strings.Builder) {
	sb.WriteString(headerPadding[padding])

	if flag.HasLongForm() {
		if flag.HasShortForm() {
			sb.WriteByte(chars.CharDash)
			sb.WriteByte(flag.ShortForm())

			if flag.HasArgument() {
				sb.WriteByte(chars.CharSpace)
				ArgumentName(flag.Argument(), sb)
			}

			sb.WriteString(flagDivider)
		}

		sb.WriteString(chars.StrDoubleDash)
		sb.WriteString(flag.LongForm())

		if flag.HasArgument() {
			sb.WriteByte(chars.CharEquals)
			ArgumentName(flag.Argument(), sb)
		}
	} else {
		sb.WriteByte(chars.CharDash)
		sb.WriteByte(flag.ShortForm())

		if flag.HasArgument() {
			sb.WriteByte(chars.CharSpace)
			ArgumentName(flag.Argument(), sb)
		}
	}

	if flag.HasDescription() {
		sb.WriteByte(chars.CharLF)
		BreakFmt(flag.Description(), descriptionPadding[padding], maxWidth, sb)
	}

	if flag.HasArgument() {
		Argument(flag.Argument(), padding+1, sb)
	}
}

func ShortestFlagLine(flag argo.Flag, sb *strings.Builder) {
	if flag.HasShortForm() {
		sb.WriteByte(chars.CharDash)
		sb.WriteByte(flag.ShortForm())

		if flag.HasArgument() {
			sb.WriteByte(chars.CharSpace)
			ArgumentName(flag.Argument(), sb)
		}
	} else {
		sb.WriteString(chars.StrDoubleDash)
		sb.WriteString(flag.LongForm())

		if flag.HasArgument() {
			sb.WriteByte(chars.CharEquals)
			ArgumentName(flag.Argument(), sb)
		}
	}
}
