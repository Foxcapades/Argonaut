package argo

// Unmarshaler defines a type that may be used as an Argument binding value to
// unmarshal the given raw string into a custom type.
//
// Example 1:
//     cli.Argument()
//         WithBinding(UnmarshalerFunc(func(raw string) error {
//             // Do something with the raw string
//             return nil
//         })
//
// Example 2:
//     type MyCustomType struct {
//         Value string
//     }
//     func (t *MyCustomType) Unmarshal(raw string) error {
//         t.Value = raw
//         return nil
//     }
//
//     var value MyCustomType
//
//     cli.Argument().
//         WithBinding(&value)
type Unmarshaler interface {

	// Unmarshal is handed the raw string for its obviously nefarious purposes and
	// may return an error.
	Unmarshal(raw string) error
}

// UnmarshalerFunc defines a function that implements the Unmarshaler interface.
type UnmarshalerFunc func(raw string) error

func (u UnmarshalerFunc) Unmarshal(raw string) error {
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
