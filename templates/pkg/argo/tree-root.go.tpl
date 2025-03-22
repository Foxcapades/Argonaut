{{ $vars := (types "CommandTree" "CommandTreeBuilder" "Tree") -}}
package argo

// CommandTree represents the root of a tree of subcommands.
//
// The command tree consists of branch and leaf nodes.  The branch nodes can be
// thought of as categories for containing sub-branches and/or leaves.  Leaf
// nodes are the actual callable command implementations.
//
// All levels of the command tree accept flags, with sub-node flags taking
// priority over parent node flags on flag collision.  Leaf nodes, however, are
// the only nodes that accept positional arguments, or passthroughs.
//
// Example command tree:
//
//	docker
//	 |- compose
//	 |   |- build
//	 |   |- down
//	 |   |- ...
//	 |- container
//	 |   |- exec
//	 |   |- ls
//	 |   |- ...
//	 |- ...
type CommandTree interface {
	ParentNode
	{{ template "CommandCommon" $vars }}
  {{- template "ParentNodeCommon" $vars }}
  // SelectedCommand returns the leaf command that was selected in the CLI call.
  SelectedCommand() LeafCommand
}

// A CommandTreeBuilder is a builder type used to construct a CommandTree
// instance.
//
// A CommandTree is a command that consists of branching subcommands.  Examples
// of such commands include the `go` command, `docker`, or `kubectl`.
//
// To use the Docker command example we have a command tree that includes the
// following:
//
//	docker
//	 |- compose
//	 |   |- build
//	 |   |- down
//	 |   |- ...
//	 |- container
//	 |   |- exec
//	 |   |- ls
//	 |   |- ...
//	 |- ...
type CommandTreeBuilder interface {
	ParentNodeBuilder
	{{ template "CommandBuilderCommon" $vars }}
  {{- template "ParentNodeBuilderCommon" $vars }}
}
