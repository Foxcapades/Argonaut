package tree

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/emit"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newCommandTreeInterpreter(args []string, command argo.CommandTree) commandTreeInterpreter {
	return commandTreeInterpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		current:  command,
		tree:     command,
		queue:    utils.NewDeque[parse.Element](2),
		flagHits: flag.NewQueue(),
	}
}

type commandTreeInterpreter struct {
	parser   parse.Parser
	current  argo.Node[any]
	boundary bool

	tree     argo.CommandTree
	branches []argo.BranchCommand
	leaf     argo.LeafCommand
	queue    utils.Deque[parse.Element]

	flagHits flag.Queue
}

func (c *commandTreeInterpreter) next() parse.Element {
	if c.queue.IsEmpty() {
		c.queue.Offer(c.parser.Next())
	}

	return c.queue.Poll()
}

func (c *commandTreeInterpreter) haveLeaf() bool {
	return c.leaf != nil
}

func (c *commandTreeInterpreter) Run() error {
	var argumentStream argument.Appender
	var err error

	unmapped := make([]string, 0, 10)

FOR:
	for {
		element := c.next()

		// If we've hit the boundary marker, then everything else becomes either an
		// argument or an unmapped input.
		if c.boundary {
			if element.Type == parse.ElementTypeEnd {
				break
			}

			if ok, err := argumentStream.Append(element.String()); err != nil {
				return err
			} else if !ok {
				unmapped = append(unmapped, element.String())
			}

			continue
		}

		switch element.Type {
		case parse.ElementTypePlainText:
			if unmapped, err = c.handlePlainText(element, argumentStream, unmapped); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagPair:
			if err = c.interpretLongPair(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagSolo:
			if err = c.interpretLongSolo(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockSolo:
			if err = c.interpretShortSolo(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockPair:
			if err = c.interpretShortPair(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeBoundary:
			c.boundary = true

		case parse.ElementTypeEnd:
			break FOR

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	errs := xerr.NewMultiError()
	var onIncomplete argo.IncompleteCommandHandler[any]

	// If the last reached node was a command leaf.
	if c.haveLeaf() {
		// append unrecognized flags
		for _, value := range unmapped {
			c.leaf.AppendUnmappedInput(value)
		}

		common.CheckRequiredArguments(c.leaf.Arguments, errs)
	} else {
		if parent, ok := c.current.(argo.ParentNode[any]); ok {
			onIncomplete = parent.IncompleteHandler()
		} else {
			errs.AppendError(fmt.Errorf("command leaf was not reached"))
		}
	}

	common.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())

	c.checkRequiredFlagsWereHit(errs)

	common.ExecuteFlagCallbacks(c.flagHits.Iterator())

	if onIncomplete != nil {
		onIncomplete(c.current.(argo.ParentNode[any]))
	}

	if len(errs.Errors()) > 0 {
		return errs
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

	return nil
}

func (c *commandTreeInterpreter) checkRequiredFlagsWereHit(errs argo.MultiError) {
	var current argo.Node[any] = c.leaf

	for {
		common.CheckRequiredFlags(current.FlagGroups, errs)

		if par, ok := current.(argo.ChildNode[any]); ok {
			current = par
		} else {
			break
		}
	}
}

func (c *commandTreeInterpreter) handlePlainText(element parse.Element, arguments argument.Appender, unmapped []string) ([]string, error) {
	// If we've hit the leaf node, then the plain text becomes an argument on
	// that node.  If we haven't yet hit the leaf node, then we must treat the
	// plaintext value as the name of the next node in the tree.  If no such
	// node exists, that is an error.
	if _, ok := c.current.(argo.LeafCommand); ok {
		if ok, err := arguments.Append(element.String()); err != nil {
			return unmapped, err
		} else if !ok {
			return append(unmapped, element.String()), nil
		}

		// argument value was accepted
		return unmapped, nil
	}

	if node, ok := c.current.(argo.ParentNode[any]); ok {
		// Lookup a child with the given input string
		if child := node.FindChild(element.String()); child != nil {
			c.current = child
			node.SelectChild(element.String())

			if branch, ok := child.(argo.BranchCommand); ok {
				c.branches = append(c.branches, branch)
			} else if leaf, ok := child.(argo.LeafCommand); ok {
				c.leaf = leaf
			}
		}

		// If node child could be found matching the input string, then print
		// out a help message about the invalid subcommand.
		return unmapped, c.invalidSubCommand(element.String())
	}

	panic("illegal state: command node was neither a leaf or a parent")
}

func (c *commandTreeInterpreter) interpretShortSolo(element *parse.Element, unmapped *[]string) error {
	remainder := element.Data[0]

	for i := 0; i < len(element.Data[0]); i++ {
		// has next
		h := i+1 < len(element.Data[0])
		// short flag byte
		b := remainder[0]
		remainder = remainder[1:]

		// Look up the flag in the short flag map
		f := c.current.FindShortFlag(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			// c.tree.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			*unmapped = append(*unmapped, chars.StrDash+string(b))
			continue
		}

		f.IncrementHitCount()
		c.flagHits.Append(f)

		if !f.HasArgument() {
			continue
		}

		arg := f.Argument()

		// If the flag we found requires an argument, eat the rest of the block and
		// pass it to the flag.Hit method.  Since the block will have been consumed
		// after this, return here.
		if arg.IsRequired() {

			// If we don't have any more characters in this short block, then we have
			// to consume the next element as the argument for this flag.
			if !h {
				nextElement := c.next()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					if f.HasArgument() && argument.IsBoolean(arg) {
						return arg.SetValue("true")
					}

					return nil
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					c.boundary = true

					if f.HasArgument() && argument.IsBoolean(arg) {
						return arg.SetValue("true")
					}

					return nil
				}

				// If we're here then we have a next element, and we're going to try and
				// sacrifice it to the flag gods.
				return arg.SetValue(nextElement.String())
			}

			if argument.IsBoolean(arg) {
				possibleNextFlag := c.current.FindShortFlag(remainder[0])

				if possibleNextFlag != nil {
					if err := arg.SetValue("true"); err != nil {
						return err
					}

					continue
				}
			}

			// So we have at least one more character in this block.  Eat that and
			// anything else as the flag argument.
			return arg.SetValue(remainder)
		}

		// If we have a next character in the block...
		if h {

			// grab the next character
			n := remainder[0]

			// test if the next character is a flag itself.  If it is, then we
			// prioritize the flag over an optional argument.
			if t := c.current.FindShortFlag(n); t != nil {
				if argument.IsBoolean(arg) {
					if err := arg.SetValue("true"); err != nil {
						return err
					}
				}

				// skip on to the next flag
				continue
			}

			// Since there is no flag matching the next character, then we have to
			// assume that the remaining text is the argument for the flag
			return arg.SetValue(remainder)
		}

		nextElement := c.next()

		switch nextElement.Type {

		case parse.ElementTypeEnd:
			if argument.IsBoolean(arg) {
				return arg.SetValue("true")
			}

		case parse.ElementTypeBoundary:
			c.boundary = true
			if argument.IsBoolean(arg) {
				return arg.SetValue("true")
			}

		case parse.ElementTypePlainText:
			// Try and give the argument the next value, if the argument doesn't like
			// it then put it back on the queue for the next parse iteration.
			if err := arg.SetValue(nextElement.String()); err != nil {
				c.queue.Offer(nextElement)
			}

		case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
			if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.queue.Offer(nextElement)
			} else if err := arg.SetValue(nextElement.String()); err != nil {
				c.queue.Offer(nextElement)
			}

		case parse.ElementTypeLongFlagPair, parse.ElementTypeLongFlagSolo:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				c.queue.Offer(nextElement)
			} else if err := arg.SetValue(nextElement.String()); err != nil {
				c.queue.Offer(nextElement)
			}

		default:
			panic("illegal state: unrecognized parser element type")
		}

		break
	}

	return nil
}

// interpretShortPair tries to make sense of a pair where the first value is a
// block of one or more short flags, and the second value is an argument value
// that was directly attached using an `=` character.
func (c *commandTreeInterpreter) interpretShortPair(element *parse.Element, unmapped *[]string) error {
	block := element.Data[0]

	if len(block) == 0 {
		// c.tree.AppendWarning("blank short flag name")
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.current.FindShortFlag(block[0]); f != nil {
			c.flagHits.Append(f)
			f.IncrementHitCount()
			if f.HasArgument() {
				return f.Argument().SetValue(element.Data[1])
			}

			// TODO: warn or error for flag given an argument when it doesn't expect
			//       one.
			return nil
		} else {
			*unmapped = append(*unmapped, element.String())
			return nil
		}
	}

	for i := 0; i < len(element.Data[0]); i++ {
		// has next character
		hasNextChar := i+1 < len(element.Data[0])
		// current character
		b := block[0]

		f := c.current.FindShortFlag(b)

		if f == nil {
			// c.tree.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			*unmapped = append(*unmapped, chars.StrDash+block[0:1])
			block = block[1:]
			continue
		}

		c.flagHits.Append(f)
		f.IncrementHitCount()

		if f.HasArgument() && f.Argument().IsRequired() {
			if hasNextChar {
				return f.Argument().SetValue(block[1:] + "=" + element.Data[1])
			} else {
				return f.Argument().SetValue(element.Data[1])
			}
		}

		// If the current flag has, but does not require an argument...
		if f.HasArgument() {
			// and there is no next character in the flag name block...
			if !hasNextChar {
				// Hit the current flag with the argument value and exit.
				return f.Argument().SetValue(element.Data[1])
			}

			// If there _is_ a next character, and it happens to be a valid short
			// flag itself, then hit the current flag and move on to the next
			// character in the block.
			if c.current.FindShortFlag(block[1]) != nil {
				block = block[1:]
				continue
			}

			// If the next character in the block does not match any known short flag,
			// assume that the whole remaining value is part of the value.
			return f.Argument().SetValue(block[1:] + "=" + element.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !hasNextChar {
			// TODO: c.tree.AppendWarning(fmt.Sprintf("flag -%c received an argument it didn't expect", b))
			return f.Argument().SetValue(element.Data[1])
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c *commandTreeInterpreter) interpretLongSolo(element *parse.Element, unmapped *[]string) error {
	f := c.current.FindLongFlag(element.Data[0])

	if f == nil {
		// TODO: c.tree.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", element.Data[0]))
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.flagHits.Append(f)
	f.IncrementHitCount()

	if !f.HasArgument() {
		return nil
	}

	arg := f.Argument()

	if arg.IsRequired() {
		nextElement := c.next()

		if nextElement.Type == parse.ElementTypeEnd {
			return nil
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			c.boundary = true
			return nil
		}

		return arg.SetValue(nextElement.String())
	}

	nextElement := c.next()

	switch nextElement.Type {

	case parse.ElementTypeEnd:
		// do nothing

	case parse.ElementTypeBoundary:
		c.boundary = true

	case parse.ElementTypePlainText:
		if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	case parse.ElementTypeLongFlagSolo, parse.ElementTypeLongFlagPair:
		if c.current.FindLongFlag(nextElement.Data[0]) != nil {
			c.queue.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
		if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
			c.queue.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	default:
		panic("illegal state: unrecognized parser element type")
	}

	return nil
}

func (c *commandTreeInterpreter) interpretLongPair(element *parse.Element, unmapped *[]string) error {
	targetFlag := c.current.FindLongFlag(element.Data[0])

	if targetFlag == nil {
		// TODO: c.tree.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", element.Data[0]))
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.flagHits.Append(targetFlag)
	targetFlag.IncrementHitCount()

	if targetFlag.HasArgument() {
		return targetFlag.Argument().SetValue(element.Data[1])
	} else {
		// TODO: c.tree.AppendWarning(fmt.Sprintf("flag --%s received an argument it didn't expect", element.Data[0]))
	}

	return nil
}

func (c *commandTreeInterpreter) invalidSubCommand(input string) error {
	type pair struct {
		depth int
		child string
	}

	matches := make([]pair, 0, 8)

	if parent, ok := c.current.(argo.ParentNode[any]); ok {
		for _, group := range parent.CommandGroups(true) {
			for _, child := range group.Branches() {
				if idx := strings.Index(child.Name(), input); idx > -1 {
					matches = append(matches, pair{idx, child.Name()})
				} else {
					for _, alias := range child.Aliases() {
						if idx := strings.Index(alias, input); idx > -1 {
							matches = append(matches, pair{idx, alias})
						}
					}
				}
			}

			for _, child := range group.Leaves() {
				if idx := strings.Index(child.Name(), input); idx > -1 {
					matches = append(matches, pair{idx, child.Name()})
				} else {
					for _, alias := range child.Aliases() {
						if idx := strings.Index(alias, input); idx > -1 {
							matches = append(matches, pair{idx, alias})
						}
					}
				}
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].depth < matches[j].depth || matches[i].child < matches[j].child
	})

	msg := new(bytes.Buffer)
	buf := bufio.NewWriter(msg)

	// TODO: Instead of using `tree.Name()` here print out the full path to the
	//       current subcommand.  So like `app foo bar: Subcommand "f" is blah..."
	utils.MustReturn(buf.WriteString(fmt.Sprintf("%s: Subcommand \"%s\" is unrecognized.", c.tree.Name(), input)))
	t := 0
	if c.current.FindShortFlag('h') != nil {
		t = 1
	} else if c.current.FindLongFlag("help") != nil {
		t = 2
	}

	if t != 0 {
		utils.MustReturn(buf.WriteString(" See available subcommands by using "))
		if t == 1 {
			utils.MustReturn(buf.WriteString("-h"))
		} else {
			utils.MustReturn(buf.WriteString("--help"))
		}
	}

	if len(matches) > 0 {
		if len(matches) == 1 {
			utils.MustReturn(buf.WriteString("\n\nPerhaps you meant:\n"))
		} else {
			utils.MustReturn(buf.WriteString("\n\nPerhaps you meant one of:\n"))
		}

		for i := range matches {
			utils.MustReturn(buf.WriteString("    "))
			utils.MustReturn(buf.WriteString(matches[i].child))
			utils.MustReturn(buf.WriteString("\n"))
		}
	} else {
		utils.MustReturn(buf.WriteString("\n"))
	}

	utils.Must(buf.Flush())

	//goland:noinspection GoUnreachableCode
	return fmt.Errorf(msg.String())
}
