package argo

import (
	. "github.com/foxcapades/argonaut/v3/pkg/argo"
)

type CommandCommon interface {
{{ define "CommandCommon" }}
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
  FlagGroups() []FlagGroup

  // HasFlagGroups indicates whether this Node has at least one populated
  // flag group.
  HasFlagGroups() bool

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

  Callback() CommandCallback[{{ .OutputType }}]

  IsHelpDisabled() bool
{{ end }}
}

type CommandBuilderCommon interface {
{{ define "CommandBuilderCommon" }}
  // WithDescription sets the description value that will be used for the built
  // cli command or subcommand.
  //
  // Descriptions are used when rendering help text.
  WithDescription(desc string) {{ .BuilderType }}

  HasDescription() bool

  Description() string

  // WithHelpDisabled disables the automatic `-h` and `--help` flags for
  // rendering help text.
  WithHelpDisabled() {{ .BuilderType }}

  IsHelpDisabled() bool

  // WithFlagGroup appends the given FlagGroupBuilder to this CLI component
  // builder.
  WithFlagGroup(group FlagGroupBuilder) {{ .BuilderType }}

  WithFlagGroups(groups ...FlagGroupBuilder) {{ .BuilderType }}

  HasFlagGroups(includeDefault bool) bool

  FlagGroups(includeDefault bool) []FlagGroupBuilder

  // WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
  // instance attached to this CLI component builder.
  WithFlag(flag FlagBuilder) {{ .BuilderType }}

  WithFlags(flags ...FlagBuilder) {{ .BuilderType }}

  HasFlags() bool

  WithCallback(callback CommandCallback[{{ .OutputType }}]) {{ .BuilderType }}

  Callback() CommandCallback[{{ .OutputType }}]

  HasCallback() bool
{{ end }}
}