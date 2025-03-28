package tree

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/emit"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Parse(tree argo.TreeCommand, args []string) (argo.ParseResult, error) {
	i := CommandTreeInterpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		current:  tree,
		tree:     tree,
		queue:    utils.NewDeque[parse.Element](2),
		flagHits: flag.NewQueue(),
	}

	res, err := i.Run()

	if err != nil && res.Error == nil {
		res.Error = err
	}

	return res, err
}

type CommandTreeInterpreter struct {
	parser   parse.Parser
	current  any
	boundary bool
	options  argo.Options

	tree     argo.TreeCommand

	// branches keeps track of the branch path that was followed in the CLI call.
	branches []argo.BranchCommand
	leaf     argo.LeafCommand
	queue    utils.Deque[parse.Element]

	flagHits flag.Queue

	result argo.ParseResult
}

func (c *CommandTreeInterpreter) Next() parse.Element {
	if c.queue.IsEmpty() {
		c.queue.Offer(c.parser.Next())
	}

	return c.queue.Poll()
}

func (c *CommandTreeInterpreter) ElementQueue() utils.Deque[parse.Element] {
	return c.queue
}

func (c *CommandTreeInterpreter) HitBoundary() {
	c.boundary = true
}

func (c *CommandTreeInterpreter) CurrentAsFlagGroupContainer() flag.GroupContainer {
	return c.current.(flag.GroupContainer)
}

func (c *CommandTreeInterpreter) FlagHits() *flag.Queue {
	return &c.flagHits
}

func (c *CommandTreeInterpreter) haveLeaf() bool {
	return c.leaf != nil
}

func (c *CommandTreeInterpreter) Run() (argo.ParseResult, error) {
	var argumentStream argument.ValueAppender
	var err error

	unmapped := make([]string, 0, 10)
	errs := xerr.NewMultiError()

FOR:
	for {
		element := c.Next()

		// If we've hit the boundary marker, then everything else becomes either an
		// argument or an unmapped input.
		if c.boundary {
			if element.Type == parse.ElementTypeEnd {
				break
			}

			if ok, err := argumentStream.Append(element.String()); err != nil {
				return c.result, err
			} else if !ok {
				unmapped = append(unmapped, element.String())
			}

			continue
		}

		switch element.Type {
		case parse.ElementTypePlainText:
			if unmapped, err = c.handlePlainText(element, argumentStream, unmapped, errs); err != nil {
				return c.result, err
			}

		case parse.ElementTypeLongFlagPair:
			if err = flag.InterpretLongPair(c, &element, &unmapped, &c.result); err != nil {
				return c.result, err
			}

		case parse.ElementTypeLongFlagSolo:
			if err = flag.InterpretLongSolo(c, &element, &unmapped, &c.result); err != nil {
				return c.result, err
			}

		case parse.ElementTypeShortBlockSolo:
			if err = flag.InterpretShortSolo(c, &element, &unmapped, &c.result); err != nil {
				return c.result, err
			}

		case parse.ElementTypeShortBlockPair:
			if err = flag.InterpretShortPair(c, &element, &unmapped, &c.result); err != nil {
				return c.result, err
			}

		case parse.ElementTypeBoundary:
			c.boundary = true

		case parse.ElementTypeEnd:
			break FOR

		default:
			panic("illegal state: unrecognized parser element type")
		}

		if len(errs.Errors()) > 0 {
			return c.result, errs
		}
	}

	var onIncomplete argo.IncompleteCommandHandler[argo.ParentNode]

	// If the last reached node was a command leaf.
	if c.haveLeaf() {
		// append unrecognized flags
		for _, value := range unmapped {
			c.leaf.AppendUnmappedInput(value)
		}

		argument.CheckRequired(c.leaf.Arguments(), errs)
	} else {
		if parent, ok := c.current.(argo.BranchCommandBuilder); ok {
			onIncomplete = utils.Cast[argo.IncompleteCommandHandler[argo.ParentNode]](parent.IncompleteHandler())
		} else if parent, ok := c.current.(argo.TreeCommandBuilder); ok {
			onIncomplete = utils.Cast[argo.IncompleteCommandHandler[argo.ParentNode]](parent.IncompleteHandler())
		} else {
			errs.AppendError(fmt.Errorf("command leaf was not reached"))
		}
	}

	c.processFlags(c.current.(flag.GroupContainer), errs)

	if onIncomplete != nil {
		onIncomplete(c.current.(argo.ParentNode))
	}

	if len(errs.Errors()) > 0 {
		return c.result, errs
	}

	if c.tree.HasCallback() {
		c.tree.Callback()(c.tree)
	}

	for _, b := range c.branches {
		if b.HasCallback() {
			b.Callback()(b)
		}
	}

	if c.leaf.HasCallback() {
		c.leaf.Callback()(c.leaf)
	}

	return c.result, nil
}
