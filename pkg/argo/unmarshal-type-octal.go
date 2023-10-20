package argo

import "strconv"

// Octal represents an untyped signed int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type Octal int

// Unmarshal implements the Unmarshaler.Unmarshal method for the Octal type.
func (o *Octal) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 8, strconv.IntSize)
	*o = Octal(tmp)
	return
}

// Octal8 represents a signed 8 bit int value that is expected
// to be input in octal notation and will be parsed from
// string in base 8.
type Octal8 int8

// Unmarshal implements the Unmarshaler.Unmarshal method for the Octal8 type.
func (o *Octal8) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 8, 8)
	*o = Octal8(tmp)
	return
}

// Octal16 represents a signed 16 bit int value that
// is expected to be input in octal notation and will be
// parsed from string in base 8.
type Octal16 int16

// Unmarshal implements the Unmarshaler.Unmarshal method for the Octal16 type.
func (o *Octal16) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 8, 16)
	*o = Octal16(tmp)
	return
}

// Octal32 represents a signed 32 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type Octal32 int32

// Unmarshal implements the Unmarshaler.Unmarshal method for the Octal32 type.
func (o *Octal32) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 8, 32)
	*o = Octal32(tmp)
	return
}

// Octal64 represents a signed 64 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type Octal64 int64

// Unmarshal implements the Unmarshaler.Unmarshal method for the Octal64 type.
func (o *Octal64) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseInt(value, 8, 64)
	*o = Octal64(tmp)
	return
}
