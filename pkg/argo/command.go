package argo

import (
	"os"
	"path/filepath"
)

// Command represents a singular, non-nested command which accepts flags and
// arguments.
type Command interface {

	// Name returns the name of the command.
	Name() string

	// Description returns the custom description for the command.
	//
	// Description values are used internally for rendering held text.
	Description() string

	// HasDescription indicates whether this command has a description value set.
	HasDescription() bool

	// FlagGroups returns the flag groups attached to this command.
	//
	// Flag groups are named categories of flags defined when building the
	// command.
	FlagGroups() []FlagGroup

	// HasFlagGroups indicates whether this command has any flag groups attached
	// to it.
	HasFlagGroups() bool

	// FindShortFlag looks up a Flag instance by its short form.
	//
	// If no such flag could be found on this command, this method will return
	// nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a Flag instance by its long form.
	//
	// If no such flag could be found on this command, this method will return
	// nil.
	FindLongFlag(name string) Flag

	// Arguments returns the positional Argument instances attached to this
	// Command.
	Arguments() []Argument

	// HasArguments indicates whether this Command has any positional arguments
	// attached.
	//
	// This method does not indicate whether those arguments were present on the
	// command line, it simply indicates whether Argument instances were attached
	// to the Command by the CommandBuilder.
	//
	// To determine whether an argument was present on the command line, test the
	// argument itself by using the Argument.WasHit method.
	HasArguments() bool

	appendArgument(rawArgument string) error

	// UnmappedInputs returns a collection of inputs that were passed to this
	// command that do not match any registered flag or argument.
	//
	// Unmapped inputs may be used to collect slices of positional arguments when
	// singular arguments can't be used.  For these situations, consider using
	// CommandBuilder.WithUnmappedLabel to set a help-text label indicating that
	// the command expects an arbitrary number of positional arguments.
	//
	// Defined positional arguments will always be hit before a value is added to
	// a command's unmapped inputs.
	UnmappedInputs() []string

	// HasUnmappedInputs indicates whether the command has collected any inputs
	// that were not mapped to any registered flag or argument.
	HasUnmappedInputs() bool

	appendUnmapped(val string)

	// PassthroughInputs are command line values that were passed after an
	// end-of-arguments boundary, "--".
	PassthroughInputs() []string

	// HasPassthroughInputs indicates whether this command has collected any
	// passthrough input values.
	HasPassthroughInputs() bool

	appendPassthrough(val string)

	// GetUnmappedLabel returns the label used when generating help text to
	// indicate the shape or purpose of unmapped inputs.
	GetUnmappedLabel() string

	// HasUnmappedLabel indicates whether an unmapped label has been set on this
	// command instance.
	HasUnmappedLabel() bool
}

type command struct {
	description   string
	unmappedLabel string
	flagGroups    []FlagGroup
	arguments     []Argument
	unmapped      []string
	passthrough   []string
}

func (c command) Name() string {
	return filepath.Base(os.Args[0])
}

func (c command) Description() string {
	return c.description
}

func (c command) HasDescription() bool {
	return len(c.description) > 0
}

func (c command) FlagGroups() []FlagGroup {
	return c.flagGroups
}

func (c command) HasFlagGroups() bool {
	return len(c.flagGroups) > 0
}

func (c command) HasUnmappedLabel() bool {
	return len(c.unmappedLabel) > 0
}

func (c command) GetUnmappedLabel() string {
	return c.unmappedLabel
}

func (c command) FindShortFlag(b byte) Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (c command) FindLongFlag(name string) Flag {
	for _, group := range c.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (c command) Arguments() []Argument {
	return c.arguments
}

func (c command) HasArguments() bool {
	return len(c.arguments) > 0
}

func (c *command) appendArgument(rawArgument string) error {
	for _, arg := range c.arguments {
		if !arg.WasHit() {
			if err := arg.setValue(rawArgument); err != nil {
				return err
			}
		}
	}

	c.unmapped = append(c.unmapped, rawArgument)
	return nil
}

func (c command) UnmappedInputs() []string {
	return c.unmapped
}

func (c command) HasUnmappedInputs() bool {
	return len(c.unmapped) > 0
}

func (c *command) appendUnmapped(val string) {
	c.unmapped = append(c.unmapped, val)
}

func (c command) PassthroughInputs() []string {
	return c.passthrough
}

func (c command) HasPassthroughInputs() bool {
	return len(c.passthrough) > 0
}

func (c *command) appendPassthrough(val string) {
	c.passthrough = append(c.passthrough, val)
}
