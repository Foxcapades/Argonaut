package command

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/emit"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func newCommandInterpreter(args []string, command argo.Command) commandInterpreter {
	return commandInterpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		command:  command,
		elements: utils.NewDeque[parse.Element](2),
		flagHits: flag.NewQueue(),
	}
}

type commandInterpreter struct {
	parser   parse.Parser
	command  argo.Command
	flagHits flag.Queue
	elements utils.Deque[parse.Element]
	boundary bool
}

func (c *commandInterpreter) nextElement() parse.Element {
	if c.elements.IsEmpty() {
		c.elements.Offer(c.parser.Next())
	}

	return c.elements.Poll()
}

func (c *commandInterpreter) Run() error {
	var argumentStream argument.Appender

FOR:
	for {
		element := c.nextElement()

		// If we've already hit the boundary marker then everything else is just a
		// passthrough.
		if c.boundary {
			if element.Type == parse.ElementTypeEnd {
				break
			}

			if ok, err := argumentStream.Append(element.String()); err != nil {
				return err
			} else if !ok {
				c.command.AppendUnmappedInput(element.String())
			}

			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if ok, err := argumentStream.Append(element.String()); err != nil {
				return err
			} else if !ok {
				c.command.AppendUnmappedInput(element.String())
			}

		case parse.ElementTypeShortBlockSolo:
			if err := c.interpretShortSolo(&element); err != nil {
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

	errs := xerr.NewMultiError()

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

	var helpFlagIndex = 0
	for flag := range c.flagHits.Iterator() {
		if flag.IsHelpFlag() {
			flag.Callback()(flag)
			break
		}
		helpFlagIndex++
	}

	var currentFlagIndex = 0
	for flag := range c.flagHits.Iterator() {
		if currentFlagIndex != helpFlagIndex {
			flag.Callback()(flag)
		}
		currentFlagIndex++
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	c.command.Callback()(c.command)

	return nil
}

func (c *commandInterpreter) interpretShortSolo(e *parse.Element) error {
	remainder := e.Data[0]

	for i := 0; i < len(e.Data[0]); i++ {
		// has next
		h := i+1 < len(e.Data[0])
		// short flag byte
		b := remainder[0]
		remainder = remainder[1:]

		// Look up the flag in the short flag map
		f := c.command.FindShortFlag(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			// TODO: c.command.AppendWarning(fmt.Sprintf("unrecognized short flag -%c", b))
			c.command.AppendUnmappedInput(chars.StrDash + string(b))
			continue
		}

		c.flagHits.Append(f)
		f.IncrementHitCount()

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
				nextElement := c.nextElement()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					if argument.IsBoolean(arg) {
						return arg.SetValue("true")
					}

					return nil
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					c.boundary = true

					if argument.IsBoolean(arg) {
						return arg.SetValue("true")
					}

					return nil
				}

				// If we're here then we have a next element, and we're going to try and
				// sacrifice it to the flag gods.
				return arg.SetValue(nextElement.String())
			}

			if argument.IsBoolean(arg) {
				possibleNextFlag := c.command.FindShortFlag(remainder[0])

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
			if t := c.command.FindShortFlag(n); t != nil {
				continue
			}

			// Since there is no flag matching the next character, then we have to
			// assume that the remaining text is the argument for the flag
			return arg.SetValue(remainder)
		}

		nextElement := c.nextElement()

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
			// If the flag expects an argument, but the value following the flag
			// cannot be parsed as the argument, then the argument value is not
			// treated as the value to the flag and is instead treated as a
			// positional argument value.
			if err := arg.SetValue(nextElement.Data[0]); err != nil {
				c.elements.Offer(nextElement)
			}

		case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
			if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				c.elements.Offer(nextElement)
			} else if err := arg.SetValue(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

		case parse.ElementTypeLongFlagPair, parse.ElementTypeLongFlagSolo:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				c.elements.Offer(nextElement)
			} else if err := arg.SetValue(nextElement.String()); err != nil {
				c.elements.Offer(nextElement)
			}

		default:
			panic("illegal state")
		}

		break
	}

	return nil
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
		// TODO: c.command.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", e.Data[0]))
		c.command.AppendUnmappedInput(e.String())
		return false, nil
	}

	c.flagHits.Append(f)

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
		// TODO: c.command.AppendWarning(fmt.Sprintf("unrecognized long flag --%s", e.Data[0]))
		c.command.AppendUnmappedInput(e.String())
	} else {
		c.flagHits.Append(flag)

		if flag.HasArgument() {
			return false, flag.hitWithArg(e.Data[1])
		}
		// TODO: c.command.AppendWarning(fmt.Sprintf("flag --%s received an argument it didn't expect", e.Data[0]))
		return false, flag.hit()
	}

	return false, nil
}
