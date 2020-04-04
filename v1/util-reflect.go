package argo

import "reflect"

func isPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}

func isNumericType(obj interface{}) bool {
	tp := reflect.TypeOf(obj).Elem().Kind()

	return tp == reflect.Int ||
		tp == reflect.Float32 ||
		tp == reflect.Int64 ||
		tp == reflect.Float64 ||
		tp == reflect.Uint64 ||
		tp == reflect.Uint ||
		tp == reflect.Int32 ||
		tp == reflect.Uint32 ||
		tp == reflect.Uint8 ||
		tp == reflect.Int8 ||
		tp == reflect.Int16 ||
		tp == reflect.Uint16
}
