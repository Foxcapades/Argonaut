package parse

import (
	"fmt"
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
	"strings"
)

func NewParser() A.Parser {
	return new(Parser)
}

type Parser struct {
	shorts map[byte]A.Flag
	longs  map[string]A.Flag
	reqs   map[uintptr]interface{}

	com A.Command

	args   []A.Argument

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

func (p *Parser) complete() {
	TraceStart("Parser.complete")
	defer TraceEnd(func() []interface{} { return nil })

	if p.waiting != nil {
		if p.isBoolArg(p.waiting.Argument()) {
			util.Must(p.com.Unmarshaler().Unmarshal("true", p.popArg().Binding()))
		} else {
			// TODO: make this a real error
			panic("missing required arg")
		}
	}

	if len(p.reqs) > 0 {
		// TODO: make this a real error
		panic("missing required params")
	}

	// assign defaults
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

// Process a short flag param
func (p *Parser) handleShortFlag() {
	TraceStart("Parser.handleShortFlag")
	defer TraceEnd(func() []interface{} { return nil })

	flag, ok := p.shorts[p.char()]
	Trace(flag, ok)

	// Invalid flag character
	if !ok {
		p.extra = append(p.extra, string([]byte{'-', p.char()}))
		if p.nextChar() {
			p.handleShortFlag()
		}
		return
	}

	// Clear this so we don't error on complete for missing
	// required argument
	delete(p.reqs, pointerFor(flag))

	flag.IncrementHits()

	// no argument, moving on
	if !flag.HasArgument() {
		if p.nextChar() {
			p.handleShortFlag()
		}
		return
	}

	arg := flag.Argument()

	if !p.nextChar() {
		Trace("ending with waiting flag")
		p.waiting = flag
		return
	}

	// Argument is required
	if arg.Required() {
		Trace("required arg")

		if p.isBoolArg(arg) {
			Trace("argument is bool")
			if util.IsBool(p.argument()[p.charI:]) {
				Trace("input is also bool")
				util.Must(p.com.Unmarshaler().Unmarshal(p.eatString(), arg.Binding()))
			} else {
				Trace("input is not bool")
				util.Must(p.com.Unmarshaler().Unmarshal("true", arg.Binding()))
				p.handleShortFlag()
			}
			return
		}
		util.Must(p.com.Unmarshaler().Unmarshal(p.eatString(), arg.Binding()))
		return
	}

	Trace("short flag arg is optional")

	if _, ok := p.shorts[p.char()]; ok {
		Trace("next char is flag")
		if p.isBoolArg(arg) {
			Trace("cur flag argument is bool")
			util.Must(p.com.Unmarshaler().Unmarshal("true", arg.Binding()))
		}
		p.handleShortFlag()
	} else {
		util.Must(p.com.Unmarshaler().Unmarshal(p.eatString(), arg.Binding()))
	}
}

func (p *Parser) handleLongFlag() {
	TraceStart("Parser.handleLongFlag")
	defer TraceEnd(func() []interface{} { return nil })

	rest := p.eatString()
	split := strings.SplitN(rest, "=", 2)

	flag, ok := p.longs[split[0]]

	// Unrecognized flag
	if !ok {
		p.extra = append(p.extra, rest)
		return
	}

	flag.IncrementHits()
	delete(p.reqs, pointerFor(flag))

	// Flag takes no argument
	if !flag.HasArgument() {
		if len(split) > 1 {
			// TODO: Make this a real error
			panic("flag " + flag.String() + " does not expect an argument")
		}
		return
	}

	arg := flag.Argument()

	// No argument provided
	if len(split) == 1 {
		if arg.Required() {
			// Boolean case
			if p.isBoolArg(arg) {
				util.Must(p.com.Unmarshaler().Unmarshal("true", arg.Binding()))
				return
			}

			// TODO: make this a real error
			panic("flag " + flag.String() + " requires an argument")
		}
		return
	}

	util.Must(p.com.Unmarshaler().Unmarshal(split[1], arg.Binding()))
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
//┃      Internal API: Positioning                                           ┃//
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


//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Helpers                                               ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) setup(args []string, com A.Command) {
	TraceStart("Parser.setup", args, com)
	defer TraceEnd(func() []interface{} { return nil })
	p.makeMaps(com)
	p.args        = com.Arguments()
	p.input       = args
	p.com         = com
	p.extra       = nil
	p.passthrough = nil
	p.argI        = 0
	p.charI       = 0
	p.waiting     = nil
}

func (p *Parser) popArg() (arg A.Argument) {
	TraceStart("Parser.popArg")
	defer TraceEnd(func() []interface{} { return []interface{}{arg} })

	if p.waiting != nil {
		arg = p.waiting.Argument()
		p.waiting = nil
	} else {
		if len(p.args) > 0 {
			arg = p.args[0]
			p.args = p.args[1:]
		}
	}
	return
}

func (p *Parser) eatString() (out string) {
	TraceStart("Parser.eatString")
	defer TraceEnd(func() []interface{} { return []interface{}{out} })
	out = p.argument()[p.charI:]
	return
}

// Returns the length of the current argument string
func (p *Parser) strLen() int {
	return len(p.argument())
}

// returns the character at the current arg and char index
func (p *Parser) char() byte {
	return p.argument()[p.charI]
}

// returns the argument at the current arg index
func (p *Parser) argument() string {
	return p.input[p.argI]
}

// Create and populate parser maps
func (p *Parser) makeMaps(command A.Command) {
	TraceStart("Parser.makeMaps", command)
	defer TraceEnd(func() []interface{} { return nil })

	p.shorts = make(map[byte]A.Flag)
	p.longs = make(map[string]A.Flag)
	p.reqs = make(map[uintptr]interface{})

	for _, group := range command.FlagGroups() {
		flags := group.Flags()
		for i := range flags {
			if flags[i].HasShort() {
				p.shorts[flags[i].Short()] = flags[i]
			}
			if flags[i].HasLong() {
				p.longs[flags[i].Long()] = flags[i]
			}
			if flags[i].Required() {
				p.reqs[pointerFor(flags[i])] = flags[i]
			}
		}
	}

	args := command.Arguments()
	for i := range args {
		if args[i].Required() {
			p.reqs[pointerFor(args[i])] = args[i]
		}
	}
}

func (p *Parser) unmarshal(arg A.Argument) {
	TraceStart("Parser.unmarshal", arg)
	defer TraceEnd(func() []interface{} { return nil })

	bind := arg.Binding()
	kind := util.GetRootValue(R.ValueOf(bind))

	Trace(arg.Parent())

	// If binding is specialized
	if cst, ok := kind.Interface().(A.SpecializedUnmarshaler); ok {
		util.Must(p.com.Unmarshaler().Unmarshal(p.argument(), bind))
		if !cst.ConsumesArguments() {
			p.argI--
		}
		return
	}

	// if binding is bool, only consume the arg if it's
	// actually a valid bool value
	if p.isBoolArg(arg) {
		if util.IsBool(p.argument()) {
			util.Must(p.com.Unmarshaler().Unmarshal(p.argument(), bind))
		} else {
			util.Must(p.com.Unmarshaler().Unmarshal("true", bind))
			p.argI--
		}
		return
	}

	util.Must(p.com.Unmarshaler().Unmarshal(p.argument(), bind))
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Helpers                                                             ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func pointerFor(v interface{}) uintptr {
	var err error
	defer func() {if err != nil {
		panic(fmt.Errorf("Failed to get the backing pointer for %s: %s", v, err))
	}}()
	defer recovery(&err)

	return R.ValueOf(v).Pointer()
}
