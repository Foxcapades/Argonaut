package argo

import "strconv"

// UHex represents an untyped unsigned int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type UHex uint

// Unmarshal implements the Unmarshaler.Unmarshal method for the UHex type.
func (h *UHex) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 16, strconv.IntSize)
	*h = UHex(tmp)
	return
}

// UHex8 represents an unsigned 8 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type UHex8 uint8

// Unmarshal implements the Unmarshaler.Unmarshal method for the UHex8 type.
func (h *UHex8) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 16, 8)
	*h = UHex8(tmp)
	return
}

// UHex16 represents an unsigned 16 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type UHex16 uint16

// Unmarshal implements the Unmarshaler.Unmarshal method for the UHex16 type.
func (h *UHex16) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 16, 16)
	*h = UHex16(tmp)
	return
}

// UHex32 represents an unsigned 32 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type UHex32 uint32

// Unmarshal implements the Unmarshaler.Unmarshal method for the UHex32 type.
func (h *UHex32) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 16, 32)
	*h = UHex32(tmp)
	return
}

// UHex64 represents an unsigned 64 bit int value that is
// expected to be input in hexadecimal notation and will be
// parsed from string in base 16.
type UHex64 uint64

// Unmarshal implements the Unmarshaler.Unmarshal method for the UHex64 type.
func (h *UHex64) Unmarshal(value string) (err error) {
	tmp, err := strconv.ParseUint(value, 16, 64)
	*h = UHex64(tmp)
	return
}
