package argo

import (
	"strconv"
	"unsafe"
)

// Consumer defines a type that may be used as an Argument binding value to
// unmarshal and consume the given raw string.
//
// Example:
//     cli.Argument()
//         WithBinding(ConsumerFunc(func(raw string) error {
//             // Do something with the raw string
//             return nil
//         })
type Consumer interface {

	// Accept is handed the raw string for its obviously nefarious purposes and
	// may return an error.
	Accept(raw string) error
}

// ConsumerFunc defines a function that implements the Consumer interface.
type ConsumerFunc func(raw string) error

func (u ConsumerFunc) Accept(raw string) error {
	return u(raw)
}

// ValueUnmarshaler is the internal unmarshaler type that is used to parse
// input string values into the given pointer (val).
//
// For arguments with a specific type, this may be implemented and passed to the
// argument builder to use instead of the default "magic" unmarshaler instance.
//
// An example ValueUnmarshaler that is used for a single Argument may safely
// make assumptions about the type of the provided pointer, however reflection
// or unsafe pointers will still be required to fill the given pointer value.
//
// If you were reasonable, you could do something like this:
//    ValueUnmarshalerFunc(func(raw string, val any) error {
//        if parsed, err := strconv.Atoi(raw); err != nil {
//            return err
//        } else {
//            rVal := reflect.ValueOf(val)
//            if rVal.Kind() != reflect.Pointer {
//                return errors.New("value must be a pointer")
//            } else {
//                rVal.Elem().Set(reflect.ValueOf(parsed))
//                return nil
//            }
//        }
//    })
//
// And if you are half-baked, you could try this:
//     ValueUnmarshalerFunc(func(raw string, val any) error {
//         if parsed, err := strconv.Atoi(raw); err != nil {
//             return err
//         } else {
//             **(**int)(unsafe.Add(unsafe.Pointer(&val), bits.UintSize/8)) = parsed
//             return nil
//         }
//     })
type ValueUnmarshaler interface {

	// Unmarshal accepts a raw value and a pointer to an arbitrary type and is
	// expected to fill the pointer value by parsing the raw string, returning an
	// error if parsing fails.
	Unmarshal(raw string, val any) error
}

// ValueUnmarshalerFunc defines a function that implements the ValueUnmarshaler
// interface.
type ValueUnmarshalerFunc func(raw string, val any) error

func (v ValueUnmarshalerFunc) Unmarshal(raw string, val any) error {
	return v(raw, val)
}

var foo = ValueUnmarshalerFunc(func(raw string, val any) error {
	ptr := *(**int)(unsafe.Pointer(&val))
	if parsed, err := strconv.Atoi(raw); err != nil {
		return err
	} else {
		*ptr = parsed
		return nil
	}
})
