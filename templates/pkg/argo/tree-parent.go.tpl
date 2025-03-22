package argo

import (
  . "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ParentNode interface {
  // CommandGroups returns the CommandGroup instances attached to this
  // ParentNode node.
  CommandGroups(includeDefault bool) []CommandGroup

  HasCommandGroups(includeDefault bool) bool

  // FindChild searches this ParentNode's CommandGroup instances for a
  // subcommand that matches the given string.
  //
  // A subcommand may match on either its name or one of its aliases.
  FindChild(name string) ChildNode

  HasSelectedChild() bool

  SelectChild(name string) bool

  SelectedChild() ChildNode[any]

  HasIncompleteHandler() bool
}

type ParentNodeBuilder interface {
  HasCommandGroups() bool

  CommandGroups() []CommandGroupBuilder

  HasSubcommands() bool

  HasIncompleteHandler() bool
}
