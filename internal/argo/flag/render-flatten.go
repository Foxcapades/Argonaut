package flag

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

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
		for _, flag := range Stream(current) {
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

		switch t := current.(type) {
		case argo.ChildNode:
			current = t.Parent().(GroupContainer)
		default:
			break LOOP
		}
	}

	return options
}

func RenderInheritedForms(forms *Forms, options opts.Options, padding uint8, sb *utils.BatchWriter) {
	sb.WriteString(render.HeaderPadding[padding])

	// If the flag has a long form name
	if forms.Long {

		// AND a short form character
		if forms.Short {
			sb.WriteByte(text.DashByte)
			sb.WriteByte(forms.Flag.ShortForm())

			if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
				sb.WriteByte(text.SpaceByte)
				argument.RenderName(forms.Flag.Argument(), sb, 0)
			}

			sb.WriteString(render.FlagDivider)
		}

		sb.WriteString(text.DoubleDash)

		sb.WriteString(forms.Flag.LongForm())

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			sb.WriteByte(text.EqualsByte)

			argument.RenderName(forms.Flag.Argument(), sb, 0)
		}
	} else {
		sb.WriteByte(text.DashByte)

		sb.WriteByte(forms.Flag.ShortForm())

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			sb.WriteByte(text.SpaceByte)

			argument.RenderName(forms.Flag.Argument(), sb, 0)
		}
	}

	if forms.Flag.HasDescription() {
		sb.WriteByte(text.LineFeedByte)

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb)
		formatter.Format(forms.Flag.Description())
	}

	if forms.Flag.HasArgument() && forms.Flag.Argument().HasDescription() {
		sb.WriteByte(text.LineFeedByte)
		RenderArgument(forms.Flag.Argument(), options, padding+1, sb)
	}
}
