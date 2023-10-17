package argo

import (
	"fmt"
)

type commandTreeInterpreter struct {
	parser   parser
	current  CommandNode
	boundary bool

	tree     CommandTree
	branches []CommandBranch
	leaf     CommandLeaf
	queue    deque[element]

	flagHits []Flag
}

func (c *commandTreeInterpreter) next() element {
	if c.queue.IsEmpty() {
		c.queue.Offer(c.parser.Next())
	}

	return c.queue.Poll()
}

func (c *commandTreeInterpreter) Run() error {
	passthroughs := make([]string, 0, 10)
	unmapped := make([]string, 0, 10)

FOR:
	for {
		element := c.next()

		// If we've hit the boundary marker, then everything else becomes a
		// passthrough element.
		if c.boundary {
			if element.Type == elementTypeEnd {
				break
			}

			passthroughs = append(passthroughs, element.String())
			continue
		}

		switch element.Type {
		case elementTypePlainText:
			// If we've hit the leaf node, then the plain text becomes an argument on
			// that node.  If we haven't yet hit the leaf node, then we must treat the
			// plaintext value as the name of the next node in the tree.  If no such
			// node exists, that is an error.
			if node, ok := c.current.(CommandLeaf); ok {
				if err := node.appendArgument(element.String()); err != nil {
					return err
				}
			} else if node, ok := c.current.(CommandParent); ok {
				if child := node.FindChild(element.String()); child != nil {
					c.current = child

					if branch, ok := child.(CommandBranch); ok {
						c.branches = append(c.branches, branch)
					} else if leaf, ok := child.(CommandLeaf); ok {
						c.leaf = leaf
					}

				} else {
					// TODO: wrap this error in an "invalid subcommand" error or some such.
					return fmt.Errorf("unrecognized subcommand %s", element.String())
				}
			} else {
				panic("illegal state: command node was neither a leaf or a parent")
			}

		case elementTypeLongFlagPair:
			if err := c.interpretLongPair(&element, &unmapped); err != nil {
				return err
			}

		case elementTypeLongFlagSolo:
			if err := c.interpretLongSolo(&element, &unmapped); err != nil {
				return err
			}

		case elementTypeShortBlockSolo:
			if err := c.interpretShortSolo(&element, &unmapped); err != nil {
				return err
			}

		case elementTypeShortBlockPair:
			if err := c.interpretShortPair(&element, &unmapped); err != nil {
				return err
			}

		case elementTypeBoundary:
			c.boundary = true

		case elementTypeEnd:
			break FOR

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	errs := newMultiError()
	var onIncomplete func()

	if node, ok := c.current.(CommandLeaf); !ok {
		if parent, ok := c.current.(CommandParent); ok {
			onIncomplete = parent.onIncomplete
		} else {
			errs.AppendError(fmt.Errorf("command leaf was not reached"))
		}
	} else {
		for _, value := range unmapped {
			node.appendUnmapped(value)
		}

		for _, value := range passthroughs {
			node.appendPassthrough(value)
		}

		for i, arg := range node.Arguments() {
			if arg.IsRequired() && !arg.WasHit() {
				if arg.HasName() {
					errs.AppendError(fmt.Errorf("argument %d (<%s>) is required", i, arg.Name()))
				} else {
					errs.AppendError(fmt.Errorf("argument %d is required", i))
				}
			}
		}
	}

	current := c.current
	for current != nil {
		for _, group := range current.FlagGroups() {
			for _, f := range group.Flags() {
				if f.IsRequired() && !f.WasHit() {
					errs.AppendError(fmt.Errorf("required flag %s was not used", printFlagNames(f)))
				} else if f.RequiresArgument() && !f.Argument().WasHit() {
					errs.AppendError(fmt.Errorf("flag %s requires an argument", printFlagNames(f)))
				}
			}
		}

		current = current.Parent()
	}

	for i, flag := range c.flagHits {
		if flag.isHelpFlag() {
			flag.executeCallback()
			c.flagHits[i] = nil
		}
	}

	for _, flag := range c.flagHits {
		if flag != nil {
			flag.executeCallback()
		}
	}

	if onIncomplete != nil {
		onIncomplete()
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	if c.tree.hasCallback() {
		c.tree.executeCallback()
	}

	for _, b := range c.branches {
		if b.hasCallback() {
			b.executeCallback()
		}
	}

	if c.leaf.hasCallback() {
		c.leaf.executeCallback()
	}

	c.tree.selectCommand(c.leaf)

	return nil
}

func (c *commandTreeInterpreter) interpretShortSolo(element *element, unmapped *[]string) error {
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
			*unmapped = append(*unmapped, strDash+remainder[0:1])
			continue
		}

		c.flagHits = append(c.flagHits, f)

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
				if nextElement.Type == elementTypeEnd {
					return f.hit()
				}

				if nextElement.Type == elementTypeBoundary {
					c.boundary = true
					return f.hit()
				}

				// If we're here then we have a next element, and we're going to
				// sacrifice it to the flag gods.
				return f.hitWithArg(nextElement.String())
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

				case elementTypeEnd:
					return f.hit()

				case elementTypeBoundary:
					c.boundary = true
					return f.hit()

				case elementTypePlainText:
					return f.hitWithArg(nextElement.Data[0])

				case elementTypeShortBlockSolo:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.queue.Offer(nextElement)
						return nil
					} else {
						return f.hitWithArg(nextElement.String())
					}

				case elementTypeShortBlockPair:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.queue.Offer(nextElement)
						return nil
					} else {
						return f.hitWithArg(nextElement.String())
					}

				case elementTypeLongFlagPair:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						c.queue.Offer(nextElement)
						return nil
					} else {
						return f.hitWithArg(nextElement.String())
					}

				case elementTypeLongFlagSolo:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						c.queue.Offer(nextElement)
						return nil
					} else {
						return f.hitWithArg(nextElement.String())
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

func (c *commandTreeInterpreter) interpretShortPair(element *element, unmapped *[]string) error {
	block := element.Data[0]

	if len(block) == 0 {
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.current.FindShortFlag(block[0]); f != nil {
			c.flagHits = append(c.flagHits, f)
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
			*unmapped = append(*unmapped, strDash+block[0:1])
			continue
		}

		c.flagHits = append(c.flagHits, f)

		if f.RequiresArgument() {
			if h {
				return f.hitWithArg(block[1:] + "=" + element.Data[1])
			} else {
				return f.hitWithArg(element.Data[1])
			}
		}

		if f.HasArgument() {
			if !h {
				return f.hitWithArg(element.Data[1])
			}

			if c.current.FindShortFlag(block[1]) != nil {
				return f.hit()
			}

			return f.hitWithArg(block[1:] + "=" + element.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !h {
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

func (c *commandTreeInterpreter) interpretLongSolo(element *element, unmapped *[]string) error {
	f := c.current.FindLongFlag(element.Data[0])

	if f == nil {
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.flagHits = append(c.flagHits, f)

	if f.RequiresArgument() {
		nextElement := c.next()

		if nextElement.Type == elementTypeEnd {
			return f.hit()
		}

		if nextElement.Type == elementTypeBoundary {
			c.boundary = true
			return nil
		}

		return f.hitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.next()

		switch nextElement.Type {

		case elementTypeEnd:
			return f.hit()

		case elementTypeBoundary:
			c.boundary = true
			return f.hit()

		case elementTypePlainText:
			return f.hitWithArg(nextElement.Data[0])

		case elementTypeLongFlagSolo:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				c.queue.Offer(nextElement)
				return nil
			} else {
				return f.hitWithArg(nextElement.String())
			}

		case elementTypeLongFlagPair:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				c.queue.Offer(nextElement)
				return nil
			} else {
				return f.hitWithArg(nextElement.String())
			}

		case elementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.queue.Offer(nextElement)
				return nil
			} else {
				return f.hitWithArg(nextElement.String())
			}

		case elementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.queue.Offer(nextElement)
				return nil
			} else {
				return f.hitWithArg(nextElement.String())
			}

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	return f.hit()
}

func (c *commandTreeInterpreter) interpretLongPair(element *element, unmapped *[]string) error {
	flag := c.current.FindLongFlag(element.Data[0])

	if flag == nil {
		*unmapped = append(*unmapped, element.String())
	} else {
		c.flagHits = append(c.flagHits, flag)

		if flag.HasArgument() {
			return flag.hitWithArg(element.Data[1])
		} else {
			// TODO: this should be a warning for a flag getting an argument it didn't expect, this also applies to shortpair
		}
	}

	return nil
}
