package text_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/text"
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

		if text.IsAlpha(b) != okays[b] {
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

		if text.IsNumeric(b) != okays[b] {
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

		if text.IsAlphanumeric(b) != okays[b] {
			t.Errorf("expected %c to be alphanumeric but it wasn't", b)
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	okays := map[byte]bool{' ': true, '\t': true, '\n': true, '\r': true}

	for i := 0; i < 256; i++ {
		b := byte(i)

		if text.IsWhitespace(b) != okays[b] {
			t.Errorf("expected %c to be whitespace but it wasn't", b)
		}
	}
}

func TestIsBlank(t *testing.T) {
	if !text.IsBlank("") {
		t.Error("expected true but got false")
	}
	if !text.IsBlank(" ") {
		t.Errorf("expected true but got false")
	}
	if text.IsBlank("           3") {
		t.Errorf("expected false but got true")
	}
}
