package parse

import (
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/emit"
)

type state uint8

const (
	stateNone state = iota
	statePass
)

type ElementType uint8

const (
	// ElementTypeLongFlagPair represents a longform flag with a value directly
	// attached by use of an equals ('=') character.
	ElementTypeLongFlagPair ElementType = iota

	// ElementTypeLongFlagSolo represents a longform flag that did not have a value
	// directly attached by use of an equals ('=') character.
	ElementTypeLongFlagSolo

	// ElementTypeShortBlockSolo represents a group of one or more characters
	// following a single dash ('-') character with no value directly attached via
	// an equals ('=') character.
	ElementTypeShortBlockSolo

	// ElementTypeShortBlockPair represents a group of one or more characters
	// following a single dash ('-') character with a value directly attached via
	// an equals ('=') character.
	ElementTypeShortBlockPair

	// ElementTypePlainText represents a plain-text argument that has no flag
	// indicators
	ElementTypePlainText

	ElementTypeBoundary

	ElementTypeEnd
)

func NewParser(e emit.Emitter) Parser {
	return Parser{emitter: e}
}

type Parser struct {
	sb      strings.Builder
	state   state
	emitter emit.Emitter
}

func (p *Parser) Next() Element {
	next := p.emitter.Next()

	if p.state == statePass && next.Kind != emit.EventKindEnd {
		return p.consumeString(next.Data)
	}

	switch next.Kind {
	case emit.EventKindDash:
		return p.handleDash()
	case emit.EventKindText:
		return p.handleText(next.Data)
	case emit.EventKindEnd:
		return p.handleEnd()
	default:
		panic("illegal state: got event kind " + next.Kind.String())
	}
}

func (p *Parser) handleDash() Element {
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
		case emit.EventKindBreak:

			// If there were 2 dashes specifically, then we hit our boundary
			if dashes == 2 {
				p.state = statePass
				return boundaryElement()
			}

			for i := 0; i < dashes; i++ {
				p.sb.WriteByte(chars.CharDash)
			}

			tmp := p.sb.String()
			p.sb.Reset()
			return textElement(tmp)

		case emit.EventKindText:
			if dashes > 1 {
				return p.consumeLongFlag(dashes, next.Data)
			}

			return p.consumeShortFlag(next.Data)

		case emit.EventKindDash:
			dashes++

		default:
			panic("illegal state")
		}
	}
}

func (p *Parser) consumeString(start string) Element {
	p.sb.WriteString(start)

	for {
		next := p.emitter.Next()

		if next.Kind == emit.EventKindBreak {
			tmp := p.sb.String()
			p.sb.Reset()
			return textElement(tmp)
		}

		p.sb.WriteString(next.Data)
	}
}

func (p *Parser) consumeLongFlag(dashes int, name string) Element {
	// If there is a whitespace in the name string, then it's not actually a flag
	// it's just a string that happened to start with "--", which is stupid, but
	// what are you gonna do?
	if idx := chars.NextWhitespace(name); idx > -1 {
		for i := 0; i < dashes; i++ {
			p.sb.WriteByte(chars.CharDash)
		}

		p.sb.WriteString(name)

		next := p.emitter.Next()

		if next.Kind == emit.EventKindBreak {
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
	case emit.EventKindBreak:
		return longSoloElement(name)

	// We hit an equals, meaning we have something like "--flag=value" as our
	// input arg.  This means we need to keep eating to get the flag argument.
	case emit.EventKindEquals:
		// continue

	default:
		panic("illegal state: expected break or equals, got " + next.Kind.String())
	}

	// If we made it here, then we had an equals character and are now expecting
	// the flag argument value.
	next = p.emitter.Next()

	if next.Kind != emit.EventKindText {
		panic("illegal state")
	}

	p.sb.WriteString(next.Data)

	// Now consume the trailing break.
	next = p.emitter.Next()
	if next.Kind != emit.EventKindBreak {
		panic("illegal state")
	}

	tmp := p.sb.String()
	p.sb.Reset()

	return longPairElement(name, tmp)
}

func (p *Parser) consumeShortFlag(flags string) Element {
	next := p.emitter.Next()

	// If there is a whitespace character in the middle of this flag group then it
	// isn't really a flag group at all.
	if idx := chars.NextWhitespace(flags); idx > -1 {
		p.sb.WriteByte(chars.CharDash)
		p.sb.WriteString(flags)

		for {
			next = p.emitter.Next()

			if next.Kind == emit.EventKindBreak || next.Kind == emit.EventKindEnd {
				tmp := p.sb.String()
				p.sb.Reset()

				return textElement(tmp)
			}

			p.sb.WriteString(next.Data)
		}
	}

	if next.Kind == emit.EventKindBreak {
		return shortSoloElement(flags)
	}

	if next.Kind != emit.EventKindEquals {
		panic("illegal state")
	}

	// Skip the equals
	next = p.emitter.Next()

	if next.Kind != emit.EventKindText {
		panic("illegal state")
	}

	data := next.Data

	next = p.emitter.Next()
	if next.Kind != emit.EventKindBreak {
		panic("illegal state: expected break, got " + next.Kind.String())
	}

	return shortPairElement(flags, data)
}

func (p *Parser) handleText(data string) Element {
	p.sb.WriteString(data)

	// Loop here because we may be followed by a break, or an equals event which
	// itself may be followed by more text.
	for {
		next := p.emitter.Next()

		if next.Kind == emit.EventKindBreak {
			tmp := p.sb.String()
			p.sb.Reset()

			return textElement(tmp)
		}

		p.sb.WriteString(next.Data)
	}
}

func (p *Parser) handleEnd() Element {
	return endElement()
}

type Element struct {
	Type ElementType
	Data []string
}

func (e Element) String() string {
	switch e.Type {
	case ElementTypeLongFlagPair:
		return "--" + e.Data[0] + "=" + e.Data[1]
	case ElementTypeShortBlockPair:
		return "-" + e.Data[0] + "=" + e.Data[1]
	case ElementTypeLongFlagSolo:
		return "--" + e.Data[0]
	case ElementTypeShortBlockSolo:
		return "-" + e.Data[0]
	case ElementTypeBoundary:
		return "--"
	case ElementTypePlainText:
		return e.Data[0]
	case ElementTypeEnd:
		return string(byte(0))
	default:
		panic("illegal state")
	}
}

func longPairElement(flag, value string) Element {
	return Element{Type: ElementTypeLongFlagPair, Data: []string{flag, value}}
}

func longSoloElement(flag string) Element {
	return Element{Type: ElementTypeLongFlagSolo, Data: []string{flag}}
}

func shortPairElement(flags, value string) Element {
	return Element{Type: ElementTypeShortBlockPair, Data: []string{flags, value}}
}

func shortSoloElement(flags string) Element {
	return Element{Type: ElementTypeShortBlockSolo, Data: []string{flags}}
}

func textElement(text string) Element {
	return Element{Type: ElementTypePlainText, Data: []string{text}}
}

func boundaryElement() Element {
	return Element{Type: ElementTypeBoundary}
}

func endElement() Element {
	return Element{Type: ElementTypeEnd}
}
