package tree

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal_old/emit"
	"github.com/foxcapades/argonaut/v3/internal_old/parse"
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
	current  argo.Node
	boundary bool

	tree     argo.CommandTree
	branches []argo.Branch
	leaf     argo.CommandLeaf
	queue    utils.Deque[parse.Element]

	flagHits flag.Queue
}

func (c *commandTreeInterpreter) next() parse.Element {
	if c.queue.IsEmpty() {
		c.queue.Offer(c.parser.Next())
	}

	return c.queue.Poll()
}

type leafCatch struct {
	leaf argo.CommandLeaf
}

func (c *commandTreeInterpreter) Run() error {
	unmapped := make([]string, 0, 10)

FOR:
	for {
		element := c.next()

		// If we've hit the boundary marker, then everything else becomes a
		// passthrough element.
		if c.boundary {

			if element.Type == parse.ElementTypeEnd {
				break
			}

			passthroughs = append(passthroughs, element.String())
			continue
		}

		switch element.Type {
		case parse.ElementTypePlainText:
			// If we've hit the leaf node, then the plain text becomes an argument on
			// that node.  If we haven't yet hit the leaf node, then we must treat the
			// plaintext value as the name of the next node in the tree.  If no such
			// node exists, that is an error.
			if node, ok := c.current.(argo.CommandLeaf); ok {
				if err := appendArgument(node, element.String()); err != nil {
					return err
				}
			} else if node, ok := c.current.(argo.ParentNode); ok {
				// Lookup a child with the given input string
				if child := node.FindChild(element.String()); child != nil {
					c.current = child

					if branch, ok := child.(argo.Branch); ok {
						c.branches = append(c.branches, branch)
					} else if leaf, ok := child.(argo.CommandLeaf); ok {
						c.leaf = leaf
					}
				} else {
					// If node child could be found matching the input string, then print
					// out a help message about the invalid subcommand.
					return c.invalidSubCommand(element.String())
				}
			} else {
				panic("illegal state: command node was neither a leaf or a parent")
			}

		case parse.ElementTypeLongFlagPair:
			if err := c.interpretLongPair(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagSolo:
			if err := c.interpretLongSolo(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockSolo:
			if err := c.interpretShortSolo(&element, &unmapped); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockPair:
			if err := c.interpretShortPair(&element, &unmapped); err != nil {
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

	errs := argo.NewMultiError()
	var onIncomplete func(parent argo.ParentNode)

	// If the last reached node was NOT a command leaf.
	if node, ok := c.current.(argo.CommandLeaf); !ok {
		if parent, ok := c.current.(argo.ParentNode); ok {
			onIncomplete = parent.IncompleteHandler()
		} else {
			errs.AppendError(fmt.Errorf("command leaf was not reached"))
		}
	} else {
		for _, value := range unmapped {
			node.AppendUnmappedInput(value)
		}

		for _, value := range passthroughs {
			node.appendPassthrough(value)
		}

		c.checkRequiredArgsWereHit(node.Arguments(), errs)
	}

	c.checkRequiredFlagsWereHit(c.current, errs)

	// Hunt the help flag down first and execute its callback.  It is likely to
	// exit the application so we don't want to run any other callbacks first.
	helpFlagIndex := 0
	for flag := range c.flagHits.Iterator() {
		if flag.IsHelpFlag() {
			flag.Callback()(flag)
			break
		}

		helpFlagIndex++
	}

	currentFlagIndex := 0
	for flag := range c.flagHits.Iterator() {
		if currentFlagIndex != helpFlagIndex {
			flag.Callback()(flag)
		}

		currentFlagIndex++
	}

	if onIncomplete != nil {
		onIncomplete(c.current.(argo.ParentNode))
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	if c.tree.HasCallback() {
		c.tree.Callback()(c.tree)
	}

	for _, b := range c.branches {
		if b.HasCallback() {
			b.executeCallback()
		}
	}

	if c.leaf.hasCallback() {
		c.leaf.executeCallback()
	}

	c.tree.selectCommand(c.leaf)

	return nil
}

func (c *commandTreeInterpreter) checkRequiredArgsWereHit(args []cli_arg.Argument, errs cli_err.MultiError) {
	for i, arg := range args {
		if arg.IsRequired() {
			if !arg.WasHit() {
				if arg.HasName() {
					errs.AppendError(fmt.Errorf("argument %d (<%s>) is required", i+1, arg.Name()))
				} else {
					errs.AppendError(fmt.Errorf("argument %d is required", i+1))
				}
			}
		} else if !arg.WasHit() && arg.HasDefault() {
			if err := arg.setToDefault(); err != nil {
				errs.AppendError(err)
			}
		}
	}
}

func (c *commandTreeInterpreter) checkRequiredFlagsWereHit(current argo.Node, errs MultiError) {
	for current != nil {
		for _, group := range current.FlagGroups() {
			for _, f := range group.Flags() {
				if f.IsRequired() {
					if !f.WasHit() {
						errs.AppendError(fmt.Errorf("required flag %s was not used", printFlagNames(f)))
					} else if f.RequiresArgument() && !f.Argument().WasHit() {
						errs.AppendError(fmt.Errorf("flag %s requires an argument", printFlagNames(f)))
					}
				} else if !f.WasHit() && f.HasArgument() && f.Argument().HasDefault() {
					if err := f.Argument().setToDefault(); err != nil {
						errs.AppendError(err)
					}
				}
			}
		}

		current = current.Parent()
	}
}

func (c *commandTreeInterpreter) interpretShortSolo(element *parse.Element, unmapped *[]string) error {
	remainder := element.Data[0]

	for i := 0; i < len(element.Data[0]); i++ {
		// has next
		h := i+1 < len(element.Data[0])
		// short flag byte
		b := remainder[0]

		// Look up the flag in the short flag map
		f := c.current.FindShortFlag(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			c.tree.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			*unmapped = append(*unmapped, chars.StrDash+remainder[0:1])
			remainder = remainder[1:]
			continue
		}

		c.flagHits.append(f)

		// If the flag we found requires an argument, eat the rest of the block and
		// pass it to the flag.Hit method.  Since the block will have been consumed
		// after this, return here.
		if f.RequiresArgument() {

			// If we don't have any more characters in this short block, then we have
			// to consume the next element as the argument for this flag.
			if !h {
				nextElement := c.next()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					if hasBooleanArgument(f) {
						return f.hitWithArg("true")
					}
					return f.hit()
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					c.boundary = true
					if hasBooleanArgument(f) {
						return f.hitWithArg("true")
					}
					return f.hit()
				}

				// If we're here then we have a next element, and we're going to try and
				// sacrifice it to the flag gods.
				return f.hitWithArg(nextElement.String())
			}

			if hasBooleanArgument(f) {
				possibleNextFlag := c.current.FindShortFlag(remainder[1])
				if possibleNextFlag != nil {
					if err := f.hitWithArg("true"); err != nil {
						return err
					}
					remainder = remainder[1:]
					continue
				}
			}

			// So we have at least one more character in this block.  Eat that and
			// anything else as the flag argument.
			return f.hitWithArg(remainder[1:])
		}

		// If the flag doesn't _require_ an argument, but may take an optional
		// one...
		if f.HasArgument() {

			// If we have a next character in the block...
			if h {

				// grab the next character
				n := remainder[1]

				// test if the next character is a flag itself.  If it is, then we
				// prioritize the flag over an optional argument.
				if t := c.current.FindShortFlag(n); t != nil {

					// hit the current flag with an empty value
					if err := f.hit(); err != nil {
						return err
					}

					// skip on to the next flag
					remainder = remainder[1:]
					continue
				} else

				// Since there is no flag matching the next character, then we have to
				// assume that the remaining text is the argument for the flag
				{
					return f.hitWithArg(remainder[1:])
				}
			} else {

				nextElement := c.next()

				switch nextElement.Type {

				case parse.ElementTypeEnd:
					if hasBooleanArgument(f) {
						return f.hitWithArg("true")
					}
					return f.hit()

				case parse.ElementTypeBoundary:
					if hasBooleanArgument(f) {
						return f.hitWithArg("true")
					}
					c.boundary = true
					return f.hit()

				case parse.ElementTypePlainText:
					if err := f.hitWithArg(nextElement.String()); err != nil {
						c.queue.Offer(nextElement)
					}

					return f.hit()

				case parse.ElementTypeShortBlockSolo:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.queue.Offer(nextElement)
						return f.hit()
					} else {
						if err := f.hitWithArg(nextElement.String()); err != nil {
							c.queue.Offer(nextElement)
						}

						return f.hit()
					}

				case parse.ElementTypeShortBlockPair:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.queue.Offer(nextElement)
						return f.hit()
					} else {
						if err := f.hitWithArg(nextElement.String()); err != nil {
							c.queue.Offer(nextElement)
						}

						return f.hit()
					}

				case parse.ElementTypeLongFlagPair:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						c.queue.Offer(nextElement)
						return f.hit()
					} else {
						if err := f.hitWithArg(nextElement.String()); err != nil {
							c.queue.Offer(nextElement)
						}

						return f.hit()
					}

				case parse.ElementTypeLongFlagSolo:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						c.queue.Offer(nextElement)
						return f.hit()
					} else {
						if err := f.hitWithArg(nextElement.String()); err != nil {
							c.queue.Offer(nextElement)
						}

						return f.hit()
					}

				default:
					panic("illegal state: unrecognized parser element type")

				}
			}
		}

		// The flag doesn't expect an argument, just hit it
		if err := f.hit(); err != nil {
			return err
		}

		// if has next
		remainder = remainder[1:]
	}

	return nil
}

// interpretShortPair tries to make sense of a pair where the first value is a
// block of one or more short flags, and the second value is an argument value
// that was directly attached using an `=` character.
func (c *commandTreeInterpreter) interpretShortPair(element *parse.Element, unmapped *[]string) error {
	block := element.Data[0]

	if len(block) == 0 {
		c.tree.AppendWarning("blank short flag name")
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.current.FindShortFlag(block[0]); f != nil {
			c.flagHits.append(f)
			return f.hitWithArg(element.Data[1])
		} else {
			*unmapped = append(*unmapped, element.String())
			return nil
		}
	}

	for i := 0; i < len(element.Data[0]); i++ {
		// has next character
		h := i+1 < len(element.Data[0])
		// current character
		b := block[0]

		f := c.current.FindShortFlag(b)

		if f == nil {
			c.tree.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			*unmapped = append(*unmapped, chars.StrDash+block[0:1])
			block = block[1:]
			continue
		}

		c.flagHits.append(f)

		if f.RequiresArgument() {
			if h {
				return f.hitWithArg(block[1:] + "=" + element.Data[1])
			} else {
				return f.hitWithArg(element.Data[1])
			}
		}

		// If the current flag has, but does not require an argument...
		if f.HasArgument() {
			// and there is no next character in the flag name block...
			if !h {
				// Hit the current flag with the argument value and exit.
				return f.hitWithArg(element.Data[1])
			}

			// If there _is_ a next character, and it happens to be a valid short
			// flag itself, then hit the current flag and move on to the next
			// character in the block.
			if c.current.FindShortFlag(block[1]) != nil {
				if err := f.hit(); err != nil {
					return err
				}
				block = block[1:]
				continue
			}

			// If the next character in the block does not match any known short flag,
			// assume that the whole remaining value is part of the value.
			return f.hitWithArg(block[1:] + "=" + element.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !h {
			c.tree.AppendWarning(fmt.Sprintf("flag -%c received an argument it didn't expect", b))
			return f.hitWithArg(element.Data[1])
		}

		// Well, now that's out of the way, we can move on to the next flag (after
		// we mark this one as hit of course).
		if err := f.hit(); err != nil {
			return err
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c *commandTreeInterpreter) interpretLongSolo(element *parse.Element, unmapped *[]string) error {
	f := c.current.FindLongFlag(element.Data[0])

	if f == nil {
		c.tree.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", element.Data[0]))
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.flagHits.append(f)

	if f.RequiresArgument() {
		nextElement := c.next()

		if nextElement.Type == parse.ElementTypeEnd {
			return f.hit()
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			c.boundary = true
			return f.hit()
		}

		return f.hitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.next()

		switch nextElement.Type {

		case parse.ElementTypeEnd:
			return f.hit()

		case parse.ElementTypeBoundary:
			c.boundary = true
			return f.hit()

		case parse.ElementTypePlainText:
			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.queue.Offer(nextElement)
			}

			return f.hit()

		case parse.ElementTypeLongFlagSolo:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				c.queue.Offer(nextElement)
				return f.hit()
			} else {
				if err := f.hitWithArg(nextElement.String()); err != nil {
					c.queue.Offer(nextElement)
				}

				return f.hit()
			}

		case parse.ElementTypeLongFlagPair:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				c.queue.Offer(nextElement)
				return f.hit()
			} else {
				if err := f.hitWithArg(nextElement.String()); err != nil {
					c.queue.Offer(nextElement)
				}

				return f.hit()
			}

		case parse.ElementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.queue.Offer(nextElement)
				return f.hit()
			} else {
				if err := f.hitWithArg(nextElement.String()); err != nil {
					c.queue.Offer(nextElement)
				}

				return f.hit()
			}

		case parse.ElementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.queue.Offer(nextElement)
				return f.hit()
			} else {
				if err := f.hitWithArg(nextElement.String()); err != nil {
					c.queue.Offer(nextElement)
				}

				return f.hit()
			}

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	return f.hit()
}

func (c *commandTreeInterpreter) interpretLongPair(element *parse.Element, unmapped *[]string) error {
	flag := c.current.FindLongFlag(element.Data[0])

	if flag == nil {
		c.tree.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", element.Data[0]))
		*unmapped = append(*unmapped, element.String())
	} else {
		c.flagHits.append(flag)

		if flag.HasArgument() {
			return flag.hitWithArg(element.Data[1])
		} else {
			c.tree.AppendWarning(fmt.Sprintf("flag --%s received an argument it didn't expect", element.Data[0]))
			return flag.hit()
		}
	}

	return nil
}

func (c *commandTreeInterpreter) invalidSubCommand(input string) error {
	type pair struct {
		depth int
		child string
	}

	matches := make([]pair, 0, 8)

	if parent, ok := c.current.(argo.ParentNode); ok {
		for _, group := range parent.CommandGroups() {
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

func appendArgument(leaf argo.CommandLeaf, value string) error {
	for _, arg := range leaf.Arguments() {
		if !arg.WasHit() {
			return arg.SetValue(value)
		}
	}

	leaf.AppendUnmappedInput(value)

	return nil
}
