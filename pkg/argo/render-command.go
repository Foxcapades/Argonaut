package argo

import (
	"bufio"
	"io"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

// CommandHelpRenderer returns a HelpRenderer instance that is suited to
// rendering help text for Command instances.
func CommandHelpRenderer() HelpRenderer[Command] {
	return comRenderer{}
}

type comRenderer struct{ renderCommandBase }

func (r comRenderer) RenderHelp(command Command, writer io.Writer) error {
	if buf, ok := writer.(*bufio.Writer); ok {
		return r.renderCommand(command, buf)
	} else {
		buf := bufio.NewWriter(writer)
		err := r.renderCommand(command, buf)
		_ = buf.Flush()
		return err
	}
}

func (r comRenderer) renderCommand(com Command, out *bufio.Writer) error {
	if err := r.renderCommandUsageBlock(com, out); err != nil {
		return err
	}
	if err := out.WriteByte(chars.CharLF); err != nil {
		return err
	}
	return r.renderCommandBackHalf(com, out)
}

func (r comRenderer) renderCommandUsageBlock(com Command, out *bufio.Writer) error {
	if _, err := out.WriteString(comPrefix); err != nil {
		return err
	}
	if _, err := out.WriteString(chars.SubLinePadding[0]); err != nil {
		return err
	}
	if _, err := out.WriteString(com.Name()); err != nil {
		return err
	}
	return r.renderCommandUsageBackHalf(com, out)
}
