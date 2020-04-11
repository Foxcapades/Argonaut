package parse

import (
	"fmt"
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
)

func NewParser() A.Parser {
	return new(Parser)
}

type Parser struct {
	shorts map[byte]A.Flag
	longs  map[string]A.Flag
	reqs   map[uintptr]interface{}

	com A.Command

	args []A.Argument

	// CLI input
	input []string

	// Unrecognized params
	extra []string

	// values after the "--" end of flags marker.
	passthrough []string

	// index of the current argument
	argI int

	// index of the current char in the current argument
	charI int

	waiting A.Flag
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Public API                                                          ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) Parse(args []string, command A.Command) (err error) {
	TraceStart("Parser.Parse", args, command)
	defer TraceEnd(func() []interface{} { return []interface{}{err} })

	defer recovery(&err)
	p.setup(args, command)

	// Skip first arg (it's the command name)
	for p.nextArg() {
		p.parseNext()
	}

	p.complete()

	return
}

func (p *Parser) Unrecognized() []string {
	return p.extra
}

func (p *Parser) Passthroughs() []string {
	return p.passthrough
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Parse Meta                                            ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) parseNext() {
	TraceStart("Parser.parseNext")
	defer TraceEnd(func() []interface{} { return nil })

	// We were intentionally given an empty string as an
	// argument, assign it and move on
	if p.strLen() == 0 {
		p.handleArg()
		return
	}

	// Handle required argument
	if p.waiting != nil && p.waiting.Argument().Required() {
		ptr := pointerFor(p.waiting) // get the pointer before we pop it
		p.unmarshal(p.popArg())
		delete(p.reqs, ptr)
		return
	}

	// If we've made it here, we know that the following must
	// be true:
	//   - the input param value is non-empty
	//   - we do not have a waiting required flag argument
	//   - the input type (flag vs arg) is as of yet
	//     undetermined
	// And the following may be true
	//   - we have a waiting optional flag argument

	// Input is an argument
	if p.char() != '-' {
		Trace("current param did not start with dash")
		p.handleArg()
		return
	}

	if p.waiting != nil && p.isBoolArg(p.waiting.Argument()) {
		util.Must(p.com.Unmarshaler().Unmarshal("true", p.popArg().Binding()))
	}

	// if we've made it this far we know that:
	//   - we have no flags to fill
	//   - the input is one of the following
	//     - just a dash ("-")
	//     - one or more short flags
	//     - a long flag
	//     - the end of flags marker ("--")
	//     - invalid

	// input was just "-"
	if !p.nextChar() {
		Trace("current param was just '-'")
		p.handleArg()
		return
	}

	// Dash followed by some other character
	if p.char() != '-' {
		Trace("current param is short flag")
		p.handleShortFlag()
		return
	}

	// Input must be:
	//   - a long flag
	//   - end of flags marker
	//   - invalid

	// End of flags marker
	if !p.nextChar() {
		for p.nextArg() {
			p.passthrough = append(p.passthrough, p.argument())
		}
		return
	}

	// Input must be:
	//   - a long flag
	//   - invalid

	// triple dash?
	if p.char() == '-' {
		p.handleArg()
		return
	}

	p.handleLongFlag()
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Parse Params                                          ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) handleArg() {
	TraceStart("Parser.handleArg")
	defer TraceEnd(func() []interface{} { return nil })

	if arg := p.popArg(); arg != nil {
		Trace("popped arg")
		p.unmarshal(arg)
		delete(p.reqs, pointerFor(arg))
	} else {
		Trace("no arg, extra input")
		p.extra = append(p.extra, p.argument())
	}
}

func (p *Parser) isBoolArg(arg A.Argument) (out bool) {
	TraceStart("isBoolArg")
	defer TraceEnd(func() []interface{} { return []interface{}{out} })

	bt := arg.BindingType().String()

	// Special case for boolean arguments.  The existence of
	// the flag can itself be a true value for a boolean arg.
	out = bt == "bool" || bt == "*bool" || bt == "[]bool" || bt == "[]*bool"
	return
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Helpers                                                             ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func pointerFor(v interface{}) uintptr {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Errorf("Failed to get the backing pointer for %s: %s", v, err))
		}
	}()
	defer recovery(&err)

	return R.ValueOf(v).Pointer()
}
