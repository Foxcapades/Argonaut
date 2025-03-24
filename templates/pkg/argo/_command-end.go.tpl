package argo

type CommandEnd interface {
{{ define "CommandEnd" }}
  // Arguments returns the positional Argument instances attached to this
  // {{ .OutputType }}.
  Arguments() []Argument

  // HasArguments indicates whether this {{ .OutputType }} has any positional
	// Arguments attached.
  //
  // This method does not indicate whether those Arguments were present on the
  // command line, it only indicates whether any Argument instances were
  // attached to the command with the {{ .BuilderType }}.
  //
  // To determine whether an Argument was present on the command line, test the
  // Argument itself by using the Argument.WasHit method.
  HasArguments() bool

  // UnmappedInputs returns a collection of inputs that were passed to this
  // {{ .OutputType }} that do not match any registered flag or argument.
  //
  // Unmapped inputs may be used to collect slices of positional arguments when
  // singular arguments can't be used.  For these situations, consider using
  // {{ .BuilderType }}.WithUnmappedLabel to set a help-text label indicating
	// that the command expects an arbitrary number of positional arguments.
  //
  // Defined positional arguments will always be hit before a value is added to
  // a command's unmapped inputs.
  UnmappedInputs() []string

  // HasUnmappedInputs indicates whether the command has collected any inputs
  // that were not mapped to any registered flag or argument.
  HasUnmappedInputs() bool

  AppendUnmappedInput(val string)

  // UnmappedInputLabel returns the label used when generating help text to
  // indicate the shape or purpose of unmapped inputs.
  UnmappedInputLabel() string

  // HasUnmappedInputLabel indicates whether an unmapped label has been set on this
  // command.
  HasUnmappedInputLabel() bool
{{ end }}
}

type CommandEndBuilder interface {
{{ define "CommandEndBuilder" }}
  // WithArgument appends the given ArgumentBuilder instance to this
  // {{ .BuilderType }}'s list of positional arguments.
  WithArgument(arg ArgumentBuilder) {{ .BuilderType }}

  // WithArguments appends the given ArgumentBuilder instances to this
  // {{ .BuilderType }}'s list of positional arguments.
  WithArguments(args ...ArgumentBuilder) {{ .BuilderType }}

  // HasArguments indicates whether any positional arguments have been set on
  // this {{ .BuilderType }} instance.
  HasArguments() bool

  // Arguments returns a slice containing the ArgumentBuilders that have been
  // set on this {{ .BuilderType }} instance.
  Arguments() []ArgumentBuilder

  // WithUnmappedInputLabel sets the help-text label for unmapped arguments.
  //
  // This is useful when your command takes an arbitrary number of argument
  // inputs, and you would like the help text to indicate as such.
  //
  // Example Config:
  //     cli.{{ .CLIFunc }}().
  //         WithUnmappedInputLabel("[FILE...]")
  //
  // Example Result:
  //     Usage:
  //       my-command [FILE...]
  WithUnmappedInputLabel(label string) {{ .BuilderType }}

  // HasUnmappedInputLabel indicates whether an unmapped input label has been
  // set on this {{ .BuilderType }} instance.
  HasUnmappedInputLabel() bool

  // UnmappedInputLabel returns the unmapped input label set on this
  // {{ .BuilderType }} instance.
  //
  // If no unmapped input label has been set, this method returns an empty
  // string.
  UnmappedInputLabel() string
{{ end }}
}