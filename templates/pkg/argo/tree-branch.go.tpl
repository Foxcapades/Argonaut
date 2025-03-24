{{ $vars := (types "BranchCommand" "BranchCommandBuilder" "Branch" true) -}}
package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

{{ if false -}}
import (
	. "github.com/foxcapades/argonaut/v3/pkg/argo"
)
{{- end -}}

// BranchCommand represents a subcommand under a CommandTree that is an
// intermediate node between the tree root and an executable LeafCommand.
//
// CommandBranches enable the organization of subcommands into categories.
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
type BranchCommand interface {
	ChildNode
	ParentNode
  {{ template "CommandCommon" $vars }}
  {{- template "ParentNodeCommon" $vars }}
}

// A BranchCommandBuilder instance may be used to configure a new BranchCommand
// instance to be built.
//
// CommandBranches are intermediate steps between the root of the CommandTree
// and the LeafCommand instances.
//
// For example, given the following command example, the tree is "foo", the
// branch is "bar", and the leaf is "fizz":
//
//	./foo bar fizz
type BranchCommandBuilder interface {
	ChildNodeBuilder
	ParentNodeBuilder
  {{ template "CommandBuilderCommon" $vars }}
  {{- template "ParentNodeBuilderCommon" $vars }}
  {{- template "ChildNodeBuilderCommon" $vars }}
}

