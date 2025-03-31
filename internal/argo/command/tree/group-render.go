package tree

import (
	"slices"

	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderCommandGroups(groups []argo.CommandGroup, options Options, padding uint8, sb *utils.BatchWriter) {
	if len(groups) == 0 {
		return
	}

	sb.WriteByte(text.LineFeedByte)

	for i, group := range groups {
		if i > 0 {
			sb.WriteString(render.ParagraphBreak)
		}

		RenderCommandGroup(group, options, padding, sb)
	}

	sb.WriteByte(text.LineFeedByte)
}

func RenderCommandGroup(group argo.CommandGroup, options Options, padding uint8, sb *utils.BatchWriter) {
	sb.WriteString(render.HeaderPadding[padding])

	if group.Name() == DefaultCommandGroupName {
		sb.WriteString(options.DefaultCommandGroupName())
	} else {
		sb.WriteString(group.Name())
	}

	sb.WriteByte(text.LineFeedByte)

	if group.HasDescription() {
		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb)
		formatter.Format(group.Description())
		sb.WriteString(render.ParagraphBreak)
	}

	ordered, lookup, maxLen := indexSubcommands(group.Branches(), group.Leaves())

	// Add padding
	maxLen += 4

	for i, name := range ordered {
		if i > 0 {
			sb.WriteByte(text.LineFeedByte)
		}

		sb.WriteString(render.HeaderPadding[padding+1])

		sb.WriteString(name)

		node := lookup[name]

		if node.HasAliases() {
			render.Pad(maxLen-len(name), sb)

			sb.WriteString("Aliases: ")

			slices.Sort(node.Aliases())

			for i, alias := range node.Aliases() {
				if i > 0 {
					sb.WriteString(", ")
				}

				sb.WriteString(alias)
			}
		}

		described := node.(interface {
			HasDescription() bool
			Description() string
		})

		if described.HasDescription() {
			sb.WriteByte(text.LineFeedByte)

			formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding+1], options.HelpTextMaxWidth(), sb)
			formatter.Format(described.Description())
		}
	}
}

func indexSubcommands(
	branches []argo.BranchCommand,
	leaves []argo.LeafCommand,
) (sortedNames []string, index map[string]argo.ChildNode, maxNameLength int) {
	sortedNames = make([]string, 0, len(leaves)+len(branches))
	index = make(map[string]argo.ChildNode, len(leaves)+len(branches))

	for _, node := range leaves {
		sortedNames = append(sortedNames, node.Name())
		index[node.Name()] = node
		if len(node.Name()) > maxNameLength {
			maxNameLength = len(node.Name())
		}
	}

	for _, node := range branches {
		sortedNames = append(sortedNames, node.Name())
		index[node.Name()] = node
		if len(node.Name()) > maxNameLength {
			maxNameLength = len(node.Name())
		}
	}

	slices.Sort(sortedNames)

	return
}
