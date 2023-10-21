package chars_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

func TestIsAlpha(t *testing.T) {
	okays := map[byte]bool{
		'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true,
		'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true,
		'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true,
		'v': true, 'w': true, 'x': true, 'y': true, 'z': true,

		'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true,
		'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true,
		'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true,
		'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,
	}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if chars.IsAlpha(b) != okays[b] {
			t.Errorf("expected %c to be alpha but it wasn't", b)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	okays := map[byte]bool{
		'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true,
		'7': true, '8': true, '9': true,
	}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if chars.IsNumeric(b) != okays[b] {
			t.Errorf("expected %c to be numeric but it wasn't", b)
		}
	}
}

func TestIsAlphanumeric(t *testing.T) {
	okays := map[byte]bool{
		'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true,
		'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true,
		'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true,
		'v': true, 'w': true, 'x': true, 'y': true, 'z': true,

		'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true,
		'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true,
		'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true,
		'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,

		'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true,
		'7': true, '8': true, '9': true,
	}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if chars.IsAlphanumeric(b) != okays[b] {
			t.Errorf("expected %c to be alphanumeric but it wasn't", b)
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	okays := map[byte]bool{' ': true, '\t': true, '\n': true, '\r': true}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if chars.IsWhitespace(b) != okays[b] {
			t.Errorf("expected %c to be whitespace but it wasn't", b)
		}
	}
}

func TestIsFlagStringSafe(t *testing.T) {
	okays := map[byte]bool{
		'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true,
		'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true,
		'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true,
		'v': true, 'w': true, 'x': true, 'y': true, 'z': true,

		'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true,
		'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true,
		'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true,
		'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,

		'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true,
		'7': true, '8': true, '9': true,

		'-': true, '_': true,
	}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if chars.IsFlagStringSafe(b) != okays[b] {
			t.Errorf("expected %c to be flag safe but it wasn't", b)
		}
	}
}

func TestValidateCommandNodeName(t *testing.T) {
	if err := chars.ValidateCommandNodeName(""); err == nil {
		t.Error("expected an error but didn't get one")
	}

	if err := chars.ValidateCommandNodeName("-who"); err == nil {
		t.Error("expected an error but didn't get one")
	}

	if err := chars.ValidateCommandNodeName("_who"); err != nil {
		t.Error("expected no error but got one")
	}

	if err := chars.ValidateCommandNodeName("who\n"); err == nil {
		t.Error("expected an error but didn't get one")
	}
}

func TestNextWhitespace(t *testing.T) {
	if idx := chars.NextWhitespace(" "); idx != 0 {
		t.Errorf("expected idx to equal 0 but it was %d", idx)
	}

	if idx := chars.NextWhitespace("a "); idx != 1 {
		t.Errorf("expected idx to equal 1 but it was %d", idx)
	}

	if idx := chars.NextWhitespace("asdfasdf"); idx != -1 {
		t.Errorf("expected idx to equal -1 but it was %d", idx)
	}
}

func TestNextEquals(t *testing.T) {
	if idx := chars.NextEquals("="); idx != 0 {
		t.Errorf("expected idx to equal 0 but it was %d", idx)
	}

	if idx := chars.NextEquals("a="); idx != 1 {
		t.Errorf("expected idx to equal 1 but it was %d", idx)
	}

	if idx := chars.NextEquals("asdfasdf"); idx != -1 {
		t.Errorf("expected idx to equal -1 but it was %d", idx)
	}
}

func TestIsBlank(t *testing.T) {
	if !chars.IsBlank("") {
		t.Error("expected true but got false")
	}
	if !chars.IsBlank(" ") {
		t.Errorf("expected true but got false")
	}
	if chars.IsBlank("           3") {
		t.Errorf("expected false but got true")
	}
}
