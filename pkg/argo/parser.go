package argo

import (
	"strings"
)

type state uint8

const (
	stateNone state = iota
	statePass
)

type elementType uint8

const (
	// elementTypeLongFlagPair represents a longform flag with a value directly
	// attached by use of an equals ('=') character.
	elementTypeLongFlagPair elementType = iota

	// elementTypeLongFlagSolo represents a longform flag that did not have a value
	// directly attached by use of an equals ('=') character.
	elementTypeLongFlagSolo

	// elementTypeShortBlockSolo represents a group of one or more characters
	// following a single dash ('-') character with no value directly attached via
	// an equals ('=') character.
	elementTypeShortBlockSolo

	// elementTypeShortBlockPair represents a group of one or more characters
	// following a single dash ('-') character with a value directly attached via
	// an equals ('=') character.
	elementTypeShortBlockPair

	// elementTypePlainText represents a plain-text argument that has no flag
	// indicators
	elementTypePlainText

	elementTypeBoundary

	elementTypeEnd
)

func newParser(e emitter) parser {
	return parser{emitter: e}
}

type parser struct {
	sb      strings.Builder
	state   state
	emitter emitter
}

func (p *parser) Next() element {
	next := p.emitter.Next()

	if p.state == statePass {
		return p.consumeString(next.Data)
	}

	switch next.Kind {
	case eventKindDash:
		return p.handleDash()
	case eventKindText:
		return p.handleText(next.Data)
	case eventKindEnd:
		return p.handleEnd()
	default:
		panic("illegal state")
	}
}

func (p *parser) handleDash() element {
	// Things that may follow a dash:
	// - Break
	// - Text
	// - Dash

	dashes := 1

	for {
		next := p.emitter.Next()

		switch next.Kind {

		// If we hit a break, then all we've seen are dashes (because we return on
		// anything else).
		case eventKindBreak:

			// If there were 2 dashes specifically, then we hit our boundary
			if dashes == 2 {
				p.state = statePass
				return boundaryElement()
			} else {
				for i := 0; i < dashes; i++ {
					p.sb.WriteByte(charDash)
				}
				tmp := p.sb.String()
				p.sb.Reset()
				return textElement(tmp)
			}

		case eventKindText:
			if dashes > 1 {
				return p.consumeLongFlag(dashes, next.Data)
			} else {
				return p.consumeShortFlag(next.Data)
			}

		case eventKindDash:
			dashes++

		default:
			panic("illegal state")
		}
	}
}

func (p *parser) consumeString(start string) element {
	p.sb.WriteString(start)

	for {
		next := p.emitter.Next()

		if next.Kind == eventKindBreak {
			tmp := p.sb.String()
			p.sb.Reset()
			return textElement(tmp)
		}

		p.sb.WriteString(next.Data)
	}
}

func (p *parser) consumeLongFlag(dashes int, name string) element {
	// If there is a whitespace in the name string, then it's not actually a flag
	// it's just a string that happened to start with "--", which is stupid, but
	// what are you gonna do?
	if idx := nextWhitespace(name); idx > -1 {
		for i := 0; i < dashes; i++ {
			p.sb.WriteByte(charDash)
		}

		p.sb.WriteString(name)

		next := p.emitter.Next()

		if next.Kind == eventKindBreak {
			tmp := p.sb.String()
			p.sb.Reset()
			p.state = stateNone
			return textElement(tmp)
		}

		panic("illegal state: expected break, but got " + next.Kind.String())
	}

	// So, it actually is a long flag.
	next := p.emitter.Next()

	// The valid things that could happen now are:
	// - equals
	// - break

	switch next.Kind {

	// We hit a break, meaning we have something like "--flag" as our input arg.
	case eventKindBreak:
		return longSoloElement(name)

	// We hit an equals, meaning we have something like "--flag=value" as our
	// input arg.  This means we need to keep eating to get the flag argument.
	case eventKindEquals:
		// continue

	default:
		panic("illegal state: expected break or equals, got " + next.Kind.String())
	}

	// If we made it here, then we had an equals character and are now expecting
	// the flag argument value.
	next = p.emitter.Next()

	if next.Kind != eventKindText {
		panic("illegal state")
	}

	p.sb.WriteString(next.Data)

	// Now consume the trailing break.
	next = p.emitter.Next()
	if next.Kind != eventKindBreak {
		panic("illegal state")
	}

	tmp := p.sb.String()
	p.sb.Reset()

	return longPairElement(name, tmp)
}

func (p *parser) consumeShortFlag(flags string) element {
	next := p.emitter.Next()

	// If there is a whitespace character in the middle of this flag group then it
	// isn't really a flag group at all.
	if idx := nextWhitespace(flags); idx > -1 {
		p.sb.WriteString(flags)

		for {
			next = p.emitter.Next()

			if next.Kind == eventKindBreak {
				tmp := p.sb.String()
				p.sb.Reset()

				return textElement(tmp)
			}

			p.sb.WriteString(next.Data)
		}
	}

	if next.Kind == eventKindBreak {
		return shortSoloElement(flags)
	}

	if next.Kind != eventKindEquals {
		panic("illegal state")
	}

	// Skip the equals
	next = p.emitter.Next()

	if next.Kind != eventKindText {
		panic("illegal state")
	}

	return shortPairElement(flags, next.Data)
}

func (p *parser) handleText(data string) element {
	p.sb.WriteString(data)

	for {
		next := p.emitter.Next()

		if next.Kind == eventKindBreak {
			tmp := p.sb.String()
			p.sb.Reset()

			return textElement(tmp)
		}

		p.sb.WriteString(next.Data)
	}
}

func (p *parser) handleEnd() element {
	return endElement()
}

type element struct {
	Type elementType
	Data []string
}

func (e element) String() string {
	switch e.Type {
	case elementTypeLongFlagPair:
		return "--" + e.Data[0] + "=" + e.Data[1]
	case elementTypeShortBlockPair:
		return "-" + e.Data[0] + "=" + e.Data[1]
	case elementTypeLongFlagSolo:
		return "--" + e.Data[0]
	case elementTypeShortBlockSolo:
		return "-" + e.Data[0]
	case elementTypeBoundary:
		return "--"
	case elementTypePlainText:
		return e.Data[0]
	case elementTypeEnd:
		return string(byte(0))
	default:
		panic("illegal state")
	}
}

func longPairElement(flag, value string) element {
	return element{Type: elementTypeLongFlagPair, Data: []string{flag, value}}
}

func longSoloElement(flag string) element {
	return element{Type: elementTypeLongFlagSolo, Data: []string{flag}}
}

func shortPairElement(flags, value string) element {
	return element{Type: elementTypeShortBlockPair, Data: []string{flags, value}}
}

func shortSoloElement(flags string) element {
	return element{Type: elementTypeShortBlockSolo, Data: []string{flags}}
}

func textElement(text string) element {
	return element{Type: elementTypePlainText, Data: []string{text}}
}

func boundaryElement() element {
	return element{Type: elementTypeBoundary}
}

func endElement() element {
	return element{Type: elementTypeEnd}
}
