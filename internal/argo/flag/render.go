package flag

import (
	"bufio"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Render(flag argo.Flag, options argo.Options, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if flag.HasLongForm() {

		// AND a short form character
		if flag.HasShortForm() {
			if err := renderShortForm(sb, flag); err != nil {
				return err
			}

			if _, err := sb.WriteString(chars.FlagDivider); err != nil {
				return err
			}
		}

		if err := renderLongForm(sb, flag); err != nil {
			return err
		}
	} else {
		if err := renderShortForm(sb, flag); err != nil {
			return err
		}
	}

	if flag.HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], options.HelpTextMaxWidth, sb)
		if err := formatter.Format(flag.Description()); err != nil {
			return err
		}
	}

	if flag.HasArgument() && flag.Argument().HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := RenderArgument(flag.Argument(), options, padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}

func RenderArgument(arg argo.Argument, options argo.Options, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(chars.SubLinePadding[padding]); err != nil {
		return err
	}

	if err := argument.RenderName(arg, out, 0); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], options.HelpTextMaxWidth, out)
		if err := formatter.Format(arg.Description()); err != nil {
			return err
		}
	}

	return nil
}

// RenderShortestLine renders the shortest form of the given flag.
//
// If the flag has a short form, that form will be rendered, otherwise the long
// form will be rendered.
//
// If the flag has an argument, that argument may also be rendered.
func RenderShortestLine(flag argo.Flag, sb *bufio.Writer) error {
	if flag.HasShortForm() {
		return renderShortForm(sb, flag)
	}

	return renderLongForm(sb, flag)
}

const (
	fgDefaultName = "General Flags"
	fgSingleName  = "Flags"
)

func RenderGroups(groups []argo.FlagGroup, options argo.Options, padding uint8, out *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := out.WriteString(chars.ParagraphBreak); err != nil {
				return err
			}
		}

		if err := RenderGroup(group, options, padding, out, len(groups) > 1); err != nil {
			return err
		}
	}

	return nil
}

func RenderGroup(
	group argo.FlagGroup,
	options argo.Options,
	padding uint8,
	out *bufio.Writer,
	multiple bool,
) error {
	if _, err := out.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	if group.Name() == chars.DefaultGroupName {
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
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], options.HelpTextMaxWidth, out)
		if err := formatter.Format(group.Description()); err != nil {
			return err
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}
	}

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			if !group.Flags()[i-1].HasDescription() {
				if err := out.WriteByte(chars.CharLF); err != nil {
					return err
				}
			}
		}
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		if err := Render(flag, options, padding+1, out); err != nil {
			return err
		}
	}

	return nil
}

type Forms struct {
	Short bool
	Long  bool
	Flag  argo.Flag
}

func FlattenInheritance(node GroupContainer) []Forms {
	shorts := make(map[byte]bool, 16)
	longs := make(map[string]bool, 16)
	options := make([]Forms, 0, 16)

	current := node
LOOP:
	for {
		for _, group := range current.FlagGroups() {
			for _, flag := range group.Flags() {
				forms := Forms{Flag: flag}

				if flag.HasLongForm() {
					if _, ok := longs[flag.LongForm()]; !ok {
						longs[flag.LongForm()] = true
						forms.Long = true
					}
				}

				if flag.HasShortForm() {
					if _, ok := shorts[flag.ShortForm()]; !ok {
						shorts[flag.ShortForm()] = true
						forms.Short = true
					}
				}

				if (forms.Short || forms.Long) && current != node {
					options = append(options, forms)
				}
			}
		}

		switch t := current.(type) {
		case argo.ChildNode:
			current = t.Parent().(GroupContainer)
		default:
			break LOOP
		}
	}

	return options
}

func RenderInherited(forms *Forms, options argo.Options, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if forms.Long {

		// AND a short form character
		if forms.Short {
			if err := sb.WriteByte(chars.CharDash); err != nil {
				return err
			}
			if err := sb.WriteByte(forms.Flag.ShortForm()); err != nil {
				return err
			}

			if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
				if err := sb.WriteByte(chars.CharSpace); err != nil {
					return err
				}
				if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
					return err
				}
			}

			if _, err := sb.WriteString(chars.FlagDivider); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(chars.StrDoubleDash); err != nil {
			return err
		}
		if _, err := sb.WriteString(forms.Flag.LongForm()); err != nil {
			return err
		}

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			if err := sb.WriteByte(chars.CharEquals); err != nil {
				return err
			}
			if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if err := sb.WriteByte(chars.CharDash); err != nil {
			return err
		}
		if err := sb.WriteByte(forms.Flag.ShortForm()); err != nil {
			return err
		}

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			if err := sb.WriteByte(chars.CharSpace); err != nil {
				return err
			}
			if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}

	if forms.Flag.HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], options.HelpTextMaxWidth, sb)
		if err := formatter.Format(forms.Flag.Description()); err != nil {
			return err
		}
	}

	if forms.Flag.HasArgument() && forms.Flag.Argument().HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := RenderArgument(forms.Flag.Argument(), options, padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}

func renderShortForm(w *bufio.Writer, f argo.Flag) error {
	if err := w.WriteByte(chars.CharDash); err != nil {
		return err
	}

	if err := w.WriteByte(f.ShortForm()); err != nil {
		return err
	}

	return tryRenderArgument(w, f, chars.CharSpace)
}

func renderLongForm(w *bufio.Writer, f argo.Flag) error {
	if _, err := w.WriteString(chars.StrDoubleDash); err != nil {
		return err
	}

	if _, err := w.WriteString(f.LongForm()); err != nil {
		return err
	}

	return tryRenderArgument(w, f, chars.CharEquals)
}

func tryRenderArgument(w *bufio.Writer, f argo.Flag, d byte) error {
	if f.HasArgument() && argument.ShouldBeRendered(f.Argument()) {
		if err := w.WriteByte(d); err != nil {
			return err
		}

		if err := argument.RenderName(f.Argument(), w, 0); err != nil {
			return err
		}
	}

	return nil
}
