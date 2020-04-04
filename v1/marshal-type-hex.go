package argo

import "strconv"

// TODO: Document special int value parsing rules around base

// Hex represents an untyped signed int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type Hex int

func (h *Hex) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 16, strconv.IntSize)
	*h = Hex(tmp)
	return
}

// Hex8 represents a signed 8 bit int value that is expected
// to be input in hexadecimal notation and will be parsed
// from string in base 16.
type Hex8 int8

func (h *Hex8) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 16, 8)
	*h = Hex8(tmp)
	return
}

// Hex16 represents a signed 16 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type Hex16 int16

func (h *Hex16) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 16, 16)
	*h = Hex16(tmp)
	return
}

// Hex32 represents a signed 32 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type Hex32 int32

func (h *Hex32) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 16, 32)
	*h = Hex32(tmp)
	return
}

// Hex64 represents a signed 64 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type Hex64 int64

func (h *Hex64) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 16, 64)
	*h = Hex64(tmp)
	return
}
