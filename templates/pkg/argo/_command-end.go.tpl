package argo

type CommandEnd interface {
{{ define "CommandEnd" }}
  // Arguments returns the positional Argument instances attached to this
  // {{ .OutputType }}.
  Arguments() []Argument

  // HasArguments indicates whether this {{ .OutputType }} has any positional
	// arguments attached.
  //
  // This method does not indicate whether those arguments were present on the
  // command line, it simply indicates whether Argument instances were attached
  // to the command by the builder.
  //
  // To determine whether an argument was present on the command line, test the
  // argument itself by using the Argument.WasHit method.
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

  // GetUnmappedLabel returns the label used when generating help text to
  // indicate the shape or purpose of unmapped inputs.
  GetUnmappedLabel() string

  // HasUnmappedLabel indicates whether an unmapped label has been set on this
  // command.
  HasUnmappedLabel() bool
{{ end }}
}

type CommandEndBuilder interface {
{{ define "CommandEndBuilder" }}
  // WithArgument appends the given ArgumentBuilder to this {{ .BuilderType }}'s
  // list of positional arguments.
  WithArgument(arg ArgumentBuilder) {{ .BuilderType }}

  WithArguments(args ...ArgumentBuilder) {{ .BuilderType }}

  HasArguments() bool

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

  HasUnmappedInputLabel() bool

  UnmappedInputLabel() string
{{ end }}
}