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
	
  // Name returns the name of the command or subcommand.
  Name() string

  // Description returns the description value assigned to this node.
  //
  // Description values are used when rendering help text.
  Description() string

  // HasDescription indicates whether this Node has a description value
  // set.
  HasDescription() bool

  // FlagGroups returns the flag groups assigned to this Node.
  //
  // This method will only return flag groups that had flags assigned to them,
  // the rest of the flag groups will have been filtered out when the node was
  // built.
  FlagGroups(includeDefault bool) []FlagGroup

  // HasFlagGroups indicates whether this Node has at least one populated
  // flag group.
  HasFlagGroups(includeDefault bool) bool

  // FindShortFlag looks up a target Flag instance by its short-form character.
  //
  // If no such flag exists on this Node or any of its parents, this
  // method will return nil.
  FindShortFlag(c byte) Flag

  // FindLongFlag looks up a target Flag instance by its long-form name.
  //
  // If no such flag exists on this Node or any of its parents, this
  // method will return nil.
  FindLongFlag(name string) Flag

  HasCallback() bool

  Callback() CommandCallback[CommandTree]

  IsHelpDisabled() bool

  IncompleteHandler() IncompleteCommandHandler[CommandTree]
  
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
	
  // WithDescription sets the description value that will be used for the built
  // cli command or subcommand.
  //
  // Descriptions are used when rendering help text.
  WithDescription(desc string) CommandTreeBuilder

  HasDescription() bool

  Description() string

  // WithHelpDisabled disables the automatic `-h` and `--help` flags for
  // rendering help text.
  WithHelpDisabled() CommandTreeBuilder

  IsHelpDisabled() bool

  // WithFlagGroup appends the given FlagGroupBuilder to this CLI component
  // builder.
  WithFlagGroup(group FlagGroupBuilder) CommandTreeBuilder

  WithFlagGroups(groups ...FlagGroupBuilder) CommandTreeBuilder

  HasFlagGroups(includeDefault bool) bool

  FlagGroups(includeDefault bool) []FlagGroupBuilder

  // WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
  // instance attached to this CLI component builder.
  WithFlag(flag FlagBuilder) CommandTreeBuilder

  WithFlags(flags ...FlagBuilder) CommandTreeBuilder

  HasFlags() bool

  WithCallback(callback CommandCallback[CommandTree]) CommandTreeBuilder

  HasCallback() bool

  Callback() CommandCallback[CommandTree]

  WithBranch(branch BranchCommandBuilder) CommandTreeBuilder
  WithBranches(branches ...BranchCommandBuilder) CommandTreeBuilder

  WithLeaf(leaf LeafCommandBuilder) CommandTreeBuilder
  WithLeaves(leaves ...LeafCommandBuilder) CommandTreeBuilder

  // WithCommandGroup appends the given command group builder to be built with
  // this command tree.
  //
  // Command groups are used for organizing subcommands into named groups that
  // are primarily used for rendering help text.
  WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder

  WithCommandGroups(groups ...CommandGroupBuilder) CommandTreeBuilder

  WithIncompleteHandler(handler IncompleteCommandHandler[CommandTree]) CommandTreeBuilder
  IncompleteHandler() IncompleteCommandHandler[CommandTree]
  
}
