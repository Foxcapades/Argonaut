package parse

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	"strings"
)

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
				arg.SetRawValue("true")
				return
			}

			// TODO: make this a real error
			panic("flag " + flag.String() + " requires an argument")
		}
		return
	}

	util.Must(p.com.Unmarshaler().Unmarshal(split[1], arg.Binding()))
	arg.SetRawValue(split[1])
}
