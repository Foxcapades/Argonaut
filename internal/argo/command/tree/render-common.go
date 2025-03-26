package tree

import (
	"bufio"
	"slices"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Aliased interface {
	HasAliases() bool
	Aliases() []string
	HasDescription() bool
	Description() string
}

func TryRenderAliases(node Aliased, w *bufio.Writer) (err error) {
	if !node.HasAliases() {
		return
	}

	if err = utils.Write2Strings(w, render.SubLinePadding[0], "Aliases: "); err != nil {
		return
	}

	aliases := node.Aliases()
	slices.Sort(aliases)

	if _, err = w.WriteString(aliases[0]); err != nil {
		return
	}

	for i := 1; i < len(aliases); i++ {
		if err = utils.Write2Strings(w, ", ", aliases[i]); err != nil {
			return
		}
	}

	return w.WriteByte(text.LineFeedByte)
}

func TryRenderDescription(node Aliased, opts argo.Options, w *bufio.Writer) (err error) {
	if !node.HasDescription() {
		return
	}

	if node.HasAliases() {
		if err = w.WriteByte(text.LineFeedByte); err != nil {
			return
		}
	}

	if err = render.NewDescriptionFormatter(render.DescriptionPadding[0], opts.HelpTextMaxWidth, w).Format(node.Description()); err != nil {
		return
	}

	return w.WriteByte(text.LineFeedByte)
}

func TryRenderFlags(node flag.GroupContainer, opts argo.Options, w *bufio.Writer) (err error) {
	if !node.HasFlagGroups() {
		return
	}

	if err = w.WriteByte(text.LineFeedByte); err != nil {
		return
	}
	if err = flag.RenderGroups(node.FlagGroups(), opts, 0, w); err != nil {
		return
	}

	return w.WriteByte(text.LineFeedByte)
}

const (
	subcommandPlaceholder = " " + string(argument.ReqPrefix) + "command" + string(argument.ReqSuffix)
)

type Named interface {
	Name() string
}

func RenderSubCommandPath(node Named, out *bufio.Writer) error {
	path := make([]string, 0, 4)

	current := node
	for {
		path = append(path, current.Name())

		if t, ok := current.(argo.ChildNode); ok {
			current = t.Parent().(Named)
		} else {
			break
		}
	}

	slices.Reverse(path)

	if _, err := out.WriteString(render.SubLinePadding[0]); err != nil {
		return err
	}

	for i, segment := range path {
		if i > 0 {
			if err := out.WriteByte(text.SpaceByte); err != nil {
				return err
			}
		}
		if _, err := out.WriteString(segment); err != nil {
			return err
		}
	}

	return nil
}
