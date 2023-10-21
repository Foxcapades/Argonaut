package argo

import (
	"bufio"
	"slices"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

const (
	defaultComGroupName = "Commands"
)

func renderCommandGroups(groups []CommandGroup, padding uint8, sb *bufio.Writer) error {
	for i, group := range groups {
		if i > 0 {
			if _, err := sb.WriteString(chars.ParagraphBreak); err != nil {
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
	if _, err := sb.WriteString(chars.HeaderPadding[padding]); err != nil {
		return err
	}

	if group.Name() == chars.DefaultGroupName {
		if _, err := sb.WriteString(defaultComGroupName); err != nil {
			return err
		}
	} else {
		if _, err := sb.WriteString(group.Name()); err != nil {
			return err
		}
	}

	if err := sb.WriteByte(chars.CharLF); err != nil {
		return err
	}

	if group.HasDescription() {
		formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding], chars.HelpTextMaxWidth, sb)
		if err := formatter.Format(group.Description()); err != nil {
			return err
		}
		if _, err := sb.WriteString(chars.ParagraphBreak); err != nil {
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
			if err := sb.WriteByte(chars.CharLF); err != nil {
				return err
			}
		}

		if _, err := sb.WriteString(chars.HeaderPadding[padding+1]); err != nil {
			return err
		}
		if _, err := sb.WriteString(name); err != nil {
			return err
		}

		node := lookup[name]

		if node.HasAliases() {
			if err := chars.Pad(maxLen-len(name), sb); err != nil {
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
			if err := sb.WriteByte(chars.CharLF); err != nil {
				return err
			}

			formatter := chars.NewDescriptionFormatter(chars.DescriptionPadding[padding+1], chars.HelpTextMaxWidth, sb)
			if err := formatter.Format(node.Description()); err != nil {
				return err
			}
		}
	}

	return nil
}
