package flag

import (
	"bufio"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
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

func RenderInheritedForms(forms *Forms, options opts.Options, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(render.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if forms.Long {

		// AND a short form character
		if forms.Short {
			if err := sb.WriteByte(text.DashByte); err != nil {
				return err
			}
			if err := sb.WriteByte(forms.Flag.ShortForm()); err != nil {
				return err
			}

			if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
				if err := sb.WriteByte(text.SpaceByte); err != nil {
					return err
				}
				if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
					return err
				}
			}

			if _, err := sb.WriteString(render.FlagDivider); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(text.DoubleDash); err != nil {
			return err
		}

		if _, err := sb.WriteString(forms.Flag.LongForm()); err != nil {
			return err
		}

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			if err := sb.WriteByte(text.EqualsByte); err != nil {
				return err
			}

			if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if err := sb.WriteByte(text.DashByte); err != nil {
			return err
		}

		if err := sb.WriteByte(forms.Flag.ShortForm()); err != nil {
			return err
		}

		if forms.Flag.HasArgument() && argument.ShouldBeRendered(forms.Flag.Argument()) {
			if err := sb.WriteByte(text.SpaceByte); err != nil {
				return err
			}

			if err := argument.RenderName(forms.Flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}

	if forms.Flag.HasDescription() {
		if err := sb.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb)
		if err := formatter.Format(forms.Flag.Description()); err != nil {
			return err
		}
	}

	if forms.Flag.HasArgument() && forms.Flag.Argument().HasDescription() {
		if err := sb.WriteByte(text.LineFeedByte); err != nil {
			return err
		}
		if err := RenderArgument(forms.Flag.Argument(), options, padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}
