package render

import (
	"bufio"
	"github.com/Foxcapades/Argonaut/v3/internal/chars"
	"github.com/Foxcapades/Argonaut/v3/pkg/argo/argument"
	cli_tree "github.com/Foxcapades/Argonaut/v3/pkg/argo/command/tree"
	"github.com/Foxcapades/Argonaut/v3/pkg/argo/flag"
	"reflect"
)

func FlagArgShouldBeRendered(arg cli_arg.Argument) bool {
	return arg.IsRequired() ||
		!arg.HasBinding() ||
		arg.BindingType().Kind() != reflect.Bool
}

func Flag(flag cli_flag.Flag, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if flag.HasLongForm() {

		// AND a short form character
		if flag.HasShortForm() {
			if err := sb.WriteByte(chars.CharDash); err != nil {
				return err
			}
			if err := sb.WriteByte(flag.ShortForm()); err != nil {
				return err
			}

			if flag.HasArgument() && FlagArgShouldBeRendered(flag.Argument()) {
				if err := sb.WriteByte(chars.CharSpace); err != nil {
					return err
				}
				if err := ArgumentName(flag.Argument(), sb, 0); err != nil {
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
		if _, err := sb.WriteString(flag.LongForm()); err != nil {
			return err
		}

		if flag.HasArgument() && FlagArgShouldBeRendered(flag.Argument()) {
			if err := sb.WriteByte(chars.CharEquals); err != nil {
				return err
			}
			if err := ArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if err := sb.WriteByte(chars.CharDash); err != nil {
			return err
		}
		if err := sb.WriteByte(flag.ShortForm()); err != nil {
			return err
		}

		if flag.HasArgument() && FlagArgShouldBeRendered(flag.Argument()) {
			if err := sb.WriteByte(chars.CharSpace); err != nil {
				return err
			}
			if err := ArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}

	if flag.HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, sb)
		if err := formatter.Format(flag.Description()); err != nil {
			return err
		}
	}

	if flag.HasArgument() && flag.Argument().HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := FlagArgument(flag.Argument(), padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}

func ShortestFlagLine(flag cli_flag.Flag, sb *bufio.Writer) error {
	if flag.HasShortForm() {
		if err := sb.WriteByte(chars.CharDash); err != nil {
			return err
		}
		if err := sb.WriteByte(flag.ShortForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(chars.CharEquals); err != nil {
				return err
			}
			if err := ArgumentName(flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if _, err := sb.WriteString(chars.StrDoubleDash); err != nil {
			return err
		}
		if _, err := sb.WriteString(flag.LongForm()); err != nil {
			return err
		}

		if flag.HasArgument() {
			if err := sb.WriteByte(chars.CharEquals); err != nil {
				return err
			}
			if err := ArgumentName(flag.Argument(), sb, 0); err != nil {
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

func FlagGroups(groups []cli_flag.FlagGroup, padding uint8, out *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := out.WriteString(chars.ParagraphBreak); err != nil {
				return err
			}
		}

		if err := FlagGroup(group, padding, out, len(groups) > 1); err != nil {
			return err
		}
	}

	return nil
}

func FlagGroup(
	group cli_flag.FlagGroup,
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

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, out)
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

		if err := Flag(flag, padding+1, out); err != nil {
			return err
		}
	}

	return nil
}

type flagForms struct {
	short bool
	long  bool
	flag  cli_flag.Flag
}

func flattenFlagInheritance(node cli_tree.CommandNode) []flagForms {
	shorts := make(map[byte]bool, 16)
	longs := make(map[string]bool, 16)
	options := make([]flagForms, 0, 16)

	current := node
	for current != nil {
		for _, group := range current.FlagGroups() {
			for _, flag := range group.Flags() {
				forms := flagForms{flag: flag}

				if flag.HasLongForm() {
					if _, ok := longs[flag.LongForm()]; !ok {
						longs[flag.LongForm()] = true
						forms.long = true
					}
				}

				if flag.HasShortForm() {
					if _, ok := shorts[flag.ShortForm()]; !ok {
						shorts[flag.ShortForm()] = true
						forms.short = true
					}
				}

				if (forms.short || forms.long) && current != node {
					options = append(options, forms)
				}
			}
		}

		current = current.Parent()
	}

	return options
}

func renderInheritedFlag(forms *flagForms, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	// If the flag has a long form name
	if forms.long {

		// AND a short form character
		if forms.short {
			if err := sb.WriteByte(chars.CharDash); err != nil {
				return err
			}
			if err := sb.WriteByte(forms.flag.ShortForm()); err != nil {
				return err
			}

			if forms.flag.HasArgument() && FlagArgShouldBeRendered(forms.flag.Argument()) {
				if err := sb.WriteByte(chars.CharSpace); err != nil {
					return err
				}
				if err := ArgumentName(forms.flag.Argument(), sb, 0); err != nil {
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
		if _, err := sb.WriteString(forms.flag.LongForm()); err != nil {
			return err
		}

		if forms.flag.HasArgument() && FlagArgShouldBeRendered(forms.flag.Argument()) {
			if err := sb.WriteByte(chars.CharEquals); err != nil {
				return err
			}
			if err := ArgumentName(forms.flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	} else {
		if err := sb.WriteByte(chars.CharDash); err != nil {
			return err
		}
		if err := sb.WriteByte(forms.flag.ShortForm()); err != nil {
			return err
		}

		if forms.flag.HasArgument() && FlagArgShouldBeRendered(forms.flag.Argument()) {
			if err := sb.WriteByte(chars.CharSpace); err != nil {
				return err
			}
			if err := ArgumentName(forms.flag.Argument(), sb, 0); err != nil {
				return err
			}
		}
	}

	if forms.flag.HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, sb)
		if err := formatter.Format(forms.flag.Description()); err != nil {
			return err
		}
	}

	if forms.flag.HasArgument() && forms.flag.Argument().HasDescription() {
		if err := sb.WriteByte(chars.CharLF); err != nil {
			return err
		}
		if err := FlagArgument(forms.flag.Argument(), padding+1, sb); err != nil {
			return err
		}
	}

	return nil
}
