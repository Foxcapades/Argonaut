package render

import (
	"bufio"
	"reflect"
	"strconv"

	"github.com/Foxcapades/Argonaut/v3/internal/chars"
	"github.com/Foxcapades/Argonaut/v3/pkg/argo/argument"
)

func Argument(arg cli_arg.Argument, padding uint8, out *bufio.Writer, argIndex int) error {
	if _, err := out.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}
	if err := ArgumentName(arg, out, argIndex); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, out)
		if err := formatter.Format(arg.Description()); err != nil {
			return err
		}
	}

	return nil
}

func ArgumentName(a cli_arg.Argument, out *bufio.Writer, argIndex int) error {
	if a.HasBinding() && a.BindingType().Kind() == reflect.Bool {
		return nil
	}

	if a.IsRequired() {
		if err := out.WriteByte(argReqPrefix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(argOptPrefix); err != nil {
			return err
		}
	}

	if _, err := out.WriteString(renderArgName(a, argIndex)); err != nil {
		return err
	}

	if a.IsRequired() {
		if err := out.WriteByte(argReqSuffix); err != nil {
			return err
		}
	} else {
		if err := out.WriteByte(argOptSuffix); err != nil {
			return err
		}
	}

	return nil
}

func FlagArgument(arg cli_arg.Argument, padding uint8, out *bufio.Writer) error {
	if _, err := out.WriteString(chars.SubLinePadding[padding]); err != nil {
		return err
	}

	if err := ArgumentName(arg, out, 0); err != nil {
		return err
	}

	if arg.HasDescription() {
		if err := out.WriteByte(chars.CharLF); err != nil {
			return err
		}

		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, out)
		if err := formatter.Format(arg.Description()); err != nil {
			return err
		}
	}

	return nil
}

func renderArgName(a cli_arg.Argument, argIndex int) string {
	if a.HasName() {
		return a.Name()
	}

	if argIndex > 0 {
		return "arg" + strconv.Itoa(argIndex)
	}

	return "arg"
}
