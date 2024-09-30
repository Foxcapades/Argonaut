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
//     func (t *MyCustomType) UnmarshalCLI(raw string) error {
//         t.Value = raw
//         return nil
//     }
//
//     var value MyCustomType
//
//     cli.Argument().
//         WithBinding(&value)
type Unmarshaler interface {

	// UnmarshalCLI parses the given raw input string into a value of the expected
	// type.
	UnmarshalCLI(raw string) error
}

// UnmarshalerFunc defines a function that implements the Unmarshaler interface.
type UnmarshalerFunc func(raw string) error

func (u UnmarshalerFunc) UnmarshalCLI(raw string) error {
	return u(raw)
}
