package impl

import (
	"fmt"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

type parserState uint8

const (
	psUnknown parserState = iota
	psExpectRequired
	psExpectOptional
)

func NewParser() argo.Parser {
	return &parser{}
}

type parser struct {
	shorts map[byte]argo.Flag
	longs  map[string]argo.Flag

	// Input arguments
	args []string

	// Current input string position
	currentInputArg int

	// Current argument in command
	currentComArg int

	// Current character index in args[currentInputArg]
	curArgPos int

	// Current character value
	curChar byte

	// in passthrough mode
	passthrough bool

	// Command to parse
	com argo.Command

	// last error
	err error

	awaitingValue argo.Flag
}

func (p *parser) Parse(args []string) error {
	//p.args = filterArgs(args)
	//
	//if p.atEoa() {
	//	return p.wrapup()
	//}
	//
	//p.shorts = make(map[byte]argo.Flag)
	//p.longs = make(map[string]argo.Flag)
	//
	//for _, arg := range p.com.Flags() {
	//	if arg.HasShort() {
	//		if p.hasShort(arg.Short()) {
	//			// TODO: Make this an official error
	//			return fmt.Errorf("duplicate entries for short flag %c", arg.Short())
	//		}
	//		p.shorts[arg.Short()] = arg
	//	}
	//
	//	if arg.HasLong() {
	//		if p.hasLong(arg.Long()) {
	//			// TODO: make this an official error
	//			return fmt.Errorf("duplicate entries for long flag %s", arg.Long())
	//		}
	//		p.longs[arg.Long()] = arg
	//	}
	//}
	//
	//p.curChar = p.args[p.curChar][p.curArgPos]
	//return p.run()
	return nil
}

func (p *parser) hasShort(f byte) bool {
	_, o := p.shorts[f]
	return o
}
func (p *parser) hasLong(f string) bool {
	_, o := p.longs[f]
	return o
}

// At end of input
func (p *parser) atEoi() bool {
	return p.currentInputArg >= len(p.args)
}

// At end of arg string
func (p *parser) atEoa() bool {
	return p.curArgPos >= len(p.args[p.currentInputArg])
}

func (p *parser) hasNextChar() bool {
	return p.curArgPos+1 < len(p.args[p.currentInputArg])
}

func (p *parser) hasNextArg() bool {
	return p.currentInputArg+1 < len(p.args)
}

func (p *parser) nextArg() {
	if !p.hasNextArg() {
		panic("invalid parser state")
	}
	p.currentInputArg++
	p.curArgPos = 0
	p.curChar = p.args[p.currentInputArg][p.curArgPos]
}

func (p *parser) nextChar() {
	if p.hasNextChar() {
		p.curArgPos++
		p.curChar = p.args[p.currentInputArg][p.curArgPos]
	}

	panic("invalid parser state")
}

func (p *parser) peek() byte {
	return p.args[p.currentInputArg][p.curArgPos+1]
}

func (p *parser) handleArg() {}

func (p *parser) handleFlag() {
	// One character string at end of arg list
	if !p.hasNextChar() {
		if p.currentComArg >= len(p.com.Arguments()) {
			//p.com.appendUnmappedInput(p.read(1))
		} else {
			//p.err = p.com.Arguments()[p.currentComArg].parse(p.read(1))
		}
		return
	}

	if p.peek() == '-' {
		p.nextChar()
		p.longFlagPassthroughOrArg()
	} else {
		p.nextChar()
		p.shortFlag()
	}
}

func (p *parser) shortFlag() {
	if !p.hasShort(p.curChar) {
		// TODO: make this an official error
		p.err = fmt.Errorf("unrecognized flag %c", p.curChar)
		return
	}

	//flag := p.shorts[p.curChar]
	//flag.hit()

	//if !flag.HasArgument() {
	//	if p.hasNextChar() {
	//		p.nextChar()
	//		p.shortFlag()
	//	}
	//	return
	//}

	//arg := flag.Argument()

	//if arg.Required() {
	//
	//}
}

func (p *parser) longFlagPassthroughOrArg() {}

func (p *parser) read(n int) (val string) {
	val = p.args[p.currentInputArg][p.curArgPos : p.curArgPos+n]
	p.curArgPos += n
	return
}

func (p *parser) run() error {
	for true {
		if p.err != nil {
			return p.err
		}

		if p.curChar == '-' {
			p.handleFlag()
		} else {
			p.handleArg()
		}

		if !p.hasNextChar() {
			break
		}
		p.nextChar()
	}
	return nil
}

// confirm we aren't missing a required arg somewhere
func (p *parser) wrapup() error {
	return nil
}

func filterArgs(args []string) []string {
	out := make([]string, 0, len(args))

	for _, a := range args {
		if len(a) > 0 {
			out = append(out, a)
		}
	}

	return out
}
