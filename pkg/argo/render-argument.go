package argo

import (
	"bufio"
	"reflect"
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

func renderArgumentName(a Argument, out *bufio.Writer) error {
	if a.hasBinding() && a.bindingType().Kind() == reflect.Bool {
		return nil
	}

	if a.IsRequired() {
		if err := out.WriteByte(argReqPrefix); err != nil {
			return err
		}
		if _, err := out.WriteString(renderArgName(a)); err != nil {
			return err
		}
		if err := out.WriteByte(argReqSuffix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(argOptPrefix); err != nil {
			return err
		}
		if _, err := out.WriteString(renderArgName(a)); err != nil {
			return err
		}
		if err := out.WriteByte(argOptSuffix); err != nil {
			return err
		}
	}

	return nil
}

func renderFlagArgument(arg Argument, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(subLinePadding[padding]); err != nil {
		return err
	}
	if err := renderArgumentName(arg, out); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out); err != nil {
			return err
		}
	}
	return nil
}

func renderArgument(arg Argument, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(headerPadding[padding]); err != nil {
		return err
	}
	if err := renderArgumentName(arg, out); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out); err != nil {
			return err
		}
	}

	return nil
}
