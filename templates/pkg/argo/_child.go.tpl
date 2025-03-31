{{- /* gotype: github.com/foxcapades/argonaut/v3/scripts/generate/data.InterfaceFields */ -}}
package argo

import (
  . "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type ChildNodeCommon interface {
{{ define "ChildNodeCommon" }}
  // FindShortFlag looks up a target Flag instance by its short-form character.
  //
  // If no such flag exists on this node, nil will be returned.
  //
  // To search this node and all parent nodes, use FindShortFlagRecursive
  FindShortFlag(c byte) Flag

  // FindShortFlagRecursive looks up a target Flag instance by its short-form
  // character recursively on this node and all ancestor nodes.
  //
  // If no such flag exists in the tree hierarchy, nil will be returned.
  //
  // The behavior of this method is not affected by the InheritParentFlags
  // TreeCommandOption setting.
  FindShortFlagRecursive(c byte) Flag

  // FindLongFlag looks up a target Flag instance by its long-form name.
  //
  // If no such flag exists on this node, the ancestor nodes will be checked
  // parent by parent.
  //
  // To search this node and all parent nodes, use FindLongFlagRecursive
  FindLongFlag(name string) Flag

  // FindLongFlagRecursive looks up a target Flag instance by its long-form
  // name recursively on this node and all ancestor nodes.
  //
  // If no such flag exists in the tree hierarchy, nil will be returned.
  //
  // The behavior of this method is not affected by the InheritParentFlags
  // TreeCommandOption setting.
  FindLongFlagRecursive(name string) Flag
{{ end }}
}

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
