package interpret

import (
	"fmt"

	"github.com/Foxcapades/Argonaut/v1/internal/chars"
	"github.com/Foxcapades/Argonaut/v1/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/v1/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/v1/internal/parse"
	"github.com/Foxcapades/Argonaut/v1/internal/xerr"
	"github.com/Foxcapades/Argonaut/v1/pkg/argo"
)

type commandTreeInterpreter struct {
	parser   parse.Parser
	current  argo.CommandNode
	boundary bool
}

func (c commandTreeInterpreter) Run() error {
	passthroughs := make([]string, 0, 10)
	unmapped := make([]string, 0, 10)

	for {
		element := c.parser.Next()

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
				if err := node.AppendArgument(element.String()); err != nil {
					return err
				}
			} else if node, ok := c.current.(argo.CommandParent); ok {
				if child := node.FindChild(element.String()); child != nil {
					c.current = child
				} else {
					// TODO: wrap this error in an "invalid subcommand" error or some such.
					return fmt.Errorf("unrecognized subcommand %s", element.String())
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
			break

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	errs := xerr.NewMultiError()

	if node, ok := c.current.(argo.CommandLeaf); !ok {
		// TODO: wrap this error in an error type that indicates that the command
		//       leaf was not hit.
		return fmt.Errorf("command leaf was not reached")
	} else {
		for _, value := range unmapped {
			node.AppendUnmapped(value)
		}

		for _, value := range passthroughs {
			node.AppendPassthrough(value)
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
					errs.AppendError(fmt.Errorf("required flag %s was not used", flagutil.PrintFlagNames(f)))
				} else if f.RequiresArgument() && !f.Argument().WasHit() {
					errs.AppendError(fmt.Errorf("flag %s requires an argument", flagutil.PrintFlagNames(f)))
				}
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return errs
	}

	return nil
}

func (c commandTreeInterpreter) interpretShortSolo(element *parse.Element, unmapped *[]string) error {
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
			*unmapped = append(*unmapped, chars.StrDash+remainder[0:1])
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
					return f.Hit()
				}

				if nextElement.Type == parse.ElementTypeBoundary {
					c.boundary = true
					return f.Hit()
				}

				// If we're here then we have a next element, and we're going to
				// sacrifice it to the flag gods.
				return f.HitWithArg(nextElement.String())
			}

			// So we have at least one more character in this block.  Eat that and
			// anything else as the flag argument.
			return f.HitWithArg(remainder[1:])
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
					if err := f.Hit(); err != nil {
						return err
					}

					// skip on to the next flag
					continue
				} else

				// Since there is no flag matching the next character, then we have to
				// assume that the remaining text is the argument for the flag
				{
					return f.HitWithArg(remainder[1:])
				}
			} else {

				nextElement := c.parser.Next()

				switch nextElement.Type {

				case parse.ElementTypeEnd:
					return f.Hit()

				case parse.ElementTypeBoundary:
					c.boundary = true
					return f.Hit()

				case parse.ElementTypePlainText:
					return f.HitWithArg(nextElement.Data[0])

				case parse.ElementTypeShortBlockSolo:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						if err := f.Hit(); err != nil {
							return err
						}

						return c.interpretShortSolo(&nextElement, unmapped)
					} else {
						return f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeShortBlockPair:
					if c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
						if err := f.Hit(); err != nil {
							return err
						}

						return c.interpretShortPair(&nextElement, unmapped)
					} else {
						return f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeLongFlagPair:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						if err := f.Hit(); err != nil {
							return err
						}

						return c.interpretLongPair(&nextElement, unmapped)
					} else {
						return f.HitWithArg(nextElement.String())
					}

				case parse.ElementTypeLongFlagSolo:
					if c.current.FindLongFlag(nextElement.Data[0]) != nil {
						if err := f.Hit(); err != nil {
							return err
						}

						return c.interpretLongSolo(&nextElement, unmapped)
					} else {
						return f.HitWithArg(nextElement.String())
					}

				default:
					panic("illegal state: unrecognized parser element type")

				}
			}
		}

		// The flag doesn't expect an argument, just hit it
		if err := f.Hit(); err != nil {
			return err
		}

		// if has next
		remainder = remainder[1:]
	}

	return nil
}

