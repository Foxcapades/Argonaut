package tree

import (
	"bufio"
	"slices"

	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func RenderCommandGroups(groups []argo.CommandGroup, options Options, padding uint8, sb *bufio.Writer) error {
	if len(groups) == 0 {
		return nil
	}

	if err := sb.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	for i, group := range groups {
		if i > 0 {
			if _, err := sb.WriteString(render.ParagraphBreak); err != nil {
				return err
			}
		}

		if err := RenderCommandGroup(group, options, padding, sb); err != nil {
			return err
		}
	}

	return sb.WriteByte(text.LineFeedByte)
}

func RenderCommandGroup(group argo.CommandGroup, options Options, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(render.HeaderPadding[padding]); err != nil {
		return err
	}

	if group.Name() == DefaultCommandGroupName {
		if _, err := sb.WriteString(options.DefaultCommandGroupName()); err != nil {
			return err
		}
	} else {
		if _, err := sb.WriteString(group.Name()); err != nil {
			return err
		}
	}

	if err := sb.WriteByte(text.LineFeedByte); err != nil {
		return err
	}

	if group.HasDescription() {
		formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding], options.HelpTextMaxWidth(), sb)
		if err := formatter.Format(group.Description()); err != nil {
			return err
		}
		if _, err := sb.WriteString(render.ParagraphBreak); err != nil {
			return err
		}
	}

	ordered, lookup, maxLen := indexSubcommands(group.Branches(), group.Leaves())

	// Add padding
	maxLen += 4

	for i, name := range ordered {
		if i > 0 {
			if err := sb.WriteByte(text.LineFeedByte); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(render.HeaderPadding[padding+1]); err != nil {
			return err
		}

		if _, err := sb.WriteString(name); err != nil {
			return err
		}

		node := lookup[name]

		if node.HasAliases() {
			if err := render.Pad(maxLen-len(name), sb); err != nil {
				return err
			}

			if _, err := sb.WriteString("Aliases: "); err != nil {
				return err
			}

			slices.Sort(node.Aliases())

			for i, alias := range node.Aliases() {
				if i > 0 {
					if _, err := sb.WriteString(", "); err != nil {
						return err
					}
				}

				if _, err := sb.WriteString(alias); err != nil {
					return err
				}
			}
		}

		described := node.(interface {
			HasDescription() bool
			Description() string
		})

		if described.HasDescription() {
			if err := sb.WriteByte(text.LineFeedByte); err != nil {
				return err
			}

			formatter := render.NewDescriptionFormatter(render.DescriptionPadding[padding+1], options.HelpTextMaxWidth(), sb)
			if err := formatter.Format(described.Description()); err != nil {
				return err
			}
		}
	}

	return nil
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
