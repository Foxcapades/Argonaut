{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.InterfaceFields */ -}}
package argo

import (
  . "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ChildNodeBuilderCommon interface {
{{ define "ChildNodeBuilderCommon" }}
  // WithAlias assigns the given alias to the target command node.
  //
  // Command aliases must be unique per level in a command tree.  This means
  // that for any given step in the tree, no alias may conflict with another
  // branch or leaf subcommand's name or aliases.
  //
  // This also applies if a subcommand node is reused at multiple levels of the
  // command tree.
  //
  // If a conflict is found between subcommand names and/or aliases, an error
  // will be returned when attempting to build the command tree.
  //
  // Example:
  //   cli.{{ .CLIFunc }}("list").WithAlias("ls")
  WithAlias(alias string) {{ .BuilderType }}

  // WithAliases assigns the given aliases to the target command node.
  //
  // Command aliases must be unique per level in a command tree.  This means
  // that for any given step in the tree, no alias may conflict with another
  // branch or leaf subcommand's name or aliases.
  //
  // This also applies if a subcommand node is reused at multiple levels of the
  // command tree.
  //
  // If a conflict is found between subcommand names and/or aliases, an error
  // will be returned when attempting to build the command tree.
  WithAliases(aliases ...string) {{ .BuilderType }}
{{ end }}
}
