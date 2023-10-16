package render

import (
	"reflect"
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/impl/argument/argutil"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const (
	argReqPrefix = '<'
	argReqSuffix = '>'
	argOptPrefix = '['
	argOptSuffix = ']'
)

func ArgumentName(a argo.Argument, out *strings.Builder) {
	if a.HasBinding() && a.BindingType().Kind() == reflect.Bool {
		return
	}

	if a.IsRequired() {
		out.WriteByte(argReqPrefix)
		out.WriteString(argutil.ArgName(a))
		out.WriteByte(argReqSuffix)
	} else {
		out.WriteByte(argOptPrefix)
		out.WriteString(argutil.ArgName(a))
		out.WriteByte(argOptSuffix)
	}
}

func Argument(arg argo.Argument, padding uint8, out *strings.Builder) {
	out.WriteString(headerPadding[padding])
	ArgumentName(arg, out)

	if arg.HasDescription() {
		out.WriteByte(chars.CharLF)
		BreakFmt(arg.Description(), descriptionPadding[padding], maxWidth, out)
	}
}
