package argo

import "reflect"

type commandInterpreter struct {
	parser   parser
	command  Command
	flagHits []Flag
	elements deque[element]
	boundary bool
}

func (c *commandInterpreter) nextElement() element {
	if c.elements.IsEmpty() {
		c.elements.Offer(c.parser.Next())
	}

	return c.elements.Poll()
}

func (c *commandInterpreter) Run() error {
	var err error

FOR:
	for {
		element := c.nextElement()

		// If we've already hit the boundary marker then everything else is just a
		// passthrough.
		if c.boundary {
			if element.Type == elementTypeEnd {
				break
			}

			c.command.appendPassthrough(element.String())
			continue
		}

		switch element.Type {

		case elementTypePlainText:
			if err = c.command.appendArgument(element.String()); err != nil {
				return err
			}

		case elementTypeShortBlockSolo:
			if c.boundary, err = c.interpretShortSolo(&element); err != nil {
				return err
			}

		case elementTypeShortBlockPair:
			if c.boundary, err = c.interpretShortPair(&element); err != nil {
				return err
			}

		case elementTypeLongFlagSolo:
			if c.boundary, err = c.interpretLongSolo(&element); err != nil {
				return err
			}

		case elementTypeLongFlagPair:
			if c.boundary, err = c.interpretLongPair(&element); err != nil {
				return err
			}

		case elementTypeBoundary:
			c.boundary = true
			continue

		case elementTypeEnd:
			break FOR

		default:
			panic("illegal state")

		}
	}

	errs := newMultiError()

	flagGroups := c.command.FlagGroups()
	for i := range flagGroups {
		flagGroup := flagGroups[i].Flags()

		for j := range flagGroup {
			f := flagGroup[j]

			if f.IsRequired() && !f.WasHit() {
				errs.AppendError(newMissingFlagError(f))
			}

			if f.WasHit() && f.RequiresArgument() && !f.Argument().WasHit() {
				errs.AppendError(newMissingRequiredFlagArgumentError(f.Argument(), f, c.command))
			}

			if !f.WasHit() && f.HasArgument() && f.Argument().HasDefault() {
				if err := f.Argument().setToDefault(); err != nil {
					errs.AppendError(err)
				}
			}
		}
	}

	arguments := c.command.Arguments()
	for i := range arguments {
		arg := arguments[i]

		if arg.IsRequired() && !arg.WasHit() {
			errs.AppendError(newMissingRequiredPositionalArgumentError(arg, c.command))
		}
		if !arg.WasHit() && arg.HasDefault() {
			if err := arg.setToDefault(); err != nil {
				errs.AppendError(err)
			}
		}
	}

	for i, flag := range c.flagHits {
		if flag.isHelpFlag() {
			flag.executeCallback()
			c.flagHits[i] = nil
			break
		}
	}

	for _, flag := range c.flagHits {
		if flag != nil {
			flag.executeCallback()
		}
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	return nil
}

func (c *commandInterpreter) interpretShortSolo(e *element) (bool, error) {
	remainder := e.Data[0]

	for i := 0; i < len(e.Data[0]); i++ {
		// has next
		h := i+1 < len(e.Data[0])
		// short flag byte
		b := remainder[0]

		// Look up the flag in the short flag map
		f := c.command.FindShortFlag(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			c.command.appendUnmapped(strDash + remainder[0:1])
			remainder = remainder[1:]
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
				nextElement := c.nextElement()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == elementTypeEnd {
					return false, f.hit()
				}

				if nextElement.Type == elementTypeBoundary {
					return true, f.hit()
				}

				// If we're here then we have a next element, and we're going to
				// sacrifice it to the flag gods.
				return false, f.hitWithArg(nextElement.String())
			}

			// So we have at least one more character in this block.  Eat that and
			// anything else as the flag argument.
			return false, f.hitWithArg(remainder[1:])
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
				if t := c.command.FindShortFlag(n); t != nil {
					remainder = remainder[1:]
					continue
				} else

				// Since there is no flag matching the next character, then we have to
				// assume that the remaining text is the argument for the flag
				{
					return false, f.hitWithArg(remainder[1:])
				}
			} else {

				nextElement := c.nextElement()

				switch nextElement.Type {

				case elementTypeEnd:
					if f.HasArgument() && f.Argument().HasBinding() && f.Argument().BindingType().Kind() == reflect.Bool {
						return false, f.hitWithArg("true")
					}
					return false, f.hit()

				case elementTypeBoundary:
					if f.HasArgument() && f.Argument().HasBinding() && f.Argument().BindingType().Kind() == reflect.Bool {
						return false, f.hitWithArg("true")
					}
					return true, f.hit()

				case elementTypePlainText:
					return false, f.hitWithArg(nextElement.Data[0])

				case elementTypeShortBlockSolo:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.elements.Offer(nextElement)
						return false, nil
					} else {
						return false, f.hitWithArg(nextElement.String())
					}

				case elementTypeShortBlockPair:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.elements.Offer(nextElement)
						return false, nil
					} else {
						return false, f.hitWithArg(nextElement.String())
					}

				case elementTypeLongFlagPair:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						c.elements.Offer(nextElement)
						return false, nil
					} else {
						return false, f.hitWithArg(nextElement.String())
					}

				case elementTypeLongFlagSolo:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						c.elements.Offer(nextElement)
						return false, nil
					} else {
						return false, f.hitWithArg(nextElement.String())
					}

				default:
					panic("illegal state")

				}
			}
		}

		// The flag doesn't expect an argument, just hit it
		if err := f.hit(); err != nil {
			return false, err
		}

		// if has next
		remainder = remainder[1:]
	}

	return false, nil
}

