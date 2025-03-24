{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.InterfaceFields */ -}}
package argo

// File Documentation:
//
// This file defines templates that are used for methods that are type dependent
// and cannot be moved to the ParentNode interface.
//
// Methods that are not type dependent (do not use the input fields) should be
// added to the ParentNode interface instead of here.

import (
  . "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ParentNodeCommon interface {
  // NOTE: HasIncompleteHandler is set on the shared ParentNode interface.
  {{ define "ParentNodeCommon" }}
  // IncompleteHandler returns the incomplete command handler function that was
  // attached to this {{ .OutputType }} through the {{ .BuilderType }}.
  //
  // If no incomplete command handler function was attached, this method returns
  // nil.
  IncompleteHandler() IncompleteCommandHandler[{{ .OutputType }}]
  {{ end }}
}

type ParentNodeBuilderCommon interface {
  {{ define "ParentNodeBuilderCommon" }}
  // WithBranch appends the given BranchCommandBuilder instance to this
  // {{ .BuilderType }} instance.
  //
  // Branch commands added to this {{ .BuilderType }} are placed in the default
  // CommandGroupBuilder.
  WithBranch(branch BranchCommandBuilder) {{ .BuilderType }}

  // WithBranches appends the given BranchCommandBuilder instances to this
  // {{ .BuilderType }} instance.
  //
  // Branch commands added to this {{ .BuilderType }} are placed in the default
  // CommandGroupBuilder.
  WithBranches(branches ...BranchCommandBuilder) {{ .BuilderType }}

  // WithLeaf appends the given LeafCommandBuilder instance to this
  // {{ .BuilderType }} instance.
  //
  // Leaf commands added to this {{ .BuilderType }} are placed in the default
  // CommandGroupBuilder.
  WithLeaf(leaf LeafCommandBuilder) {{ .BuilderType }}

  // WithLeaves appends the given LeafCommandBuilder instances to this
  // {{ .BuilderType }} instance.
  //
  // Leaf commands added to this {{ .BuilderType }} are placed in the default
  // CommandGroupBuilder.
  WithLeaves(leaves ...LeafCommandBuilder) {{ .BuilderType }}

  // WithCommandGroup appends the given CommandGroupBuilder instance to this
  // {{ .BuilderType }} instance.
  //
  // Command groups are used for organizing subcommands into named groups that
  // are primarily used for rendering help text.
  //
  // Example usage:
  //   cli.Tree().
  //       WithCommandGroup(cli.CommandGroup("My Command Group").
  //       WithBranch(cli.Branch("foo").
  //           WithLeaf(cli.Leaf("bar"))))
  //
  // Resulting help text:
  //
  WithCommandGroup(group CommandGroupBuilder) {{ .BuilderType }}

  WithCommandGroups(groups ...CommandGroupBuilder) {{ .BuilderType }}

  WithIncompleteHandler(handler IncompleteCommandHandler[{{ .OutputType }}]) {{ .BuilderType }}
  IncompleteHandler() IncompleteCommandHandler[{{ .OutputType }}]
  {{ end }}
}
