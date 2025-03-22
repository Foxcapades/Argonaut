package argo

import (
  . "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ParentNodeCommon interface {
  {{ define "ParentNodeCommon" }}
  IncompleteHandler() IncompleteCommandHandler[{{ .OutputType }}]
  {{ end }}
}

type ParentNodeBuilderCommon interface {
  {{ define "ParentNodeBuilderCommon" }}
  WithBranch(branch BranchCommandBuilder) {{ .BuilderType }}
  WithBranches(branches ...BranchCommandBuilder) {{ .BuilderType }}

  WithLeaf(leaf LeafCommandBuilder) {{ .BuilderType }}
  WithLeaves(leaves ...LeafCommandBuilder) {{ .BuilderType }}

  // WithCommandGroup appends the given command group builder to be built with
  // this command tree.
  //
  // Command groups are used for organizing subcommands into named groups that
  // are primarily used for rendering help text.
  WithCommandGroup(group CommandGroupBuilder) {{ .BuilderType }}

  WithCommandGroups(groups ...CommandGroupBuilder) {{ .BuilderType }}

  WithIncompleteHandler(handler IncompleteCommandHandler[{{ .OutputType }}]) {{ .BuilderType }}
  IncompleteHandler() IncompleteCommandHandler[{{ .OutputType }}]
  {{ end }}
}
