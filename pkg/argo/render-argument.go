package argo

import (
	"reflect"
	"strings"
	"time"
)

func renderArgName(a Argument) string {
	if a.HasName() {
		return a.Name()
	} else {
		return "arg"
	}
}

func renderFormatArgType(kind reflect.Type) string {
	switch kind.Kind() {
	case reflect.Map:
		return renderFormatArgType(kind.Key()) + "=" + renderFormatArgType(kind.Elem())
	case reflect.Slice:
		if kind.Elem().Kind() == reflect.Uint8 {
			return "bytes"
		}
		return renderFormatArgType(kind.Elem())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "int"
	case reflect.Int64:
		if kind.AssignableTo(reflectGetRootValue(reflect.ValueOf(time.Duration(0))).Type()) {
			return "duration"
		} else {
			return "int"
		}
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:

	}
	return kind.String()
}

const (
	argReqPrefix = '<'
	argReqSuffix = '>'
	argOptPrefix = '['
	argOptSuffix = ']'
)

func renderArgumentName(a Argument, out *strings.Builder) {
	if a.hasBinding() && a.bindingType().Kind() == reflect.Bool {
		return
	}

	if a.IsRequired() {
		out.WriteByte(argReqPrefix)
		out.WriteString(renderArgName(a))
		out.WriteByte(argReqSuffix)
	} else {
		out.WriteByte(argOptPrefix)
		out.WriteString(renderArgName(a))
		out.WriteByte(argOptSuffix)
	}
}

func renderFlagArgument(arg Argument, padding uint8, out *strings.Builder) {
	out.WriteString(subLinePadding[padding])
	renderArgumentName(arg, out)

	if arg.HasDescription() {
		out.WriteByte(charLF)
		breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out)
	}
}

func renderArgument(arg Argument, padding uint8, out *strings.Builder) {
	out.WriteString(headerPadding[padding])
	renderArgumentName(arg, out)

	if arg.HasDescription() {
		out.WriteByte(charLF)
		breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out)
	}
}
