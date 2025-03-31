package flag

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ContainingInterpreter interface {
	FlagIndex() Index

	Next() parse.Element

	FlagHits() *Queue

	HitBoundary()

	ElementQueue() utils.Deque[parse.Element]
}

func InterpretShortSolo(c ContainingInterpreter, element *parse.Element, unmapped *[]string, warnings *[]argo.InputWarning) error {
	remainder := element.Data[0]

	for i := 0; i < len(element.Data[0]); i++ {
		// has next
		h := i+1 < len(element.Data[0])
		// short flag byte
		b := remainder[0]
		remainder = remainder[1:]

		// Look up the flag in the short flag map
		f := c.FlagIndex().ByShortName(b)

		// If the flag was not found, append the arg to the unmapped slice and move
		// on to the next character.
		if f == nil {
			AppendUnrecognizedShortFlag(warnings, b)
			*unmapped = append(*unmapped, text.DashString+string(b))
			continue
		}

		f.IncrementHitCount()
		c.FlagHits().Append(f)

		if !f.HasArgument() {
			continue
		}

		arg := f.Argument()

		// If the flag we found requires an argument, eat the rest of the block and
		// pass it to the Hit method.  Since the block will have been consumed after
		// this, return here.
		if arg.IsRequired() {

			// If we don't have any more characters in this short block, then we have
			// to consume the next element as the argument for this flag.
			if !h {
				nextElement := c.Next()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					if f.HasArgument() && argument.IsBoolean(arg) {
						return arg.SetValue("true")
					}

					return nil
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					c.HitBoundary()

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
				possibleNextFlag := c.FlagIndex().ByShortName(remainder[0])

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
			if t := c.FlagIndex().ByShortName(n); t != nil {
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

		return tryConsumeNextElementForOptionalArgument(c, arg)
	}

	return nil
}

// InterpretShortPair tries to make sense of a pair where the first value is a
// block of one or more short flags, and the second value is an argument value
// that was directly attached using an `=` character.
func InterpretShortPair(c ContainingInterpreter, element *parse.Element, unmapped *[]string, warnings *[]argo.InputWarning) error {
	block := element.Data[0]

	if len(block) == 0 {
		*warnings = append(*warnings, argo.InputWarning{
			Type:    argo.UnrecognizedFlag,
			Message: "blank short flag name",
		})
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.FlagIndex().ByShortName(block[0]); f != nil {
			c.FlagHits().Append(f)
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

		f := c.FlagIndex().ByShortName(b)

		if f == nil {
			AppendUnrecognizedShortFlag(warnings, b)

			// If we are on the last byte of the flag name segment, and we didn't find
			// anything, then include the value in the unmapped, and return here as
			// there is nothing left to do.
			if !hasNextChar {
				*unmapped = append(*unmapped, text.DashString+block[0:1]+element.Data[1])
				return nil
			}

			*unmapped = append(*unmapped, text.DashString+block[0:1])
			block = block[1:]
			continue
		}

		c.FlagHits().Append(f)
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
			if c.FlagIndex().ByShortName(block[1]) != nil {
				block = block[1:]
				continue
			}

			// If the next character in the block does not match any known short flag,
			// assume that the whole remaining value is part of the value.
			return f.Argument().SetValue(block[1:] + "=" + element.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		if !hasNextChar {
			*warnings = append(*warnings, argo.InputWarning{
				Type:    argo.UnexpectedFlagArgument,
				Message: fmt.Sprintf("flag -%c received an argument it didn't expect", b),
			})
			return nil
		}

		block = block[1:]
	}

	panic("illegal state")
}

func InterpretLongSolo(c ContainingInterpreter, element *parse.Element, unmapped *[]string, warnings *[]argo.InputWarning) error {
	f := c.FlagIndex().ByLongName(element.Data[0])

	if f == nil {
		AppendUnrecognizedLongFlag(warnings, element.Data[0])
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.FlagHits().Append(f)
	f.IncrementHitCount()

	if !f.HasArgument() {
		return nil
	}

	arg := f.Argument()

	if arg.IsRequired() {
		nextElement := c.Next()

		if nextElement.Type == parse.ElementTypeEnd {
			return nil
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			c.HitBoundary()
			return nil
		}

		return arg.SetValue(nextElement.String())
	}

	return tryConsumeNextElementForOptionalArgument(c, arg)
}

func InterpretLongPair(c ContainingInterpreter, element *parse.Element, unmapped *[]string, warnings *[]argo.InputWarning) error {
	targetFlag := c.FlagIndex().ByLongName(element.Data[0])

	if targetFlag == nil {
		AppendUnrecognizedLongFlag(warnings, element.Data[0])
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	c.FlagHits().Append(targetFlag)
	targetFlag.IncrementHitCount()

	if targetFlag.HasArgument() {
		return targetFlag.Argument().SetValue(element.Data[1])
	} else {
		*warnings = append(*warnings, argo.InputWarning{
			Type:    argo.UnexpectedFlagArgument,
			Message: fmt.Sprintf("flag --%s received an argument it didn't expect", element.Data[0]),
		})
	}

	return nil
}

func tryConsumeNextElementForOptionalArgument(c ContainingInterpreter, arg argo.Argument) error {
	nextElement := c.Next()

	switch nextElement.Type {

	case parse.ElementTypeBoundary:
		c.HitBoundary()
		fallthrough

	case parse.ElementTypeEnd:
		if argument.IsBoolean(arg) {
			return arg.SetValue("true")
		}

	case parse.ElementTypePlainText:
		if err := arg.SetValue(nextElement.String()); err != nil {
			c.ElementQueue().Offer(nextElement)
		}

	case parse.ElementTypeShortBlockSolo, parse.ElementTypeShortBlockPair:
		if len(nextElement.Data[0]) > 0 && c.FlagIndex().ByShortName(nextElement.Data[0][0]) != nil {
			c.ElementQueue().Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.ElementQueue().Offer(nextElement)
		}

	case parse.ElementTypeLongFlagSolo, parse.ElementTypeLongFlagPair:
		if c.FlagIndex().ByLongName(nextElement.Data[0]) != nil {
			c.ElementQueue().Offer(nextElement)
		} else if err := arg.SetValue(nextElement.String()); err != nil {
			c.ElementQueue().Offer(nextElement)
		}

	default:
		panic("illegal state: unrecognized parser element type")
	}

	if !arg.WasHit() && argument.IsBoolean(arg) {
		return arg.SetValue("true")
	}

	return nil
}

func AppendUnrecognizedLongFlag(warnings *[]argo.InputWarning, name string) {
	*warnings = append(*warnings, argo.InputWarning{
		Type:    argo.UnrecognizedFlag,
		Message: fmt.Sprintf("unrecognized long flag --%s", name),
	})
}

func AppendUnrecognizedShortFlag(warnings *[]argo.InputWarning, name byte) {
	*warnings = append(*warnings, argo.InputWarning{
		Type:    argo.UnrecognizedFlag,
		Message: fmt.Sprintf("unrecognized short flag -%c", name),
	})
}
