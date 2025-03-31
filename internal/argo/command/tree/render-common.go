package tree

import (
	"slices"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Described interface {
	HasDescription() bool
	Description() string
}

type Aliased interface {
	HasAliases() bool
	Aliases() []string
}

func TryRenderAliases(node Aliased, w *utils.BatchWriter) {
	if !node.HasAliases() {
		return
	}
	w.WriteString(render.SubLinePadding[0])
	w.WriteString("Aliases: ")

	aliases := node.Aliases()
	slices.Sort(aliases)

	w.WriteString(aliases[0])

	for i := 1; i < len(aliases); i++ {
		w.WriteString(", ")
		w.WriteString(aliases[i])
	}

	w.WriteByte(text.LineFeedByte)
}

func TryRenderDescription(node Described, opts Options, w *utils.BatchWriter, preNl bool) {
	if node.HasDescription() {
		if preNl {
			w.WriteByte(text.LineFeedByte)
		}

		render.NewDescriptionFormatter(render.DescriptionPadding[0], opts.HelpTextMaxWidth(), w).Format(node.Description())
		w.WriteByte(text.LineFeedByte)
	}
}

func TryRenderFlags(node flag.GroupContainer, opts Options, w *utils.BatchWriter) {
	if node.HasFlagGroups() {
		w.WriteByte(text.LineFeedByte)
		flag.RenderGroups(node.FlagGroups(), opts, 0, w)
	}
}

const (
	subcommandPlaceholder = " " + string(argument.ReqPrefix) + "command" + string(argument.ReqSuffix)
)

type Named interface {
	Name() string
}

func RenderSubCommandPath(node Named, out *utils.BatchWriter) {
	path := make([]Named, 0, 4)

	current := node
	for {
		path = append(path, current)

		if t, ok := current.(argo.ChildNode); ok {
			current = t.Parent().(Named)
		} else {
			break
		}
	}

	slices.Reverse(path)

	out.WriteString(render.SubLinePadding[0])

	for i, segment := range path {
		if i > 0 {
			out.WriteByte(text.SpaceByte)
		}

		out.WriteString(segment.Name())

		for _, f := range flag.Stream(segment.(flag.GroupContainer)) {
			if f.IsRequired() {
				out.WriteByte(text.SpaceByte)
				flag.RenderShortestForUsage(f, out)
			}
		}
	}
}

type FlagInheritor interface {
	argo.ChildNode
	flag.GroupContainer
}

func TryRenderInheritedFlags(node FlagInheritor, options Options, out *utils.BatchWriter) {
	if options.InheritParentFlags() == argo.FlagInheritanceEnabled {
		renderInheritedFlagsFlat(node, options, out)
		return
	}

	if options.InheritParentFlags() == argo.FlagInheritanceGrouped {
		renderInheritedFlagsGrouped(node, options, out)
		return
	}
}

func renderInheritedFlagsFlat(node FlagInheritor, options Options, out *utils.BatchWriter) {
	inherited := flag.FlattenInheritance(node)

	if len(inherited) == 0 {
		return
	}

	out.WriteByte(text.LineFeedByte)

	out.WriteString("\nInherited Flags")

	for i := range inherited {
		if i > 0 && !inherited[i-1].Flag.HasDescription() {
			out.WriteByte(text.LineFeedByte)
		}
		out.WriteByte(text.LineFeedByte)
		flag.RenderInheritedForms(&inherited[i], options, 1, out)
	}
}

func renderInheritedFlagsGrouped(node FlagInheritor, options Options, sb *utils.BatchWriter) {
	grouped := flag.GroupInheritance(node)

	for i := range grouped {
		flag.RenderGroupedInheritance(&grouped[i], options, 1, sb)
	}
}
