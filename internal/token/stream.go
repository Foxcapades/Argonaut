package token

import (
	"strings"

	"github.com/foxcapades/argonaut/internal/util/xstr"
	"github.com/foxcapades/argonaut/pkg/argo"
)

func NewStream(inputs []string, options *argo.Config) Stream {
	return Stream{rawInput: inputs, options: options}
}

// Stream produces a stream of Token values by doing very basic examination of
// the CLI input strings.
//
// Once the input stream is exhausted, repeated usages of the Stream will just
// return end of input tokens.
type Stream struct {
	rawInput []string
	rawIndex uint32

	reachedEndOfOptions bool

	options *argo.Config
}

// Next returns the next token in the CLI input strings.
//
// If an end-of-options marker is encountered, all following input strings will
// be marked as arguments without examination.
//
// Once the input stream has been exhausted, this method will only return end of
// input tokens.
func (s *Stream) Next() Token {
	if s.rawIndex >= s.rawLen() {
		return Token{Type: TypeEndOfInput}
	}

	tmp := s.rawInput[s.rawIndex]
	s.rawIndex++

	if s.reachedEndOfOptions {
		return Token{TypeArgument, tmp, 0}
	}

	// if the length of the string is 0, it's an empty argument
	if len(tmp) < 1 {
		return Token{TypeArgument, tmp, 0}
	}

	// if the value is the end of options marker
	if tmp == s.options.EndOfOptionsMarker {
		s.reachedEndOfOptions = true
		return Token{TypeEndOfOptions, tmp, 0}
	}

	if strings.HasPrefix(tmp, s.options.LongFlagPrefix) {
		pos := xstr.IndexOfByte(tmp, s.options.LongFlagValueSeparator, len(s.options.LongFlagPrefix))

		if pos < 0 || pos > 65535 {
			pos = 0
		}

		return Token{TypeLong, tmp, uint16(pos)}
	}

	if tmp[0] == s.options.ShortFlagPrefix && len(tmp) > 1 {
		pos := xstr.IndexOfByte(tmp, s.options.LongFlagValueSeparator, 1)

		if pos < 0 || pos > 65535 {
			pos = 0
		}

		return Token{TypeShort, tmp, uint16(pos)}
	}

	return Token{TypeArgument, tmp, 0}
}

func (s *Stream) rawLen() uint32 {
	return uint32(len(s.rawInput))
}
