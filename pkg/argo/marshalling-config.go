package argo

type UnmarshallingConfig struct {
	Integers IntegerUnmarshallingConfig

	Maps MapUnmarshallingConfig

	Slices SliceUnmarshallingConfig

	Time TimeUnmarshallingConfig
}

type IntegerUnmarshallingConfig struct {
	// OctalLeaders defines the prefixes used to signify that a value should be
	// parsed as octal.
	//
	// An empty slice will disable octal value parsing.
	//
	// Default: ["0o", "0O", "o", "O", "0"]
	//
	// Example octal values using the default prefixes:
	//
	//     o666    O666
	//     o0666   O0666
	//     0o666   0O666
	//     0o0666  0O0666
	//     0666
	OctalLeaders []string

	// HexLeaders defines the prefixes used to signify that a value should be
	// parsed as hexadecimal.
	//
	// An empty slice will disable hex value parsing.
	//
	// Default: ["0x", "0X", "x", "X"]
	//
	// Example hex values using the default prefixes:
	//
	//     xFA9    XFA9
	//     xfa9    Xfa9
	//     0xFA9   0XFA9
	//     0xfa9   0Xfa9
	HexLeaders []string

	// The integer base to use when no prefix is present.
	//
	// Default: base 10
	DefaultBase int
}

type MapUnmarshallingConfig struct {
	// KeyValSeparatorChars defines characters used to separate a key from a value
	// in an individual mapping entry.
	//
	// This character can be escaped with a '\' (backslash) character.
	//
	// The first unescaped instance of one of the defined characters in the
	// individual entry will be used as the divider, and any subsequent
	// appearances in the entry will be included in the value.
	//
	// Default: "=:" (equals, colon)
	//
	// Example key/value pairs using the default divider characters.  The second
	// column is a JSON representation of the parsed map
	//
	//     key:value            {"key": "value"}
	//     key=value            {"key": "value"}
	//     key\\:foo:value      {"key:foo": "value"}
	//     key\\=bar=value      {"key=bar": "value"}
	//     key\\:baz=value:a=b  {"key:baz": "value:a=b"}
	//     key:value=c:d        {"key": "value=c:d"}
	KeyValSeparatorChars string

	// Default: ",;" (comma, semicolon)
	EntrySeparatorChars string
}

type SliceUnmarshallingConfig struct {
	// Scanner is a function that provides a text scanner that will be used to
	// break an argument string into segments representing individual
	// unmarshalable values.
	//
	// The values split by the scanner will then be parsed into the type expected
	// by the argument binding slice.
	//
	// The default scanner breaks strings on comma characters.
	//
	// For example, given the input string "foo,bar,fizz,buzz", use of the default
	// scanner would result in a slice containing the values:
	//     [ foo, bar, fizz, buzz ]
	Scanner StringScannerFactory

	// The ByteSliceParser property controls the parser that will be used to
	// deserialize a raw argument string into a byte slice.
	//
	// This function will be called with the raw CLI input and will be expected to
	// return either the parsed byte slice or an error.  If the function returns
	// an error it will be passed up like any other parsing error.
	//
	// The default ByteSliceParser implementation takes the raw string and
	// directly casts it to a byte slice.
	//
	// WARNING: In v2 the default behavior will be changed to expect base64 input.
	ByteSliceParser ByteSliceParser
}

type TimeUnmarshallingConfig struct {
	// DateFormats configures the date-time formats that the unmarshaler will use
	// when attempting to parse a date value.
	//
	// By default, the RFC3339 and RFC3339 nano patterns are used.
	DateFormats []string
}
