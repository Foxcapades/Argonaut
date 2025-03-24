package tree

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func (c *CommandTreeInterpreter) interpretShortSolo(element *parse.Element, unmapped *[]string) error {
	curNode := c.current.(flag.GroupContainer)
	remainder := element.Data[0]

	for i := 0; i < len(element.Data[0]); i++ {
		// has next
		h := i+1 < len(element.Data[0])
		// short flag byte
		b := remainder[0]
		remainder = remainder[1:]

		// Look up the flag in the short flag map
		f := curNode.FindShortFlag(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			common.AppendUnrecognizedShortFlag(&c.result, b)
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
				possibleNextFlag := curNode.FindShortFlag(remainder[0])

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
			if t := curNode.FindShortFlag(n); t != nil {
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

		return c.tryConsumeNextElementForOptionalArgument(curNode, arg)
	}

	return nil
}

// interpretShortPair tries to make sense of a pair where the first value is a
// block of one or more short flags, and the second value is an argument value
// that was directly attached using an `=` character.
func (c *CommandTreeInterpreter) interpretShortPair(element *parse.Element, unmapped *[]string) error {
	block := element.Data[0]
	curNode := c.current.(flag.GroupContainer)

	if len(block) == 0 {
		common.AppendWarning(&c.result, "blank short flag name", argo.UnrecognizedFlag)
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := curNode.FindShortFlag(block[0]); f != nil {
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

		f := curNode.FindShortFlag(b)

		if f == nil {
			common.AppendUnrecognizedShortFlag(&c.result, b)
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
			if curNode.FindShortFlag(block[1]) != nil {
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
			common.AppendWarning(&c.result, fmt.Sprintf("flag -%c received an argument it didn't expect", b), argo.UnexpectedFlagArgument)
			return f.Argument().SetValue(element.Data[1])
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c *CommandTreeInterpreter) interpretLongSolo(element *parse.Element, unmapped *[]string) error {
	curNode := c.current.(flag.GroupContainer)
	f := curNode.FindLongFlag(element.Data[0])

	if f == nil {
		common.AppendUnrecognizedLongFlag(&c.result, element.Data[0])
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

	return c.tryConsumeNextElementForOptionalArgument(curNode, arg)
}

func (c *CommandTreeInterpreter) interpretLongPair(element *parse.Element, unmapped *[]string) error {
	targetFlag := c.current.(flag.GroupContainer).FindLongFlag(element.Data[0])

	if targetFlag == nil {
		common.AppendUnrecognizedLongFlag(&c.result, element.Data[0])
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.flagHits.Append(targetFlag)
	targetFlag.IncrementHitCount()

	if targetFlag.HasArgument() {
		return targetFlag.Argument().SetValue(element.Data[1])
	} else {
		common.AppendWarning(
			&c.result,
			fmt.Sprintf("flag --%s received an argument it didn't expect", element.Data[0]),
			argo.UnexpectedFlagArgument,
		)
	}

	return nil
}

func (c *CommandTreeInterpreter) tryConsumeNextElementForOptionalArgument(curNode flag.GroupContainer, arg argo.Argument) error {
	nextElement := c.next()

	switch nextElement.Type {

	case parse.ElementTypeBoundary:
		c.boundary = true
		fallthrough

	case parse.ElementTypeEnd:
		if argument.IsBoolean(arg) {
			return arg.SetValue("true")
		}

	case parse.ElementTypePlainText:
		if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
		if len(nextElement.Data[0]) > 0 && curNode.FindShortFlag(nextElement.Data[0][0]) != nil {
			c.queue.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	case parse.ElementTypeLongFlagSolo, parse.ElementTypeLongFlagPair:
		if curNode.FindLongFlag(nextElement.Data[0]) != nil {
			c.queue.Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.queue.Offer(nextElement)
		}

	default:
		panic("illegal state: unrecognized parser element type")
	}

	return nil
}
