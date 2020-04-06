package impl

import (
	"errors"
	"fmt"
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
	//defer recovery(&err)
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

	// We were intentionally given an empty string as an
	// argument, assign it and move on
	if p.strLen() == 0 {
		p.handleArg()
		return
	}

	// Handle required argument
	if p.waiting != nil && p.waiting.Argument().Required() {
		p.handleRequiredArg()
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
		p.handleArg()
		return
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
		p.handleArg()
		return
	}

	// Dash followed by some other character
	if p.char() != '-' {
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
}


//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Parse Params                                          ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//


func (p *Parser) handleArg() {
	if arg := p.popArg(); arg != nil {
		util.Must(p.com.Unmarshaler().Unmarshal(p.argument(), arg.Binding()))
		delete(p.reqs, pointerFor(arg))
	} else {
		p.extra = append(p.extra, p.argument())
	}
	p.nextArg()
}

// Process a short flag param
func (p *Parser) handleShortFlag() {
	flag, ok := p.shorts[p.char()]

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
		p.waiting = flag
		return
	}

	// Argument is required
	if arg.Required() {
		util.Must(p.com.Unmarshaler().Unmarshal(p.eatString(), arg.Binding()))
		return
	}

	if _, ok := p.shorts[p.char()]; ok {
		p.handleShortFlag()
	} else {
		util.Must(p.com.Unmarshaler().Unmarshal(p.eatString(), arg.Binding()))
	}
}

func (p *Parser) handleLongFlag() {
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

func (p *Parser) isBoolArg(arg A.Argument) bool {
	bt := arg.BindingType().String()

	// Special case for boolean arguments.  The existence of
	// the flag can itself be a true value for a boolean arg.
	return bt == "bool" || bt == "*bool" || bt == "[]bool" || bt == "[]*bool"
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Positioning                                           ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//


// Increments the current character index and returns whether or not we've
// passed the end of the arg string
func (p *Parser) nextChar() bool {
	p.charI++
	return p.charI < p.strLen()
}

// Increments the current argument index, resets the current character index and
// returns whether or not we've passed the end of the arg list
func (p *Parser) nextArg() bool {
	p.argI++
	p.charI = 0
	return p.argI < len(p.input)
}


//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Helpers                                               ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) setup(args []string, com A.Command) {
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

func (p *Parser) handleRequiredArg() {
	// Special case for boolean arguments.  The existence of
	// the flag can itself be a true value for a boolean arg.
	if p.isBoolArg(p.waiting.Argument()) {
		if util.IsBool(p.argument()) {
			p.handleArg()
			return
		} else {
			util.Must(p.com.Unmarshaler().Unmarshal("true", p.popArg().Binding()))
		}
	}

	p.handleArg()
	return
}

func (p *Parser) popArg() (arg A.Argument) {
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

func (p *Parser) eatString() string {
	return p.argument()[p.charI:]
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

// Catch panics and set the panic value to the given error
// pointer if the panic value is an instance of error or
// string.
func recovery(err *error) {
	rec := recover()

	// No panics, nothing to do
	if rec == nil {
		return
	}

	// If the panic was due to an error, pass it up and
	// return.
	if tmp, ok := rec.(error); ok {
		*err = tmp
		//return
	}

	// If the panic was a string, convert it to an error, pass
	// it up and return.
	if tmp, ok := rec.(string); ok {
		*err = errors.New(tmp)
		//return
	}

	// If the panic was some unknown type, it didn't come from
	// us, panic with it again
	panic(rec)
}
