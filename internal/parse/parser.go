package parse

import (
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/event"
)

type ChunkConsumer = func(Element)

func NewParser(e event.Emitter) Parser {
	return Parser{emitter: e}
}

type Parser struct {
	sb      strings.Builder
	state   state
	emitter event.Emitter
}

func (p *Parser) Next() Element {
	next := p.emitter.Next()

	if p.state == statePass {
		return p.consumeString(next.Data)
	}

	switch next.Kind {
	case event.KindDash:
		return p.handleDash()
	case event.KindText:
		return p.handleText(next.Data)
	case event.KindEnd:
		return p.handleEnd()
	default:
		panic("illegal state")
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
		case event.KindBreak:

			// If there were 2 dashes specifically, then we hit our boundary
			if dashes == 2 {
				p.state = statePass
				return boundaryElement()
			} else {
				for i := 0; i < dashes; i++ {
					p.sb.WriteByte(chars.CharDash)
				}
				tmp := p.sb.String()
				p.sb.Reset()
				return textElement(tmp)
			}

		case event.KindText:
			if dashes > 1 {
				return p.consumeLongFlag(dashes, next.Data)
			} else {
				return p.consumeShortFlag(next.Data)
			}

		case event.KindDash:
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

		if next.Kind == event.KindBreak {
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

		if next.Kind == event.KindBreak {
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
	case event.KindBreak:
		return longSoloElement(name)

	// We hit an equals, meaning we have something like "--flag=value" as our
	// input arg.  This means we need to keep eating to get the flag argument.
	case event.KindEquals:
		// continue

	default:
		panic("illegal state: expected break or equals, got " + next.Kind.String())
	}

	// If we made it here, then we had an equals character and are now expecting
	// the flag argument value.
	next = p.emitter.Next()

	if next.Kind != event.KindText {
		panic("illegal state")
	}

	p.sb.WriteString(next.Data)

	// Now consume the trailing break.
	next = p.emitter.Next()
	if next.Kind != event.KindBreak {
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
		p.sb.WriteString(flags)

		for {
			next = p.emitter.Next()

			if next.Kind == event.KindBreak {
				tmp := p.sb.String()
				p.sb.Reset()

				return textElement(tmp)
			}

			p.sb.WriteString(next.Data)
		}
	}

	if next.Kind == event.KindBreak {
		return shortSoloElement(flags)
	}

	if next.Kind != event.KindEquals {
		panic("illegal state")
	}

	// Skip the equals
	next = p.emitter.Next()

	if next.Kind != event.KindText {
		panic("illegal state")
	}

	return shortPairElement(flags, next.Data)
}

func (p *Parser) handleText(data string) Element {
	p.sb.WriteString(data)

	for {
		next := p.emitter.Next()

		if next.Kind == event.KindBreak {
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
