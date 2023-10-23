package xunsafe

import (
	"unsafe"
)

func HardCast[T any](value any) T {
	if val, ok := value.(T); ok {
		return val
	}

	return *((*T)(unsafe.Pointer(&value)))
}
