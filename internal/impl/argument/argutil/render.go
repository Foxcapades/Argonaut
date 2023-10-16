package argutil

import (
	"reflect"
	"time"

	"github.com/Foxcapades/Argonaut/internal/xref"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func ArgName(a argo.Argument) string {
	if a.HasName() {
		return a.Name()
	} else {
		return "arg"
	}
}

func FormatArgType(kind reflect.Type) string {
	switch kind.Kind() {
	case reflect.Map:
		return FormatArgType(kind.Key()) + "=" + FormatArgType(kind.Elem())
	case reflect.Slice:
		if kind.Elem().Kind() == reflect.Uint8 {
			return "bytes"
		}
		return FormatArgType(kind.Elem())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "int"
	case reflect.Int64:
		if kind.AssignableTo(xref.GetRootValue(reflect.ValueOf(time.Duration(0))).Type()) {
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
