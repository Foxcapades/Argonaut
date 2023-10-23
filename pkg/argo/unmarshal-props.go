package argo

import "time"

// UnmarshalProps defines configuration options for the included "magic"
// Unmarshaler implementation.
//
// This configuration may be used to customize how the unmarshaler behaves or
// parses input values.
type UnmarshalProps struct {

	// Integers defines settings for parsing integral types
	// (int, int*, uint, uint*).
	Integers UnmarshalIntegerProps `json:"integers"`

	// Maps defines settings to use when parsing mappings from
	// the command line
	Maps UnmarshalMapProps `json:"maps"`

	Slices UnmarshalSliceProps `json:"slices"`

	Time UnmarshalTimeProps `json:"time"`
}

// UnmarshalIntegerProps defines int parsing specific options for the "magic"
// unmarshaler.
type UnmarshalIntegerProps struct {

	// OctalLeaders defines the prefixes used to signify
	// that a value should be parsed as octal.
	//
	// An empty slice will disable octal value parsing.
	//
	// Default: ["0o", "0O", "o", "O"]
	//
	// Example octal values using the default prefixes:
	//
	//     o666    O666
	//     o0666   O0666
	//     0o666   0O666
	//     0o0666  0O0666
	//     0666
	OctalLeaders []string `json:"octalLeaderChars"`

	// HexLeaders defines the prefixes used to signify that
	// a value should be parsed as hexadecimal.
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
	HexLeaders []string `json:"hexLeaderChars"`

	// The integer base to use when no prefix is present.
	//
	// Default: base 10
	DefaultBase int `json:"defaultBase"`
}

// UnmarshalMapProps defines map specific unmarshalling options for the "magic"
// unmarshaler.
type UnmarshalMapProps struct {

	// KeyValSeparatorChars defines characters used to
	// separate a key from a value in an individual mapping
	// entry.
	//
	// This character can be escaped with a '\' (backslash)
	// character.
	//
	// The first unescaped instance of one of the defined
	// characters in the individual entry will be used as
	// the divider, and any subsequent appearances in the
	// entry will be included in the value.
	//
	// Default: "=:" (equals, colon)
	//
	// Example key/value pairs using the default divider
	// characters.  The second column is a JSON
	// representation of the parsed map
	//
	//     key:value            {"key": "value"}
	//     key=value            {"key": "value"}
	//     key\\:foo:value      {"key:foo": "value"}
	//     key\\=bar=value      {"key=bar": "value"}
	//     key\\:baz=value:a=b  {"key:baz": "value:a=b"}
	//     key:value=c:d        {"key": "value=c:d"}
	KeyValSeparatorChars string

	// Default: ",; " (comma, semicolon, space)
	EntrySeparatorChars string
}

// UnmarshalSliceProps defines options for slice unmarshalling in the "magic"
// unmarshaler.
type UnmarshalSliceProps struct {
}

type UnmarshalTimeProps struct {
	// DateFormats configures the date-time formats that the unmarshaler will use
	// when attempting to parse a date value.
	//
	// By default, the RFC3339 and RFC3339 nano patterns are used.
	DateFormats []string
}

// ////////////////////////////////////////////////////////////////////////// //

// DefaultUnmarshalProps returns an UnmarshalProps instance with the default
// values configured.
func DefaultUnmarshalProps() UnmarshalProps {
	return defaultUnmarshalProps
}

var defaultUnmarshalProps = UnmarshalProps{
	Integers: UnmarshalIntegerProps{
		OctalLeaders: []string{"0o", "0O", "o", "O", "0"},
		HexLeaders:   []string{"0x", "0X", "x", "X"},
		DefaultBase:  10,
	},
	Maps: UnmarshalMapProps{
		KeyValSeparatorChars: "=:",
		EntrySeparatorChars:  ",; ",
	},
	Slices: UnmarshalSliceProps{},
	Time: UnmarshalTimeProps{
		DateFormats: []string{time.RFC3339, time.RFC3339Nano},
	},
}
