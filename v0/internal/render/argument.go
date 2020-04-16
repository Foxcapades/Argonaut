package render

import (
	R "reflect"
	"strings"
)

const (
	argReqPrefix = '<'
	argReqSuffix = '>'
	argOptPrefix = '['
	argOptSuffix = ']'
)

func FormattedArgName(a aa, out *strings.Builder) {
	if a.HasBinding() && a.BindingType().Kind() == R.Bool {
		return
	}

	if a.Required() {
		out.WriteByte(argReqPrefix)
		out.WriteString(ArgName(a))
		out.WriteByte(argReqSuffix)
	} else {
		out.WriteByte(argOptPrefix)
		out.WriteString(ArgName(a))
		out.WriteByte(argOptSuffix)
	}
}

func ArgName(a aa) string {
	if a.HasName() {
		return a.Name()
	} else {
		return FormatArgType(a.BindingType())
	}
}

func FormatArgType(kind R.Type) string {
	switch kind.Kind() {
	case R.Map:
		return FormatArgType(kind.Key()) + "=" + FormatArgType(kind.Elem())
	case R.Slice:
		if kind.Elem().Kind() == R.Uint8 {
			return "bytes"
		}
		return FormatArgType(kind.Elem())
	case R.Uint, R.Uint8, R.Uint16, R.Uint32, R.Uint64:
		return "uint"
	case R.Int, R.Int8, R.Int16, R.Int32, R.Int64:
		return "int"
	case R.Float32, R.Float64:
		return "float"
	case R.Bool:

	}
	return kind.String()
}
