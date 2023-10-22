package argo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseBool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true", "t", "yes", "y", "on", "1":
		return true, nil
	case "false", "f", "no", "n", "off", "0":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %s as bool", raw)
	}
}

func parseInt(v string, bits int, opt *UnmarshalIntegerProps) (int64, error) {
	var neg string
	// TODO: Wrap this error

	if v[0] == '-' {
		neg = "-"
		v = v[1:]
	}

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseInt(neg+v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseInt(neg+v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseInt(neg+v, 10, bits)
}

func parseUInt(v string, bits int, opt *UnmarshalIntegerProps) (uint64, error) {
	// TODO: Wrap this error

	for i := range opt.HexLeaders {
		if strings.HasPrefix(v, opt.HexLeaders[i]) {
			return strconv.ParseUint(v[len(opt.HexLeaders[i]):], 16, bits)
		}
	}

	for i := range opt.OctalLeaders {
		if strings.HasPrefix(v, opt.OctalLeaders[i]) {
			return strconv.ParseUint(v[len(opt.OctalLeaders[i]):], 8, bits)
		}
	}

	return strconv.ParseUint(v, 10, bits)
}

func parseMapEntry(parser *mapElementParser) (k string, v string, err error) {
	if parser.HasNext() {
		k = parser.Next()
	} else {
		err = errors.New("could not parse map key, no values left in input")
		return
	}

	if parser.HasNext() {
		v = parser.Next()
	} else {
		err = errors.New("could not parse map value, no values left in input")
		return
	}

	return
}

func newMapElementParser(props UnmarshalMapProps, rawValue string) mapElementParser {
	return mapElementParser{properties: props, raw: rawValue}
}

type mapElementParser struct {
	properties     UnmarshalMapProps
	raw            string
	position       int
	expectingValue bool
}

func (m mapElementParser) HasNext() bool {
	return m.position < len(m.raw)
}

const charBS byte = '\\'

// Next returns the next split string from the input.
//
// The values are returned one at a time, starting with the first key, moving
// on to the first value, then the second key, and so on.
//
// It is up to the interpreter to join the keys to the values however it sees
// fit.
func (m *mapElementParser) Next() string {
	// there are 2 states, expecting a key, and expecting a value.
	//
	// When expecting a key, an entry separator will be ignored.  When a key-value
	// separator is encountered, if it was not preceded by a backslash, it will
	// break the raw string and shift the state to expecting a value.
	//
	// When expecting a value, a key-value separator will be ignored.  When an
	// entry separator is encountered, if it was not preceded by a backslash, it
	// will break the raw string and shift the state back to expecting a key.
	//
	// If a backslash is encountered it will be counted and a flag will be set.
	// When parsing the next byte from the input, if the flag is set, it may apply
	// to the value being parsed if it is a key-value separator or an entry
	// separator (depending on the parser state).
	if m.expectingValue {
		return m.parseValue()
	} else {
		return m.parseKey()
	}
}

func (m *mapElementParser) parseValue() string {
	return m.parse(m.isEntrySeparator)
}

func (m *mapElementParser) parseKey() string {
	return m.parse(m.isKeyValSeparator)
}

func (m *mapElementParser) parse(test func(byte) bool) string {
	start := m.position
	lastWasBS := false

	for m.position < len(m.raw) {
		c := m.raw[m.position]

		if test(c) {
			if lastWasBS {
				m.position++
			} else {
				out := m.raw[start:m.position]
				m.position++
				m.expectingValue = !m.expectingValue
				return out
			}
		} else {
			lastWasBS = c == charBS
			m.position++
		}
	}

	if start < m.position {
		m.expectingValue = !m.expectingValue
		return m.raw[start:m.position]
	}

	m.expectingValue = !m.expectingValue
	return ""
}

func (m mapElementParser) isKeyValSeparator(c byte) bool {
	for i := range m.properties.KeyValSeparatorChars {
		if c == m.properties.KeyValSeparatorChars[i] {
			return true
		}
	}
	return false
}

func (m mapElementParser) isEntrySeparator(c byte) bool {
	for i := range m.properties.EntrySeparatorChars {
		if c == m.properties.EntrySeparatorChars[i] {
			return true
		}
	}
	return false
}
