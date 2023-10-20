package argo

import (
	"bufio"
	"reflect"
	"strconv"
)

const (
	argReqPrefix = '<'
	argReqSuffix = '>'
	argOptPrefix = '['
	argOptSuffix = ']'
)

type renderBase struct{}

func (r renderBase) renderArgName(a Argument, argIndex int) string {
	if a.HasName() {
		return a.Name()
	} else {
		if argIndex > 0 {
			return "arg" + strconv.Itoa(argIndex)
		} else {
			return "arg"
		}
	}
}

func (r renderBase) renderFlagArgument(arg Argument, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(subLinePadding[padding]); err != nil {
		return err
	}
	if err := r.renderArgumentName(arg, out, 0); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out); err != nil {
			return err
		}
	}
	return nil
}

func (r renderBase) renderArgument(arg Argument, padding uint8, out *bufio.Writer, argIndex int) error {
	if _, err := out.WriteString(headerPadding[padding]); err != nil {
		return err
	}
	if err := r.renderArgumentName(arg, out, argIndex); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(arg.Description(), descriptionPadding[padding], helpTextMaxWidth, out); err != nil {
			return err
		}
	}

	return nil
}

func (r renderBase) renderArgumentName(a Argument, out *bufio.Writer, argIndex int) error {
	if a.HasBinding() && a.BindingType().Kind() == reflect.Bool {
		return nil
	}

	if a.IsRequired() {
		if err := out.WriteByte(argReqPrefix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(argOptPrefix); err != nil {
			return err
		}
	}

	if _, err := out.WriteString(r.renderArgName(a, argIndex)); err != nil {
		return err
	}

	if a.IsRequired() {
		if err := out.WriteByte(argReqSuffix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(argOptSuffix); err != nil {
			return err
		}
	}

	return nil
}

func (r renderBase) renderFlag(flag Flag, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(headerPadding[padding]); err != nil {
		return err
	}

	if flag.HasLongForm() {
		if flag.HasShortForm() {
			if err := sb.WriteByte(charDash); err != nil {
				return err
			}
			if err := sb.WriteByte(flag.ShortForm()); err != nil {
				return err
			}

			if flag.HasArgument() {
				if err := sb.WriteByte(charSpace); err != nil {
					return err
				}
				if err := r.renderArgumentName(flag.Argument(), sb, 0); err != nil {
					return err
				}
			}

			if _, err := sb.WriteString(flagDivider); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(strDoubleDash); err != nil {
			return err
		}
		if _, err := sb.WriteString(flag.LongForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(charEquals); err != nil {
				return err
			}
			if err := r.renderArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if err := sb.WriteByte(charDash); err != nil {
			return err
		}
		if err := sb.WriteByte(flag.ShortForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(charSpace); err != nil {
				return err
			}
			if err := r.renderArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}

	if flag.HasDescription() {
		if err := sb.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(flag.Description(), descriptionPadding[padding], helpTextMaxWidth, sb); err != nil {
			return err
		}
	}

	if flag.HasArgument() && flag.Argument().HasDescription() {
		if err := sb.WriteByte(charLF); err != nil {
			return err
		}
		if err := r.renderFlagArgument(flag.Argument(), padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}

func (r renderBase) renderShortestFlagLine(flag Flag, sb *bufio.Writer) error {
	if flag.HasShortForm() {
		if err := sb.WriteByte(charDash); err != nil {
			return err
		}
		if err := sb.WriteByte(flag.ShortForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(charSpace); err != nil {
				return err
			}
			if err := r.renderArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if _, err := sb.WriteString(strDoubleDash); err != nil {
			return err
		}
		if _, err := sb.WriteString(flag.LongForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(charEquals); err != nil {
				return err
			}
			if err := r.renderArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

const (
	fgDefaultName = "General Flags"
	fgSingleName  = "Flags"
)

func (r renderBase) renderFlagGroups(groups []FlagGroup, padding uint8, out *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := out.WriteString(paragraphBreak); err != nil {
				return err
			}
		}

		if err := r.renderFlagGroup(group, padding, out, len(groups) > 1); err != nil {
			return err
		}
	}

	return nil
}

func (r renderBase) renderFlagGroup(
	group FlagGroup,
	padding uint8,
	out *bufio.Writer,
	multiple bool,
) error {
	if _, err := out.WriteString(headerPadding[padding]); err != nil {
		return err
	}

	if group.Name() == defaultGroupName {
		if multiple {
			if _, err := out.WriteString(fgDefaultName); err != nil {
				return err
			}
		} else {
			if _, err := out.WriteString(fgSingleName); err != nil {
				return err
			}
		}
	} else {
		if _, err := out.WriteString(group.Name()); err != nil {
			return err
		}
	}

	// If the group has a description, print it out.
	if group.HasDescription() {
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
		if err := breakFmt(group.Description(), descriptionPadding[padding], helpTextMaxWidth, out); err != nil {
			return err
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}
	}

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			if !group.Flags()[i-1].HasDescription() {
				if err := out.WriteByte(charLF); err != nil {
					return err
				}
			}
		}
		if err := out.WriteByte(charLF); err != nil {
			return err
		}

		if err := r.renderFlag(flag, padding+1, out); err != nil {
			return err
		}
	}

	return nil
}
