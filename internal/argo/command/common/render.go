package common

import (
	"bufio"
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const (
	CommandRenderPrefix        = "Usage:\n"
	CommandRenderOptionalFlags = " " + string(argument.OptPrefix) + "options" + string(argument.OptSuffix)
	CommandRenderArgs          = "Arguments"
)

type HelpRenderer[T any] = func(T, *bufio.Writer) error

type CallEndNode interface {
	flag.GroupContainer
	HasDescription() bool
	Description() string
	HasArguments() bool
	Arguments() []argo.Argument
}

func RenderCommandBackHalf(com CallEndNode, options argo.Options, out *bufio.Writer) error {
	// If the command has a description, append it.
	if com.HasDescription() {
		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[0], options.HelpTextMaxWidth, out)
		if err := formatter.Format(com.Description()); err != nil {
			return err
		}
	}

	// Figure out if we have any printable arguments.
	//
	// Showing the arguments is conditional based on whether any of the arguments
	// have a description value.  If none do, then there is no value in rendering
	// them as they already appear in the usage line.
	//
	// This is calculated ahead of time as it informs whether the flag group
	// headers should be printed.
	writeArgs := false
	if com.HasArguments() {
		for _, arg := range com.Arguments() {
			if arg.HasDescription() {
				writeArgs = true
			}
		}
	}

	if com.HasFlagGroups() {
		if com.HasDescription() {
			if err := out.WriteByte(text.LineFeedByte); err != nil {
				return err
			}
		}

		if err := out.WriteByte(text.LineFeedByte); err != nil {
			return err
		}
		if err := flag.RenderGroups(com.FlagGroups(), options, 0, out); err != nil {
			return err
		}
	}

	if writeArgs {
		if _, err := out.WriteString(render.ParagraphBreak); err != nil {
			return err
		}
		if _, err := out.WriteString(render.HeaderPadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString(CommandRenderArgs); err != nil {
			return err
		}

		multiArgs := len(com.Arguments()) > 1

		for i, arg := range com.Arguments() {
			if i > 0 {
				if err := out.WriteByte(text.LineFeedByte); err != nil {
					return err
				}
			}

			if err := out.WriteByte(text.LineFeedByte); err != nil {
				return err
			}

			if multiArgs {
				if err := argument.Render(arg, options, 1, out, i+1); err != nil {
					return err
				}
			} else {
				if err := argument.Render(arg, options, 1, out, 0); err != nil {
					return err
				}
			}
		}
	}

	return out.WriteByte(text.LineFeedByte)
}

func RenderCommandUsageBackHalf(com argo.Command, out *bufio.Writer) error {
	// If the command has flag groups
	if com.HasFlagGroups() {
		hasOptionalFlags := false

		// For all the required flags, append their name (and argument name if
		// required) to the cli example text.
		for _, group := range com.FlagGroups() {
			for _, f := range group.Flags() {
				if f.IsRequired() {
					if err := out.WriteByte(text.SpaceByte); err != nil {
						return err
					}
					if err := flag.RenderShortestLine(f, out); err != nil {
						return err
					}
				} else {
					hasOptionalFlags = true
				}
			}
		}

		// If there are any optional flags append the general "[OPTIONS]" text.
		if hasOptionalFlags {
			if _, err := out.WriteString(CommandRenderOptionalFlags); err != nil {
				return err
			}
		}
	}

	// After all the flag groups have been rendered, append the argument names.
	if com.HasArguments() {
		multiArgs := len(com.Arguments()) > 1

		for i, arg := range com.Arguments() {
			if err := out.WriteByte(text.SpaceByte); err != nil {
				return err
			}
			if multiArgs {
				if err := argument.RenderName(arg, out, i+1); err != nil {
					return err
				}
			} else {
				if err := argument.RenderName(arg, out, 0); err != nil {
					return err
				}
			}
		}
	}

	if com.HasUnmappedInputLabel() {
		if _, err := fmt.Fprintf(out, " %c%s%c", argument.OptPrefix, com.UnmappedInputLabel(), argument.OptSuffix); err != nil {
			return err
		}
	}

	return nil
}
