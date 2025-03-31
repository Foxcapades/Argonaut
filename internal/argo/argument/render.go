package argument

import (
	"strconv"

	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const (
	ReqPrefix = text.OpenAngleBracket
	ReqSuffix = text.CloseAngleBracket
	OptPrefix = text.OpenSquareBracket
	OptSuffix = text.CloseSquareBracket
)

func ShouldBeRendered(arg argo.Argument) bool {
	return arg.IsRequired() || !IsBoolean(arg)
}

func Render(arg argo.Argument, options opts.Options, padding uint8, out *utils.BatchWriter, argIndex int) {
	out.WriteString(render.HeaderPadding[padding])
	if !IsBoolean(arg) {
		RenderName(arg, out, argIndex)
	}

	if arg.HasDescription() {
		out.WriteByte(text.LineFeedByte)

		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), out)
		formatter.Format(arg.Description())
	}
}

func RenderName(a argo.Argument, out *utils.BatchWriter, argIndex int) {
	if IsBoolean(a) {
		return
	}

	if a.IsRequired() {
		out.WriteByte(ReqPrefix)
		out.WriteString(renderArgName(a, argIndex))
		out.WriteByte(ReqSuffix)
	} else {
		out.WriteByte(OptPrefix)
		out.WriteString(renderArgName(a, argIndex))
		out.WriteByte(OptSuffix)
	}
}

func RenderForUsageLine(arguments []argo.Argument, writer *utils.BatchWriter) {
	multiArgs := len(arguments) > 1

	for i, arg := range arguments {
		writer.WriteByte(text.SpaceByte)
		if multiArgs {
			RenderName(arg, writer, i+1)
		} else {
			RenderName(arg, writer, 0)
		}
	}
}

func renderArgName(a argo.Argument, argIndex int) string {
	if a.HasName() {
		return a.Name()
	}

	if argIndex > 0 {
		return "arg" + strconv.Itoa(argIndex)
	}

	return "arg"
}
