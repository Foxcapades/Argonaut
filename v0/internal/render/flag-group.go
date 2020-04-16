package render

import (
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

const (
	fgDefaultName = "Flag Group"
	fgShortPrefix = '-'
	fgLongPrefix  = "--"
	fgShortAssign = sngSpace
	fgLongAssign  = '='
	fgDivider     = " | "
	fgEmptyDiv    = "   "
	fgPadding     = "  "
)

func flagGroup(fg A.FlagGroup, out *strings.Builder) {
	if !fg.HasFlags() {
		return
	}

	if fg.HasName() {
		out.WriteString(fg.Name())
	} else {
		out.WriteString(fgDefaultName)
	}

	out.WriteByte(sngLineBreak)

	result := flagNames(fg.Flags())
	pad := result.sLen + result.lLen + len(fgDivider) + (len(fgPadding) * 2)

	for i, flag := range fg.Flags() {
		hasShort := result.shorts[i] != ""
		hasLong := result.longs[i] != ""

		out.WriteString(fgPadding)

		WritePadded(result.shorts[i], result.sLen, out)

		if hasLong && hasShort {
			out.WriteString(fgDivider)
		} else {
			out.WriteString(fgEmptyDiv)
		}

		WritePadded(result.longs[i], result.lLen, out)

		if flag.HasDescription() {
			out.WriteString(fgPadding)
			BreakFmt(flag.Description(), pad, maxWidth, out)
		}
		out.WriteByte(sngLineBreak)
	}
}

type flagResult struct {
	shorts []string
	longs  []string
	sLen   int
	lLen   int
}

func flagNames(flags []af) (out flagResult) {
	ln := len(flags)
	bld := strings.Builder{}

	out.shorts = make([]string, ln)
	out.longs = make([]string, ln)

	bld.Grow(10)
	bld.Reset()

	for i, flag := range flags {
		var arg string

		if flag.HasArgument() {
			FormattedArgName(flag.Argument(), &bld)
			arg = bld.String()
			bld.Reset()
		}

		if flag.HasShort() {
			bld.WriteByte(fgShortPrefix)
			bld.WriteByte(flag.Short())

			if len(arg) > 0 {
				bld.WriteByte(fgShortAssign)
				bld.WriteString(arg)
			}

			out.shorts[i] = bld.String()
			bld.Reset()

			if len(out.shorts[i]) > out.sLen {
				out.sLen = len(out.shorts[i])
			}
		}

		if flag.HasLong() {
			bld.WriteString(fgLongPrefix)
			bld.WriteString(flag.Long())

			if len(arg) > 0 {
				bld.WriteByte(fgLongAssign)
				bld.WriteString(arg)
			}

			out.longs[i] = bld.String()
			bld.Reset()

			if len(out.longs[i]) > out.lLen {
				out.lLen = len(out.longs[i])
			}
		}
	}

	return
}
