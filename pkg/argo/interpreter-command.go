package argo

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/parse"
	"github.com/Foxcapades/Argonaut/internal/util"
)

type commandInterpreter struct {
	parser   parse.Parser
	command  Command
	flagHits flagQueue
	elements util.Deque[parse.Element]
	boundary bool
}

func (c *commandInterpreter) nextElement() parse.Element {
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
			if element.Type == parse.ElementTypeEnd {
				break
			}

			c.command.appendPassthrough(element.String())
			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if err = c.command.appendArgument(element.String()); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockSolo:
			if c.boundary, err = c.interpretShortSolo(&element); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockPair:
			if c.boundary, err = c.interpretShortPair(&element); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagSolo:
			if c.boundary, err = c.interpretLongSolo(&element); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagPair:
			if c.boundary, err = c.interpretLongPair(&element); err != nil {
				return err
			}

		case parse.ElementTypeBoundary:
			c.boundary = true
			continue

		case parse.ElementTypeEnd:
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

	var it = c.flagHits.iterator()
	var hf = 0
	for it.hasNext() {
		flag := it.next()
		if flag.isHelpFlag() {
			flag.executeCallback()
			break
		}
		hf++
	}

	it = c.flagHits.iterator()
	var cf = 0
	for it.hasNext() {
		if cf != hf {
			it.next().executeCallback()
		}
		cf++
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	c.command.executeCallback()

	return nil
}

func (c *commandInterpreter) interpretShortSolo(e *parse.Element) (bool, error) {
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
			c.command.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			c.command.appendUnmapped(chars.StrDash + remainder[0:1])
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
				nextElement := c.nextElement()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					if hasBooleanArgument(f) {
						return false, f.hitWithArg("true")
					}
					return false, f.hit()
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					if hasBooleanArgument(f) {
						return true, f.hitWithArg("true")
					}
					return true, f.hit()
				}

				// If we're here then we have a next element, and we're going to try and
				// sacrifice it to the flag gods.
				if err := f.hitWithArg(nextElement.String()); err != nil {
					c.elements.Offer(nextElement)

					if hasBooleanArgument(f) {
						return false, f.hitWithArg("true")
					}

					return false, f.hit()
				} else {
					return false, nil
				}

				return false, f.hitWithArg(nextElement.String())
			}

			if hasBooleanArgument(f) {
				possibleNextFlag := c.command.FindShortFlag(remainder[1])
				if possibleNextFlag != nil {
					if err := f.hitWithArg("true"); err != nil {
						return false, err
					}
					continue
				}
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
					if err := f.hit(); err != nil {
						return false, err
					}
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

				case parse.ElementTypeEnd:
					if hasBooleanArgument(f) {
						return false, f.hitWithArg("true")
					}
					return false, f.hit()

				case parse.ElementTypeBoundary:
					if hasBooleanArgument(f) {
						return true, f.hitWithArg("true")
					}
					return true, f.hit()

				case parse.ElementTypePlainText:
					// If the flag expects an argument, but the value following the flag
					// cannot be parsed as the argument, then the argument value is not
					// treated as the value to the flag and is instead treated as a
					// positional argument value.
					if err := f.hitWithArg(nextElement.Data[0]); err != nil {
						c.elements.Offer(nextElement)
					}

					return false, f.hit()

				case parse.ElementTypeShortBlockSolo:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.elements.Offer(nextElement)
						return false, f.hit()
					}

					if err := f.hitWithArg(nextElement.String()); err != nil {
						c.elements.Offer(nextElement)
					}

					return false, f.hit()

				case parse.ElementTypeShortBlockPair:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						c.elements.Offer(nextElement)
						return false, f.hit()
					}

					if err := f.hitWithArg(nextElement.String()); err != nil {
						c.elements.Offer(nextElement)
					}

					return false, f.hit()

				case parse.ElementTypeLongFlagPair:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						c.elements.Offer(nextElement)
						return false, f.hit()
					}

					if err := f.hitWithArg(nextElement.String()); err != nil {
						c.elements.Offer(nextElement)
					}

					return false, f.hit()

				case parse.ElementTypeLongFlagSolo:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						c.elements.Offer(nextElement)
						return false, f.hit()
					}

					if err := f.hitWithArg(nextElement.String()); err != nil {
						c.elements.Offer(nextElement)
					}

					return false, f.hit()

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

func (c *commandInterpreter) interpretShortPair(e *parse.Element) (bool, error) {
	block := e.Data[0]

	if len(block) == 0 {
		c.command.appendUnmapped(e.String())
		return false, nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.command.FindShortFlag(block[0]); f != nil {
			c.flagHits.append(f)
			return false, f.hitWithArg(e.Data[1])
		}

		c.command.appendUnmapped(e.String())
		return false, nil
	}

	for i := 0; i < len(e.Data[0]); i++ {
		// has next character
		h := i+1 < len(e.Data[0])
		// current character
		b := block[0]

		f := c.command.FindShortFlag(b)

		if f == nil {
			c.command.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			c.command.appendUnmapped(chars.StrDash + block[0:1])
			continue
		}

		c.flagHits.append(f)

		if f.RequiresArgument() {
			if h {
				return false, f.hitWithArg(block[1:] + "=" + e.Data[1])
			}

			return false, f.hitWithArg(e.Data[1])
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
			c.command.AppendWarning(fmt.Sprintf("flag -%c received an argument it didn't expect", b))
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

func (c *commandInterpreter) interpretLongSolo(e *parse.Element) (bool, error) {
	f := c.command.FindLongFlag(e.Data[0])

	if f == nil {
		c.command.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", e.Data[0]))
		c.command.appendUnmapped(e.String())
		return false, nil
	}

	c.flagHits.append(f)

	if f.RequiresArgument() {
		nextElement := c.parser.Next()

		if nextElement.Type == parse.ElementTypeEnd {
			return false, f.hit()
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			return true, f.hit()
		}

		return false, f.hitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.nextElement()

		switch nextElement.Type {

		case parse.ElementTypeEnd:
			return false, f.hit()

		case parse.ElementTypeBoundary:
			return true, f.hit()

		case parse.ElementTypePlainText:
			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

			return false, f.hit()

		case parse.ElementTypeLongFlagSolo:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			}

			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

			return false, f.hit()

		case parse.ElementTypeLongFlagPair:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			}

			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

			return false, f.hit()

		case parse.ElementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			}

			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

			return false, f.hit()

		case parse.ElementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.elements.Offer(nextElement)
				return false, f.hit()
			}

			if err := f.hitWithArg(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

			return false, f.hit()

		default:
			panic("illegal state")
		}
	}

	return false, f.hit()
}

func (c *commandInterpreter) interpretLongPair(e *parse.Element) (bool, error) {
	flag := c.command.FindLongFlag(e.Data[0])

	if flag == nil {
		c.command.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", e.Data[0]))
		c.command.appendUnmapped(e.String())
	} else {
		c.flagHits.append(flag)

		if flag.HasArgument() {
			return false, flag.hitWithArg(e.Data[1])
		}
		c.command.AppendWarning(fmt.Sprintf("flag --%s received an argument it didn't expect", e.Data[0]))
		return false, flag.hit()
	}

	return false, nil
}