func (c commandTreeInterpreter) interpretShortPair(element *parse.Element, unmapped *[]string) error {
	block := element.Data[0]

	if len(block) == 0 {
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	// If the flag key block is a single character in length, then we can do this
	// in a simple check.
	if len(block) == 1 {
		if f := c.current.FindShortFlag(block[0]); f != nil {
			return f.HitWithArg(element.Data[1])
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
			*unmapped = append(*unmapped, chars.StrDash+block[0:1])
			continue
		}

		if f.RequiresArgument() {
			if h {
				return f.HitWithArg(block[1:] + "=" + element.Data[1])
			} else {
				return f.HitWithArg(element.Data[1])
			}
		}

		if f.HasArgument() {
			if !h {
				return f.HitWithArg(element.Data[1])
			}

			if c.current.FindShortFlag(block[1]) != nil {
				return f.Hit()
			}

			return f.HitWithArg(block[1:] + "=" + element.Data[1])
		}

		// So the flag doesn't expect an argument at all.
		// Well let's see what we have to say about that.  It may be, if this is the
		// last character in the block, that it has to have one anyway.
		if !h {
			return f.HitWithArg(element.Data[1])
		}

		// Well, now that's out of the way, we can move on to the next flag (after
		// we mark this one as hit of course).
		if err := f.Hit(); err != nil {
			return err
		}

		block = block[1:]
	}

	panic("illegal state")
}

func (c commandTreeInterpreter) interpretLongSolo(element *parse.Element, unmapped *[]string) error {
	f := c.current.FindLongFlag(element.Data[0])

	if f == nil {
		*unmapped = append(*unmapped, element.String())
		return nil
	}

	if f.RequiresArgument() {
		nextElement := c.parser.Next()

		if nextElement.Type == parse.ElementTypeEnd {
			return f.Hit()
		}

		if nextElement.Type == parse.ElementTypeBoundary {
			c.boundary = true
			return f.Hit()
		}

		return f.HitWithArg(nextElement.String())
	}

	if f.HasArgument() {
		nextElement := c.parser.Next()

		switch nextElement.Type {

		case parse.ElementTypeEnd:
			return f.Hit()

		case parse.ElementTypeBoundary:
			c.boundary = true
			return f.Hit()

		case parse.ElementTypePlainText:
			return f.HitWithArg(nextElement.Data[0])

		case parse.ElementTypeLongFlagSolo:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				if err := f.Hit(); err != nil {
					return err
				}

				return c.interpretLongSolo(&nextElement, unmapped)
			} else {
				return f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeLongFlagPair:
			if c.current.FindLongFlag(nextElement.Data[0]) != nil {
				if err := f.Hit(); err != nil {
					return err
				}

				return c.interpretLongPair(&nextElement, unmapped)
			} else {
				return f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeShortBlockSolo:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				if err := f.Hit(); err != nil {
					return err
				}

				return c.interpretShortSolo(&nextElement, unmapped)
			} else {
				return f.HitWithArg(nextElement.String())
			}

		case parse.ElementTypeShortBlockPair:
			if len(nextElement.Data[0]) > 0 && c.current.FindShortFlag(nextElement.Data[0][0]) != nil {
				if err := f.Hit(); err != nil {
					return err
				}

				return c.interpretShortPair(&nextElement, unmapped)
			} else {
				return f.HitWithArg(nextElement.String())
			}

		default:
			panic("illegal state: unrecognized parser element type")
		}
	}

	return f.Hit()
}

func (c commandTreeInterpreter) interpretLongPair(element *parse.Element, unmapped *[]string) error {
	if found, err := c.current.TryFlag(flag.NewLongFlagRef(element.Data[0], element.Data[1], true)); err != nil {
		return err
	} else if !found {
		*unmapped = append(*unmapped, element.String())
	}

	return nil
}
