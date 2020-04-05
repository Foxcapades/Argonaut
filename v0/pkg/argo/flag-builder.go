package argo

type FlagBuilder interface {
	// Short sets this flag's short notation to the given
	// character.
	//
	// Short flags will be preceded by a single `-`(dash) on
	// the command line, and may be chained (combined with a
	// single leading dash).
	//
	// A short flag character must be alphanumeric.
	//
	// Example Flags
	//
	//     -r -F -0
	//     -rF0
	//
	// Short flags can be passed an argument using one of the
	// following formats:
	//
	//     -lValue
	//     -l Value
	//     -l=Value
	//
	// Prioritization
	//
	// Due to the allowed syntax of short flag notation, we
	// run into some decision points when deciding whether a
	// character immediately following a short flag should be
	// treated as an additional flag or a value for the
	// preceding flag.  The decision making process will be as
	// follows:
	//
	// Given the input -abcd
	//
	// If `-a` requires an argument, then:
	//   Use `bcd` as the argument value for the `-a` flag.
	//
	// If `-a` has an optional argument:
	//   And there exists a flag `-b`, then:
	//     Mark `-a` as hit and move on to `-
	//   And there is no known flag `-b`, then
	//     Use `bcd` as the argument value for the `-a` flag.
	//
	// If `-a` expects no argument:
	//   Mark `-a` as hit and move on to `-b`
	Short(byte) FlagBuilder

	GetShort() byte

	HasShort() bool

	// Long sets this flag's long notation to the given
	// string.
	//
	// Long flags will be preceded by two `-`(dash) characters
	// on the command line.
	//
	// A long flag must start with an alphanumeric character
	// followed by zero or more alphanumeric characters,
	// dashes, and/or underscores.
	//
	// Example Usage
	//
	//     --verbose
	//     --some-flag
	//     --using_underscores
	//
	// Long flags can be passed an argument using the `=`
	// (equals) character:
	//
	//     --log-level=WARN
	Long(string) FlagBuilder

	GetLong() string

	HasLong() bool

	Description(string) FlagBuilder

	GetDescription() string

	HasDescription() bool

	Arg(ArgumentBuilder) FlagBuilder

	GetArg() ArgumentBuilder

	HasArg() bool

	// Bind is a convenience shorthand for passing a nameless
	// argument to `Arg()`
	//
	// If this flag already has an argument present, it will
	// be updated with these values.
	//
	// See ArgumentBuilder.Bind() and
	// ArgumentBuilder.Required() for more information.
	Bind(ptr interface{}, required bool) FlagBuilder

	// Default is a convenience shorthand for passing a
	// nameless argument to `Arg()`
	//
	// If this flag already has an argument present, it will
	// be updated with this value.
	//
	// See ArgumentBuilder.Default() for more information.
	Default(val interface{}) FlagBuilder

	Build() (Flag, error)

	MustBuild() Flag
}
