package parse

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	R "reflect"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internal API: Helpers                                               ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

func (p *Parser) setup(args []string, com A.Command) {
	TraceStart("Parser.setup", args, com)
	defer TraceEnd(func() []interface{} { return nil })
	p.makeMaps(com)
	p.args = com.Arguments()
	p.input = args
	p.com = com
	p.argI = 0
	p.charI = 0
	p.waiting = nil
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
		arg.SetRawValue(p.argument())
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
			arg.SetRawValue(p.argument())
		} else {
			util.Must(p.com.Unmarshaler().Unmarshal("true", bind))
			arg.SetRawValue("true")
			p.argI--
		}
		return
	}

	util.Must(p.com.Unmarshaler().Unmarshal(p.argument(), bind))
	arg.SetRawValue(p.argument())
}
