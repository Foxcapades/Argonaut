package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Interpreter struct {
	boundary bool
	parser   parse.Parser
	flagHits flag.Queue
	result   argo.ParseResult
	command  argo.Command
	elements utils.Deque[parse.Element]
}

func (c *Interpreter) CurrentAsFlagGroupContainer() flag.GroupContainer {
	return c.command
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

func (c *Interpreter) Next() parse.Element {
	if c.elements.IsEmpty() {
		c.elements.Offer(c.parser.Next())
	}

	return c.elements.Poll()
}

func (c *Interpreter) Run() (argo.ParseResult, error) {
	argumentStream := argument.NewValueAppender(c.command.Arguments())

	var interpretFn func(c flag.ContainingInterpreter, element *parse.Element, unmapped *[]string, result *argo.ParseResult) error
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
				return xerr.AppendError(c.result, err)
			} else if !ok {
				unmapped = append(unmapped, element.String())
			}

			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if ok, err := argumentStream.Append(element.String()); err != nil {
				return xerr.AppendError(c.result, err)
			} else if !ok {
				unmapped = append(unmapped, element.String())
				continue
			}

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

		if err := interpretFn(c, &element, &unmapped, &c.result); err != nil {
			return xerr.AppendError(c.result, err)
		}
	}

	for _, u := range unmapped {
		c.command.AppendUnmappedInput(u)
	}

	errs := xerr.NewMultiError()

	common.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())
	common.CheckRequiredFlags(c.command.FlagGroups(), errs)
	common.CheckRequiredArguments(c.command.Arguments(), errs)

	if len(errs.Errors()) > 0 {
		return xerr.AppendError(c.result, errs)
	}

	c.command.Callback()(c.command)

	return c.result, nil
}
