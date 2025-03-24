package argument

import (
	"bufio"
	"reflect"
	"strconv"

	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

const (
	ReqPrefix = '<'
	ReqSuffix = '>'
	OptPrefix = '['
	OptSuffix = ']'
)

func ShouldBeRendered(arg argo.Argument) bool {
	return arg.IsRequired() ||
		!arg.HasBinding() ||
		arg.BindingType().Kind() != reflect.Bool
}

func Render(arg argo.Argument, options argo.Options, padding uint8, out *bufio.Writer, argIndex int) error {
	if _, err := out.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}
	if err := RenderName(arg, out, argIndex); err != nil {
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

func RenderName(a argo.Argument, out *bufio.Writer, argIndex int) error {
	if a.HasBinding() && a.BindingType().Kind() == reflect.Bool {
		return nil
	}

	if a.IsRequired() {
		if err := out.WriteByte(ReqPrefix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(OptPrefix); err != nil {
			return err
		}
	}

	if _, err := out.WriteString(renderArgName(a, argIndex)); err != nil {
		return err
	}

	if a.IsRequired() {
		if err := out.WriteByte(ReqSuffix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(OptSuffix); err != nil {
			return err
		}
	}

	return nil
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
