package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Interpreter struct {
	boundary  bool
	parser    parse.Parser
	flagHits  flag.Queue
	warnings  []argo.InputWarning
	command   argo.Command
	elements  utils.Deque[parse.Element]
	flagIndex flag.Index
}

func (c *Interpreter) FlagHits() *flag.Queue {
	return &c.flagHits
}

func (c *Interpreter) HitBoundary() {
	c.boundary = true
}

func (c *Interpreter) ElementQueue() utils.Deque[parse.Element] {
	return c.elements
}

func (c *Interpreter) FlagIndex() flag.Index {
	return c.flagIndex
}

func (c *Interpreter) Next() parse.Element {
	if c.elements.IsEmpty() {
		c.elements.Offer(c.parser.Next())
	}

	return c.elements.Poll()
}

func (c *Interpreter) Run() ([]argo.InputWarning, error) {
	argumentStream := argument.NewValueAppender(c.command.Arguments())

	var interpretFn func(flag.ContainingInterpreter, *parse.Element, *[]string, *[]argo.InputWarning) error
	var unmapped []string

FOR:
	for {
		element := c.Next()

		// If we've already hit the boundary marker then everything else is to be
		// considered an argument.
		if c.boundary {
			if element.Type == parse.ElementTypeEnd {
				break
			}

			if ok, err := argumentStream.Append(element.String()); err != nil {
				return c.warnings, err
			} else if !ok {
				unmapped = append(unmapped, element.String())
			}

			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if ok, err := argumentStream.Append(element.String()); err != nil {
				return c.warnings, err
			} else if !ok {
				unmapped = append(unmapped, element.String())
			}
			continue

		case parse.ElementTypeShortBlockSolo:
			interpretFn = flag.InterpretShortSolo

		case parse.ElementTypeShortBlockPair:
			interpretFn = flag.InterpretShortPair

		case parse.ElementTypeLongFlagSolo:
			interpretFn = flag.InterpretLongSolo

		case parse.ElementTypeLongFlagPair:
			interpretFn = flag.InterpretLongPair

		case parse.ElementTypeBoundary:
			c.boundary = true
			continue

		case parse.ElementTypeEnd:
			break FOR

		default:
			panic("illegal state")
		}

		if err := interpretFn(c, &element, &unmapped, &c.warnings); err != nil {
			return c.warnings, err
		}
	}

	for _, u := range unmapped {
		c.command.AppendUnmappedInput(u)
	}

	errs := xerr.NewMultiError()

	flag.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())
	flag.ExecuteFlagCallbacks(c.flagHits.Iterator())
	flag.CheckRequired(c.command.FlagGroups(), errs)
	argument.CheckRequired(c.command.Arguments(), errs)

	if len(errs.Errors()) > 0 {
		return c.warnings, errs
	}

	if c.command.HasCallback() {
		c.command.Callback()(c.command)
	}

	return c.warnings, nil
}
