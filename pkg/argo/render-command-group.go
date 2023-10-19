package argo

import (
	"bufio"
	"slices"
)

const (
	defaultComGroupName = "Commands"
)

func renderCommandGroups(groups []CommandGroup, padding uint8, sb *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := sb.WriteString(paragraphBreak); err != nil {
				return err
			}
		}

		if err := renderCommandGroup(group, padding, sb); err != nil {
			return err
		}
	}

	return nil
}

func renderCommandGroup(group CommandGroup, padding uint8, sb *bufio.Writer) error {
	if _, err := sb.WriteString(headerPadding[padding]); err != nil {
		return err
	}

	if group.Name() == defaultGroupName {
		if _, err := sb.WriteString(defaultComGroupName); err != nil {
			return err
		}
	} else {
		if _, err := sb.WriteString(group.Name()); err != nil {
			return err
		}
	}

	if err := sb.WriteByte(charLF); err != nil {
		return err
	}

	if group.HasDescription() {
		if err := breakFmt(group.Description(), descriptionPadding[padding], helpTextMaxWidth, sb); err != nil {
			return err
		}
		if _, err := sb.WriteString(paragraphBreak); err != nil {
			return err
		}
	}

	ordered := make([]string, 0, len(group.Leaves())+len(group.Branches()))
	lookup := make(map[string]CommandChild, len(group.Leaves())+len(group.Branches()))
	maxLen := 0

	for _, node := range group.Leaves() {
		ordered = append(ordered, node.Name())
		lookup[node.Name()] = node
		if len(node.Name()) > maxLen {
			maxLen = len(node.Name())
		}
	}
	for _, node := range group.Branches() {
		ordered = append(ordered, node.Name())
		lookup[node.Name()] = node
		if len(node.Name()) > maxLen {
			maxLen = len(node.Name())
		}
	}

	slices.Sort(ordered)
	maxLen += 4

	for i, name := range ordered {
		if i > 0 {
			if err := sb.WriteByte(charLF); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(headerPadding[padding+1]); err != nil {
			return err
		}
		if _, err := sb.WriteString(name); err != nil {
			return err
		}

		node := lookup[name]

		if node.HasAliases() {
			if err := pad(maxLen-len(name), sb); err != nil {
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

		if node.HasDescription() {
			if err := sb.WriteByte(charLF); err != nil {
				return err
			}
			if err := breakFmt(node.Description(), descriptionPadding[padding+1], helpTextMaxWidth, sb); err != nil {
				return err
			}
		}
	}

	return nil
}