func (c *commandInterpreter) interpretShortPair(e *element) (bool, error) {
	block := e.Data[0]

	if len(block) == 0 {
		c.command.appendUnmapped(e.String())
		return false, nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.command.FindShortFlag(block[0]); f != nil {
			return false, f.hitWithArg(e.Data[1])
		} else {
			c.command.appendUnmapped(e.String())
			return false, nil
		}
	}

	for i := 0; i < len(e.Data[0]); i++ {
		// has next character
		h := i+1 < len(e.Data[0])
		// current character
		b := block[0]

		f := c.command.FindShortFlag(b)

		if f == nil {
			c.command.appendUnmapped(strDash + block[0:1])
			continue
		}

		if f.RequiresArgument() {
			if h {
				return false, f.hitWithArg(block[1:] + "=" + e.Data[1])
			} else {
				return false, f.hitWithArg(e.Data[1])
			}
		}

		if f.HasArgument() {
			if !h {
				return false, f.hitWithArg(e.Data[1])
			}

			if c.command.FindShortFlag(block[1]) != nil {
				return false, f.hit()
			}

			return false, f.hitWithArg(block[1:] + "=" + e.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !h {
			return false, f.hitWithArg(e.Data[1])
		}

		// Well, now that's out of the way, we can move on to the next flag (after
		// we mark this one as hit of course).
		if err := f.hit(); err != nil {
			return false, err
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c *commandInterpreter) interpretLongSolo(e *element) (bool, error) {
	f := c.command.FindLongFlag(e.Data[0])

	if f == nil {
		c.command.appendUnmapped(e.String())
		return false, nil
	}

	if f.RequiresArgument() {
		nextElement := c.parser.Next()

		if nextElement.Type == elementTypeEnd {
			return false, f.hit()
		}

		if nextElement.Type == elementTypeBoundary {
			return true, f.hit()
		}

		return false, f.hitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.nextElement()

		switch nextElement.Type {

		case elementTypeEnd:
			return false, f.hit()

		case elementTypeBoundary:
			return true, f.hit()

		case elementTypePlainText:
			return false, f.hitWithArg(nextElement.Data[0])

		case elementTypeLongFlagSolo:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			} else {
				return false, f.hitWithArg(nextElement.String())
			}

		case elementTypeLongFlagPair:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			} else {
				return false, f.hitWithArg(nextElement.String())
			}

		case elementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			} else {
				return false, f.hitWithArg(nextElement.String())
			}

		case elementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			} else {
				return false, f.hitWithArg(nextElement.String())
			}

		default:
			panic("illegal state")
		}
	}

	return false, f.hit()
}

func (c *commandInterpreter) interpretLongPair(e *element) (bool, error) {
	flag := c.command.FindLongFlag(e.Data[0])

	if flag == nil {
		c.command.appendUnmapped(e.String())
	} else {
		c.flagHits = append(c.flagHits, flag)

		if flag.HasArgument() {
			return false, flag.hitWithArg(e.Data[1])
		}
		// TODO: this should be a warning, the flag didn't expect an argument (same applies to shortpair)
	}

	return false, nil
}
