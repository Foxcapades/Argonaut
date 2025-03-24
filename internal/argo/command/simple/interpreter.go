package command

import (
	"fmt"

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

func NewInterpreter(args []string, command argo.Command) Interpreter {
	return Interpreter{
		parser:   parse.NewParser(emit.NewEmitter(args)),
		command:  command,
		elements: utils.NewDeque[parse.Element](2),
		flagHits: flag.NewQueue(),
	}
}

type Interpreter struct {
	boundary bool
	parser   parse.Parser
	flagHits flag.Queue
	result   argo.ParseResult
	command  argo.Command
	elements utils.Deque[parse.Element]
}

func (c *Interpreter) nextElement() parse.Element {
	if c.elements.IsEmpty() {
		c.elements.Offer(c.parser.Next())
	}

	return c.elements.Poll()
}

func (c *Interpreter) Run() (argo.ParseResult, error) {
	var argumentStream argument.Appender
	var interpretFn func(*parse.Element) error

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
				return common.AppendError(c.result, err)
			} else if !ok {
				c.command.AppendUnmappedInput(element.String())
			}

			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if ok, err := argumentStream.Append(element.String()); err != nil {
				return common.AppendError(c.result, err)
			} else if !ok {
				c.command.AppendUnmappedInput(element.String())
				continue
			}

		case parse.ElementTypeShortBlockSolo:
			interpretFn = c.interpretShortSolo

		case parse.ElementTypeShortBlockPair:
			interpretFn = c.interpretShortPair

		case parse.ElementTypeLongFlagSolo:
			interpretFn = c.interpretLongSolo

		case parse.ElementTypeLongFlagPair:
			interpretFn = c.interpretLongPair

		case parse.ElementTypeBoundary:
			c.boundary = true
			continue

		case parse.ElementTypeEnd:
			break FOR

		default:
			panic("illegal state")
		}

		if err := interpretFn(&element); err != nil {
			return common.AppendError(c.result, err)
		}
	}

	common.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())

	errs := xerr.NewMultiError()

	common.CheckRequiredFlags(c.command.FlagGroups(true), errs)
	common.CheckRequiredArguments(c.command.Arguments(), errs)

	if len(errs.Errors()) > 0 {
		return common.AppendError(c.result, errs)
	}

	c.command.Callback()(c.command)

	return c.result, nil
}

func (c *Interpreter) interpretShortSolo(e *parse.Element) error {
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
			common.AppendWarning(&c.result, fmt.Sprintf("unrecognized short flag -%c", b), argo.UnrecognizedFlag)
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

func (c *Interpreter) interpretShortPair(cliElement *parse.Element) error {
	flagBlock := cliElement.Data[0]

	if len(flagBlock) == 0 {
		c.command.AppendUnmappedInput(cliElement.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(flagBlock) == 1 {
		if f := c.command.FindShortFlag(flagBlock[0]); f != nil {
			c.flagHits.Append(f)
			f.IncrementHitCount()
			if f.HasArgument() {
				return f.Argument().SetValue(cliElement.Data[1])
			}

			common.AppendWarning(&c.result, fmt.Sprintf("targetFlag -%c received an argument it didn't expect", flagBlock[0]), argo.UnexpectedFlagArgument)
			return nil
		}

		c.command.AppendUnmappedInput(cliElement.String())
		return nil
	}

	for i := 0; i < len(cliElement.Data[0]); i++ {
		// has next character
		hasNextChar := i+1 < len(cliElement.Data[0])

		// current character
		currentChar := flagBlock[0]

		// Remove the current character from the remaining flag character block.
		flagBlock = flagBlock[1:]

		targetFlag := c.command.FindShortFlag(currentChar)

		if targetFlag == nil {
			common.AppendWarning(&c.result, fmt.Sprintf("unrecognized short flag -%c", currentChar), argo.UnrecognizedFlag)
			c.command.AppendUnmappedInput(chars.StrDash + string(currentChar))
			continue
		}

		c.flagHits.Append(targetFlag)
		targetFlag.IncrementHitCount()

		if !targetFlag.HasArgument() {
			common.AppendWarning(&c.result, fmt.Sprintf("targetFlag -%c received an argument it didn't expect", currentChar), argo.UnexpectedFlagArgument)
			return nil
		}

		arg := targetFlag.Argument()

		if arg.IsRequired() {
			if hasNextChar {
				return arg.SetValue(flagBlock + "=" + cliElement.Data[1])
			}

			return arg.SetValue(cliElement.Data[1])
		}

		if !hasNextChar {
			return arg.SetValue(cliElement.Data[1])
		}

		// Look up the next character in the block and check if it is a flag too.
		// If so, prioritize that flag over an optional arg value.
		if c.command.FindShortFlag(flagBlock[0]) != nil {
			return nil
		}

		return arg.SetValue(flagBlock + "=" + cliElement.Data[1])
	}

	panic("illegal state")
}

func (c *Interpreter) interpretLongSolo(cliElement *parse.Element) error {
	targetFlag := c.command.FindLongFlag(cliElement.Data[0])

	if targetFlag == nil {
		common.AppendWarning(&c.result, fmt.Sprintf("unrecognized long flag --%s", cliElement.Data[0]), argo.UnrecognizedFlag)
		c.command.AppendUnmappedInput(cliElement.String())
		return nil
	}

	c.flagHits.Append(targetFlag)
	targetFlag.IncrementHitCount()

	if !targetFlag.HasArgument() {
		return nil
	}

	arg := targetFlag.Argument()

	if arg.IsRequired() {
		nextElement := c.parser.Next()

		switch nextElement.Type {

		case parse.ElementTypeBoundary:
			c.boundary = true
			fallthrough

		case parse.ElementTypeEnd:
			if argument.IsBoolean(arg) {
				return arg.SetValue("true")
			}

		default:
			return arg.SetValue(nextElement.String())
		}
	}

	nextElement := c.nextElement()

	switch nextElement.Type {

	case parse.ElementTypeBoundary:
		c.boundary = true

	case parse.ElementTypeEnd:
		// nothing to do

	case parse.ElementTypePlainText:
		if err := arg.SetValue(nextElement.String()); err != nil {
			c.elements.Offer(nextElement)
		}

	case parse.ElementTypeLongFlagSolo, parse.ElementTypeLongFlagPair:
		if c.command.FindLongFlag(nextElement.Data[0]) != nil {
			c.elements.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.elements.Offer(nextElement)
		}

	case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
		if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
			c.elements.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.elements.Offer(nextElement)
		}

	default:
		panic("illegal state")
	}

	return nil
}

func (c *Interpreter) interpretLongPair(e *parse.Element) error {
	targetFlag := c.command.FindLongFlag(e.Data[0])

	if targetFlag == nil {
		common.AppendWarning(&c.result, fmt.Sprintf("unrecognized long targetFlag --%s", e.Data[0]), argo.UnrecognizedFlag)
		c.command.AppendUnmappedInput(e.String())
	} else {
		c.flagHits.Append(targetFlag)
		targetFlag.IncrementHitCount()

		if targetFlag.HasArgument() {
			return targetFlag.Argument().SetValue(e.Data[1])
		}

		common.AppendWarning(&c.result, fmt.Sprintf("targetFlag --%s received an argument it didn't expect", e.Data[0]), argo.UnexpectedFlagArgument)
	}

	return nil
}
