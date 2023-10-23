package unmarshal

import "reflect"

var numericKinds = map[reflect.Kind]bool{
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
}

func reflectIsUnmarshaler(t reflect.Type, ut reflect.Type) bool {
	return t.AssignableTo(ut)
}

func reflectIsInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

func reflectIsBasicKind(k reflect.Kind) bool {
	return k == reflect.String || k == reflect.Bool || reflectIsNumericKind(k)
}

func reflectIsNumericKind(k reflect.Kind) bool {
	return numericKinds[k]
}

func reflectIsByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		t.Elem().Kind() == reflect.Uint8
}

func reflectIsBasicSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		reflectIsBasicKind(t.Elem().Kind())
}
