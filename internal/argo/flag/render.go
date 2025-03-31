package flag

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Render(flag argo.Flag, options opts.Options, padding uint8, sb *utils.BatchWriter) {
	sb.WriteString(render.HeaderPadding[padding])

	// If the flag has a long form name
	if flag.HasLongForm() {

		// AND a short form character
		if flag.HasShortForm() {
			renderShortForm(sb, flag, text.SpaceByte)

			sb.WriteString(render.FlagDivider)
		}

		renderLongForm(sb, flag, text.EqualsByte)
	} else {
		renderShortForm(sb, flag, text.SpaceByte)
	}

	if flag.HasDescription() {
		sb.WriteByte(text.LineFeedByte)

		render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb).
			Format(flag.Description())
	}

	if flag.HasArgument() && flag.Argument().HasDescription() {
		sb.WriteByte(text.LineFeedByte)
		RenderArgument(flag.Argument(), options, padding+1, sb)
	}
}

func RenderArgument(arg argo.Argument, options opts.Options, padding uint8, out *utils.BatchWriter) {
	out.WriteString(render.SubLinePadding[padding])

	argument.RenderName(arg, out, 0)

	if arg.HasDescription() {
		out.WriteByte(text.LineFeedByte)

		render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), out).
			Format(arg.Description())
	}
}

// RenderShortestForUsage renders the shortest form of the given flag for the
// command usage line.
//
// If the flag has a short form, that form will be rendered, otherwise the long
// form will be rendered.
//
// If the flag has an argument, that argument may also be rendered.
func RenderShortestForUsage(flag argo.Flag, sb *utils.BatchWriter) {
	if flag.HasShortForm() {
		renderShortForm(sb, flag, text.EqualsByte)
	} else {
		renderLongForm(sb, flag, text.EqualsByte)
	}
}

const (
	fgSingleName = "Flags"
)

func RenderGroups(groups []argo.FlagGroup, options opts.Options, padding uint8, out *utils.BatchWriter) {
	for i, group := range groups {
		if i > 0 {
			out.WriteString(render.ParagraphBreak)
		}

		RenderGroup(group, options, padding, out, len(groups) > 1)
	}
}

func RenderGroup(
	group argo.FlagGroup,
	options opts.Options,
	padding uint8,
	out *utils.BatchWriter,
	multiple bool,
) {
	out.WriteString(render.HeaderPadding[padding])

	if group.Name() == DefaultFlagGroupName {
		if multiple {
			out.WriteString(options.MetaFlagGroupName())
		} else {
			out.WriteString(fgSingleName)
		}
	} else {
		out.WriteString(group.Name())
	}

	// If the group has a description, print it out.
	if group.HasDescription() {
		out.WriteByte(text.LineFeedByte)

		render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), out).
			Format(group.Description())

		out.WriteByte(text.LineFeedByte)
	}

	// Render every flag in the group.
	for i, flag := range group.Flags() {
		if i > 0 {
			if !group.Flags()[i-1].HasDescription() {
				out.WriteByte(text.LineFeedByte)
			}
		}
		out.WriteByte(text.LineFeedByte)

		Render(flag, options, padding+1, out)
	}
}

func renderShortForm(w *utils.BatchWriter, f argo.Flag, divider byte) {
	w.WriteByte(text.DashByte)
	w.WriteByte(f.ShortForm())
	tryRenderArgument(w, f, divider)
}

func renderLongForm(w *utils.BatchWriter, f argo.Flag, divider byte) {
	w.WriteString(text.DoubleDash)
	w.WriteString(f.LongForm())
	tryRenderArgument(w, f, divider)
}

func tryRenderArgument(w *utils.BatchWriter, f argo.Flag, d byte) {
	if f.HasArgument() && argument.ShouldBeRendered(f.Argument()) {
		w.WriteByte(d)
		argument.RenderName(f.Argument(), w, 0)
	}
}
