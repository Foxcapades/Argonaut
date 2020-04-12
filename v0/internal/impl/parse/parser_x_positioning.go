package parse

import . "github.com/Foxcapades/Argonaut/v0/internal/log"

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Parser Internal API: Positioning                                    ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// Increments the current character index and returns whether or not we've
// passed the end of the arg string
func (p *Parser) nextChar() (out bool) {
	TraceStart("Parser.nextChar")
	defer TraceEnd(func() []interface{} { return []interface{}{out} })

	p.charI++
	out = p.charI < p.strLen()
	return
}

// Increments the current argument index, resets the current character index and
// returns whether or not we've passed the end of the arg list
func (p *Parser) nextArg() (out bool) {
	p.argI++
	p.charI = 0
	out = p.argI < len(p.input)
	return
}
