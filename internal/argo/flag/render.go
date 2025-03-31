package flag

import (
	"bufio"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Render(flag argo.Flag, options opts.Options, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(render.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if flag.HasLongForm() {

		// AND a short form character
		if flag.HasShortForm() {
			if err := renderShortForm(sb, flag, text.SpaceByte); err != nil {
				return err
			}

			if _, err := sb.WriteString(render.FlagDivider); err != nil {
				return err
			}
		}

		if err := renderLongForm(sb, flag, text.EqualsByte); err != nil {
			return err
		}
	} else {
		if err := renderShortForm(sb, flag, text.SpaceByte); err != nil {
			return err
		}
	}

	if flag.HasDescription() {
		if err := sb.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb)
		if err := formatter.Format(flag.Description()); err != nil {
			return err
		}
	}

	if flag.HasArgument() && flag.Argument().HasDescription() {
		if err := sb.WriteByte(text.LineFeedByte); err != nil {
			return err
		}
		if err := RenderArgument(flag.Argument(), options, padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}

func RenderArgument(arg argo.Argument, options opts.Options, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(render.SubLinePadding[padding]); err != nil {
		return err
	}

	if err := argument.RenderName(arg, out, 0); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), out)
		if err := formatter.Format(arg.Description()); err != nil {
			return err
		}
	}

	return nil
}

// RenderShortestForUsage renders the shortest form of the given flag for the
// command usage line.
//
// If the flag has a short form, that form will be rendered, otherwise the long
// form will be rendered.
//
// If the flag has an argument, that argument may also be rendered.
func RenderShortestForUsage(flag argo.Flag, sb *bufio.Writer) error {
	if flag.HasShortForm() {
		return renderShortForm(sb, flag, text.EqualsByte)
	}

	return renderLongForm(sb, flag, text.EqualsByte)
}

const (
	fgSingleName = "Flags"
)

func RenderGroups(groups []argo.FlagGroup, options opts.Options, padding uint8, out *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := out.WriteString(render.ParagraphBreak); err != nil {
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
	options opts.Options,
	padding uint8,
	out *bufio.Writer,
	multiple bool,
) error {
	if _, err := out.WriteString(render.HeaderPadding[padding]); err != nil {
		return err
	}

	if group.Name() == DefaultFlagGroupName {
		if multiple {
			if _, err := out.WriteString(options.MetaFlagGroupName()); err != nil {
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
		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), out)
		if err := formatter.Format(group.Description()); err != nil {
			return err
		}
		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}
	}

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			if !group.Flags()[i-1].HasDescription() {
				if err := out.WriteByte(text.LineFeedByte); err != nil {
					return err
				}
			}
		}
		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		if err := Render(flag, options, padding+1, out); err != nil {
			return err
		}
	}

	return nil
}

func renderShortForm(w *bufio.Writer, f argo.Flag, divider byte) error {
	if err := w.WriteByte(text.DashByte); err != nil {
		return err
	}

	if err := w.WriteByte(f.ShortForm()); err != nil {
		return err
	}

	return tryRenderArgument(w, f, divider)
}

func renderLongForm(w *bufio.Writer, f argo.Flag, divider byte) error {
	if _, err := w.WriteString(text.DoubleDash); err != nil {
		return err
	}

	if _, err := w.WriteString(f.LongForm()); err != nil {
		return err
	}

	return tryRenderArgument(w, f, divider)
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
