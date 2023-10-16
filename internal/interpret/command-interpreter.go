package interpret

import (
	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/impl/argument/argerr"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/parse"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type commandInterpreter struct {
	parser  parse.Parser
	command argo.Command
}

func (c commandInterpreter) Run() error {
	var err error
	boundary := false

FOR:
	for {
		element := c.parser.Next()

		// If we've already hit the boundary marker then everything else is just a
		// passthrough.
		if boundary {
			if element.Type == parse.ElementTypeEnd {
				break
			}

			c.command.AppendPassthrough(element.String())
			continue
		}

		switch element.Type {

		case parse.ElementTypePlainText:
			if err = c.command.AppendArgument(element.String()); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockSolo:
			if boundary, err = c.interpretShortSolo(&element); err != nil {
				return err
			}

		case parse.ElementTypeShortBlockPair:
			if boundary, err = c.interpretShortPair(&element); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagSolo:
			if boundary, err = c.interpretLongSolo(&element); err != nil {
				return err
			}

		case parse.ElementTypeLongFlagPair:
			if boundary, err = c.interpretLongPair(&element); err != nil {
				return err
			}

		case parse.ElementTypeBoundary:
			boundary = true
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
				errs.AppendError(flag.MissingFlagError(f))
			}

			if f.WasHit() && f.RequiresArgument() && !f.Argument().WasHit() {
				errs.AppendError(argerr.MissingRequiredFlagArgumentError(f.Argument(), f, c.command))
			}
		}
	}

	arguments := c.command.Arguments()
	for i := range arguments {
		arg := arguments[i]

		if arg.IsRequired() && !arg.WasHit() {
			errs.AppendError(argerr.MissingRequiredPositionalArgumentError(arg, c.command))
		}
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	return nil
}

func (c commandInterpreter) interpretShortSolo(e *parse.Element) (bool, error) {
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
			c.command.AppendUnmapped(chars.StrDash + remainder[0:1])
			continue
		}

		// If the flag we found requires an argument, eat the rest of the block and
		// pass it to the flag.Hit method.  Since the block will have been consumed
		// after this, return here.
		if f.RequiresArgument() {

			// If we don't have any more characters in this short block, then we have
			// to consume the next element as the argument for this flag.
			if !h {
				nextElement := c.parser.Next()

				// If the next element is literally the end of the cli args, then we
				// obviously can't set an argument on this flag.  Tough luck, dude.
				if nextElement.Type == parse.ElementTypeEnd {
					return false, f.Hit()
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					return true, f.Hit()
				}

				// If we're here then we have a next element, and we're going to
				// sacrifice it to the flag gods.
				return false, f.HitWithArg(nextElement.String())
			}

			// So we have at least one more character in this block.  Eat that and
			// anything else as the flag argument.
			return false, f.HitWithArg(remainder[1:])
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

					// hit the current flag with an empty value
					if err := f.Hit(); err != nil {
						return false, err
					}

					// skip on to the next flag
					continue
				} else

				// Since there is no flag matching the next character, then we have to
				// assume that the remaining text is the argument for the flag
				{
					return false, f.HitWithArg(remainder[1:])
				}
			} else {

				nextElement := c.parser.Next()

				switch nextElement.Type {

				case parse.ElementTypeEnd:
					return false, f.Hit()

				case parse.ElementTypeBoundary:
					return true, f.Hit()

				case parse.ElementTypePlainText:
					return false, f.HitWithArg(nextElement.Data[0])

				case parse.ElementTypeShortBlockSolo:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						if err := f.Hit(); err != nil {
							return false, err
						}

						return c.interpretShortSolo(&nextElement)
					} else {
						return false, f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeShortBlockPair:
					if c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
						if err := f.Hit(); err != nil {
							return false, err
						}

						return c.interpretShortPair(&nextElement)
					} else {
						return false, f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeLongFlagPair:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						if err := f.Hit(); err != nil {
							return false, err
						}

						return c.interpretLongPair(&nextElement)
					} else {
						return false, f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeLongFlagSolo:
					if c.command.FindLongFlag(nextElement.Data[0]) != nil {
						if err := f.Hit(); err != nil {
							return false, err
						}

						return c.interpretLongSolo(&nextElement)
					} else {
						return false, f.HitWithArg(nextElement.String())
					}

				default:
					panic("illegal state")

				}
			}
		}

		// The flag doesn't expect an argument, just hit it
		if err := f.Hit(); err != nil {
			return false, err
		}

		// if has next
		remainder = remainder[1:]
	}

	return false, nil
}

func (c commandInterpreter) interpretShortPair(e *parse.Element) (bool, error) {
	block := e.Data[0]

	if len(block) == 0 {
		c.command.AppendUnmapped(e.String())
		return false, nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.command.FindShortFlag(block[0]); f != nil {
			return false, f.HitWithArg(e.Data[1])
		} else {
			c.command.AppendUnmapped(e.String())
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
			c.command.AppendUnmapped(chars.StrDash + block[0:1])
			continue
		}

		if f.RequiresArgument() {
			if h {
				return false, f.HitWithArg(block[1:] + "=" + e.Data[1])
			} else {
				return false, f.HitWithArg(e.Data[1])
			}
		}

		if f.HasArgument() {
			if !h {
				return false, f.HitWithArg(e.Data[1])
			}

			if c.command.FindShortFlag(block[1]) != nil {
				return false, f.Hit()
			}

			return false, f.HitWithArg(block[1:] + "=" + e.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !h {
			return false, f.HitWithArg(e.Data[1])
		}

		// Well, now that's out of the way, we can move on to the next flag (after
		// we mark this one as hit of course).
		if err := f.Hit(); err != nil {
			return false, err
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c commandInterpreter) interpretLongSolo(e *parse.Element) (bool, error) {
	f := c.command.FindLongFlag(e.Data[0])

	if f == nil {
		c.command.AppendUnmapped(e.String())
		return false, nil
	}

	if f.RequiresArgument() {
		nextElement := c.parser.Next()

		if nextElement.Type == parse.ElementTypeEnd {
			return false, f.Hit()
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			return true, f.Hit()
		}

		return false, f.HitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.parser.Next()

		switch nextElement.Type {

		case parse.ElementTypeEnd:
			return false, f.Hit()

		case parse.ElementTypeBoundary:
			return true, f.Hit()

		case parse.ElementTypePlainText:
			return false, f.HitWithArg(nextElement.Data[0])

		case parse.ElementTypeLongFlagSolo:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				if err := f.Hit(); err != nil {
					return false, err
				}

				return c.interpretLongSolo(&nextElement)
			} else {
				return false, f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeLongFlagPair:
			if c.command.FindLongFlag(nextElement.Data[0]) != nil {
				if err := f.Hit(); err != nil {
					return false, err
				}

				return c.interpretLongPair(&nextElement)
			} else {
				return false, f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				if err := f.Hit(); err != nil {
					return false, err
				}

				return c.interpretShortSolo(&nextElement)
			} else {
				return false, f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.command.FindShortFlag(nextElement.Data[0][0]) != nil {
				if err := f.Hit(); err != nil {
					return false, err
				}

				return c.interpretShortPair(&nextElement)
			} else {
				return false, f.HitWithArg(nextElement.String())
			}

		default:
			panic("illegal state")
		}
	}

	return false, f.Hit()
}

func (c commandInterpreter) interpretLongPair(e *parse.Element) (bool, error) {
	if found, err := c.command.TryFlag(flag.NewLongFlagRef(e.Data[0], e.Data[1], true)); err != nil {
		return false, err
	} else if !found {
		c.command.AppendUnmapped(e.String())
	}

	return false, nil
}
