package xref

import "unsafe"

func HardCast[I any, O any](value I) O {
	return *((*O)(unsafe.Pointer(&value)))
}
