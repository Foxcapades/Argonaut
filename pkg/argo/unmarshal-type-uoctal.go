package argo

import "strconv"

// UOctal represents an untyped unsigned int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type UOctal uint

// Unmarshal implements the Unmarshaler.Unmarshal method for the UOctal type.
func (o *UOctal) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 8, strconv.IntSize)
	*o = UOctal(tmp)
	return
}

// UOctal8 represents an unsigned 8 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type UOctal8 uint8

// Unmarshal implements the Unmarshaler.Unmarshal method for the UOctal8 type.
func (o *UOctal8) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 8, 8)
	*o = UOctal8(tmp)
	return
}

// UOctal16 represents an unsigned 16 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type UOctal16 uint16

// Unmarshal implements the Unmarshaler.Unmarshal method for the UOctal16 type.
func (o *UOctal16) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 8, 16)
	*o = UOctal16(tmp)
	return
}

// UOctal32 represents an unsigned 32 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type UOctal32 uint32

// Unmarshal implements the Unmarshaler.Unmarshal method for the UOctal32 type.
func (o *UOctal32) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 8, 32)
	*o = UOctal32(tmp)
	return
}

// UOctal64 represents an unsigned 64 bit int value that is
// expected to be input in octal notation and will be parsed
// from string in base 8.
type UOctal64 uint64

// Unmarshal implements the Unmarshaler.Unmarshal method for the UOctal64 type.
func (o *UOctal64) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 8, 64)
	*o = UOctal64(tmp)
	return
}
