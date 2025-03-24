package command

import (
	"bufio"
	"io"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderHelp(command argo.Command, opts argo.Options, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return renderCommand(command, opts, buf)
	} else {
		buf := bufio.NewWriter(writer)
		err := renderCommand(command, opts, buf)
		_ = buf.Flush()
		return err
	}
}

func renderCommand(com argo.Command, opts argo.Options, out *bufio.Writer) error {
	if err := renderCommandUsageBlock(com, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}
	return common.RenderCommandBackHalf(com, opts, out)
}

func renderCommandUsageBlock(com argo.Command, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(chars.SubLinePadding[0]); err != nil {
		return err
	}
	if _, err := out.WriteString(com.Name()); err != nil {
		return err
	}
	return common.RenderCommandUsageBackHalf(com, out)
}
