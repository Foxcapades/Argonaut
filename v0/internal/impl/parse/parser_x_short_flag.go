package parse

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
)

// Process a short flag param
func (p *Parser) handleShortFlag() {
	TraceStart("Parser.handleShortFlag")
	defer TraceEnd(func() []interface{} { return nil })

	flag, ok := p.shorts[p.char()]
	Trace(flag, ok)

	// Invalid flag character
	if !ok {
		p.com.AppendUnmapped(string([]byte{'-', p.char()}))
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
				str := p.eatString()
				util.Must(p.com.Unmarshaler().Unmarshal(str, arg.Binding()))
				arg.SetRawValue(str)
			} else {
				Trace("input is not bool")
				util.Must(p.com.Unmarshaler().Unmarshal("true", arg.Binding()))
				arg.SetRawValue("true")
				p.handleShortFlag()
			}
			return
		}
		str := p.eatString()
		util.Must(p.com.Unmarshaler().Unmarshal(str, arg.Binding()))
		arg.SetRawValue(str)
		return
	}

	Trace("short flag arg is optional")

	if _, ok := p.shorts[p.char()]; ok {
		Trace("next char is flag")
		if p.isBoolArg(arg) {
			Trace("cur flag argument is bool")
			util.Must(p.com.Unmarshaler().Unmarshal("true", arg.Binding()))
			arg.SetRawValue("true")
		}
		p.handleShortFlag()
	} else {
		str := p.eatString()
		util.Must(p.com.Unmarshaler().Unmarshal(str, arg.Binding()))
		arg.SetRawValue(str)
	}
}
