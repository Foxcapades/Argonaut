package tree

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func (c *CommandTreeInterpreter) handlePlainText(
	element parse.Element,
	unmapped []string,
	errs argo.MultiError, // Only use for execution errors!  Invalid command structure errors should be passed up!
) ([]string, error) {
	// If we've hit the leaf node, then the plain text becomes an argument on
	// that node.  If we haven't yet hit the leaf node, then we must treat the
	// plaintext value as the name of the next node in the tree.  If no such
	// node exists, that is an error.
	if _, ok := c.current.(argo.LeafCommand); ok {
		if ok, err := c.appender.Append(element.String()); err != nil {
			return unmapped, err
		} else if !ok {
			return append(unmapped, element.String()), nil
		}

		// argument value was accepted
		return unmapped, nil
	}

	if node, ok := c.current.(argo.ParentNode); ok {
		// Lookup a child with the given input string
		if child := node.FindChild(element.String()); child != nil {
			// If flag inheritance is disabled, process the already given flags and
			// clear the queue.
			if c.tree.Options().InheritParentFlags == argo.FlagInheritanceDisabled {
				c.shiftFlagQueue(c.current.(flag.GroupContainer), errs)
				c.flagIndex = flag.BuildFlagIndex(child.(flag.GroupContainer))
			} else {
				c.flagIndex.Overlay(child.(flag.GroupContainer))
			}

			c.current = child
			node.SelectChild(element.String())

			if branch, ok := child.(argo.BranchCommand); ok {
				c.branches = append(c.branches, branch)
			} else if leaf, ok := child.(argo.LeafCommand); ok {
				c.leaf = leaf
				c.appender = argument.NewValueAppender(c.leaf.Arguments())
			} else {
				panic(fmt.Sprintf("unrecognized child node type %T", child))
			}

			return unmapped, nil
		} else {
			// If node child could be found matching the input string, then print
			// out a help message about the invalid subcommand.
			return unmapped, c.invalidSubCommand(element.String())
		}
	}

	panic("illegal state: command node was neither a leaf or a parent")
}
