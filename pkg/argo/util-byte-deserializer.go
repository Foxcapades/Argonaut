package argo

import (
	"encoding/base64"
	"encoding/hex"
)

// ByteSliceParser defines a function that expects a string input and attempts
// to deserialize that input into a byte slice.
type ByteSliceParser = func(string) ([]byte, error)

// ByteSliceParserRaw is a ByteSliceParser function that casts the input string
// into a byte slice.
//
// This function will never return an error.
func ByteSliceParserRaw(s string) ([]byte, error) {
	return []byte(s), nil
}

// ByteSliceParserBase64 is a ByteSliceParser function that decodes the given
// string as a base64 string.
func ByteSliceParserBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// ByteSliceParserHex is a ByteSliceParser function that decodes the given
// string as a hex string.
func ByteSliceParserHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
