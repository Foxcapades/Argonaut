package common

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const (
	CommandRenderPrefix        = "Usage:\n"
	CommandRenderOptionalFlags = " " + string(argument.OptPrefix) + "options" + string(argument.OptSuffix)
	CommandRenderArgs          = "Arguments"
)

type HelpRenderer[T any] = func(T, *utils.BatchWriter) error

type CallEndNode interface {
	flag.GroupContainer
	HasDescription() bool
	Description() string
	HasArguments() bool
	Arguments() []argo.Argument
}

func RenderCommandBackHalf(com CallEndNode, options opts.Options, out *utils.BatchWriter) {
	// If the command has a description, append it.
	if com.HasDescription() {
		out.WriteByte(text.LineFeedByte)

		render.NewDescriptionFormatter(render.DescriptionPadding[0], options.HelpTextMaxWidth(), out).
			Format(com.Description())
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
			out.WriteByte(text.LineFeedByte)
		}

		out.WriteByte(text.LineFeedByte)
		flag.RenderGroups(com.FlagGroups(), options, 0, out)
	}

	if writeArgs {
		out.WriteString(render.ParagraphBreak)
		out.WriteString(render.HeaderPadding[0])
		out.WriteString(CommandRenderArgs)

		multiArgs := len(com.Arguments()) > 1

		for i, arg := range com.Arguments() {
			if i > 0 {
				out.WriteByte(text.LineFeedByte)
			}

			out.WriteByte(text.LineFeedByte)

			if multiArgs {
				argument.Render(arg, options, 1, out, i+1)
			} else {
				argument.Render(arg, options, 1, out, 0)
			}
		}
	}

	out.WriteByte(text.LineFeedByte)
}

type CommandBackHalf interface {
	flag.GroupContainer
	argument.Container
}
