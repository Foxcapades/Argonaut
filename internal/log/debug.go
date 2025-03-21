//go:build dev

package log

import (
	"fmt"
	"os"
	"reflect"
)

func DebugLn1[T any](format string, arg T) {
	var a any

	if reflect.TypeOf(arg).Kind() == reflect.Func {
		a = reflect.ValueOf(arg).Call(nil)
	} else {
		a = arg
	}

	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a)
}

func DebugLn2[A, B any](format string, arg1 A, arg2 B) {
	var a any
	var b any

	if reflect.TypeOf(arg1).Kind() == reflect.Func {
		a = reflect.ValueOf(arg1).Call(nil)
	} else {
		a = arg1
	}
	if reflect.TypeOf(arg2).Kind() == reflect.Func {
		b = reflect.ValueOf(arg2).Call(nil)
	} else {
		b = arg2
	}

	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a, b)
}
