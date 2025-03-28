package tree

import (
	"bufio"
	"io"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderLeafHelp(leaf argo.LeafCommand, options argo.Options, writer io.Writer) error {
	if leaf.IsHelpDisabled() {
		return nil
	}

	var buf *bufio.Writer

	if b, ok := writer.(*bufio.Writer); ok {
		buf = b
	} else {
		buf = bufio.NewWriter(writer)
	}

	defer utils.DisregardError(buf.Flush)
	return renderCommandLeaf(leaf, options, buf)
}

func MakeRenderLeafHelpCallback(leaf argo.LeafCommand, options argo.Options) argo.FlagCallback {
	return func(flag argo.Flag) {
		utils.Must(RenderLeafHelp(leaf, options, os.Stdout))
		os.Exit(0)
	}
}

func renderCommandLeaf(leaf argo.LeafCommand, options argo.Options, out *bufio.Writer) error {
	if err := renderCommandLeafUsage(leaf, out); err != nil {
		return err
	}

	if err := out.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	if err := TryRenderAliases(leaf, out); err != nil {
		return err
	}

	return renderCommandLeafBackHalf(leaf, options, out)
}

func renderCommandLeafUsage(leaf argo.LeafCommand, out *bufio.Writer) error {
	if _, err := out.WriteString(common.CommandRenderPrefix); err != nil {
		return err
	}
	if err := RenderSubCommandPath(leaf, out); err != nil {
		return err
	}
	return common.RenderCommandUsageLineBackHalf(leaf, out)
}

func renderCommandLeafBackHalf(com argo.LeafCommand, options argo.Options, out *bufio.Writer) error {

	// If the command has a description, append it.
	if err := TryRenderDescription(com, options, out, com.HasAliases()); err != nil {
		return err
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

	if err := TryRenderFlags(com, options, out); err != nil {
		return err
	}

	if err := TryRenderInheritedFlags(com, options, out); err != nil {
		return err
	}

	if writeArgs {
		if _, err := out.WriteString(render.ParagraphBreak); err != nil {
			return err
		}
		if _, err := out.WriteString(render.HeaderPadding[0]); err != nil {
			return err
		}
		if _, err := out.WriteString(common.CommandRenderArgs); err != nil {
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
