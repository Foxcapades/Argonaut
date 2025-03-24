package tree

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/argument"
	"github.com/foxcapades/argonaut/v3/internal/parse"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func (c *CommandTreeInterpreter) handlePlainText(
	element parse.Element,
	arguments argument.Appender,
	unmapped *[]string,
	errs argo.MultiError, // Only use for execution errors!  Invalid command structure errors should be passed up!
) ([]string, error) {
	// If we've hit the leaf node, then the plain text becomes an argument on
	// that node.  If we haven't yet hit the leaf node, then we must treat the
	// plaintext value as the name of the next node in the tree.  If no such
	// node exists, that is an error.
	if _, ok := c.current.(argo.LeafCommand); ok {
		if ok, err := arguments.Append(element.String()); err != nil {
			return *unmapped, err
		} else if !ok {
			return append(*unmapped, element.String()), nil
		}

		// argument value was accepted
		return *unmapped, nil
	}

	if node, ok := c.current.(argo.ParentNode); ok {
		// Lookup a child with the given input string
		if child := node.FindChild(element.String()); child != nil {
			if !c.options.InheritParentFlags {
				c.shiftFlagQueue(errs)
			}

			c.current = child
			node.SelectChild(element.String())

			if branch, ok := child.(argo.BranchCommand); ok {
				c.branches = append(c.branches, branch)
			} else if leaf, ok := child.(argo.LeafCommand); ok {
				c.leaf = leaf
			}
		}

		// If node child could be found matching the input string, then print
		// out a help message about the invalid subcommand.
		return *unmapped, c.invalidSubCommand(element.String())
	}

	panic("illegal state: command node was neither a leaf or a parent")
}
