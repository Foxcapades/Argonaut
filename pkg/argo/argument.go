package argo

//

type Argument interface {
	// IsRequired returns an indicator for whether the Argument was marked as
	// required when configured.
	IsRequired() bool

	// UsedDefault returns an indicator for whether the Argument's default value
	// was used.
	//
	// If the returned value is `true`, then no value was provided for this
	// Argument in the CLI call.
	//
	// If the returned value is `false`, then the CLI call provided a value for
	// this Argument.
	UsedDefault() bool

	// HasValue indicates whether this Argument has a value set.
	//
	// If the Argument was configured with a default value, this method will
	// always return true.
	//
	// If this method returns `false`, calls to Value will panic.
	HasValue() bool

	// Value returns the value set for this Argument, if one was set.
	//
	// If no value is set, this method will panic.  Test whether a value is
	// available by calling HasValue before calling this method.
	//
	// Also see ValueOrNil.
	Value() any

	// ValueOrNil returns either the value set for this Argument, or `nil` if this
	// argument does not have a value.
	//
	// If HasValue returns `false`, this method will return `nil`.  If HasValue
	// returns `true`, this method will return the set value, which may also be
	// `nil`.
	//
	// Also see Value.
	ValueOrNil() any
}

//

type ArgumentSpec interface {
	// Name returns the name configured on this ArgumentSpec.
	Name() string

	// Summary returns a short summary text describing this ArgumentSpec.
	//
	// If this ArgumentSpec has no help text available, this method will return an
	// empty string.  This may be tested for using HasHelpText.
	Summary() string

	// Description returns the longform description of this ArgumentSpec.
	//
	// If this ArgumentSpec has no description available, this method will return
	// an empty string.  This may be tested for using HasHelpText.
	Description() string

	// HasHelpText returns an indicator as to whether this ArgumentSpec has help
	// text available.
	//
	// If the ArgumentSpec has any help text, both Summary and Description will
	// return a value, though the value may be the same for both.
	HasHelpText() bool

	// IsRequired returns an indicator as to whether this ArgumentSpec has been
	// marked as required.
	IsRequired() bool

	// HasValue returns whether this ArgumentSpec has had a value provided either
	// from the CLI call or from a default.
	HasValue() bool

	// Value returns the value set on this ArgumentSpec.
	//
	// If no value has been set on this ArgumentSpec, calls to this method will
	// panic.  Whether a value is available should be tested by using HasValue.
	Value() any

	// PreValidate performs any configured pre-parse validation on a raw CLI input
	// string before any attempt is made to parse the value.
	//
	// This method is separate from ProcessInput as it may be used to decide what
	// ArgumentSpec instance should be given a CLI input value when there is
	// ambiguity in the CLI definition.
	PreValidate(rawInput string) error

	// ProcessInput parses the raw CLI input string and performs any post-parse
	// validation on it before setting it as the value for this ArgumentSpec.
	//
	// If this method returns an error, a value parsed from the given input MUST
	// NOT be set on the ArgumentSpec instance.  For example, if HasValue returns
	// `false` before this method is called, and this method returns an error,
	// HasValue MUST still return `false` after this method returns.
	//
	// If this method does not return an error, the value parsed from the given
	// input MUST be used as the value for the ArgumentSpec instance.  Meaning,
	// after this method is called and returns `nil`, HasValue MUST return `true`.
	ProcessInput(rawInput string) error

	ToArgument() Argument
}

//

type ArgumentSpecBuilder interface {
	// Build validates the state of this ArgumentSpecBuilder instance and produces
	// a new ArgumentSpec based on the options configured on the builder.
	//
	// If the ArgumentSpecBuilder has been configured incorrectly, an error will
	// be returned.
	Build(config Config) (ArgumentSpec, error)
}

//

type TypedArgumentSpecBuilder[T any] interface {
	ArgumentSpecBuilder

	// WithName sets a display name for this ArgumentSpecBuilder which will be used when
	// rendering help text.
	//
	// If no name is provided for the ArgumentSpecBuilder, a default representation will
	// be used in help text output.
	WithName(name string) TypedArgumentSpecBuilder[T]

	// WithDescription sets a longform description for this ArgumentSpecBuilder which may
	// be used when rendering full help text.
	//
	// If no description is set, but a summary is provided, the summary will also
	// be used as the description.
	WithDescription(desc string) TypedArgumentSpecBuilder[T]

	// WithSummary sets a short description for this ArgumentSpecBuilder which may be
	// used when rendering short help text.
	//
	// If no summary is set, but a description is provided, the first line of that
	// description will be used as the summary.
	WithSummary(summary string) TypedArgumentSpecBuilder[T]

	// WithRawDefault sets a raw/unparsed default value for the ArgumentSpecBuilder.
	//
	// If the Argument is not provided a value at runtime, this raw default will
	// be parsed into a value of type T, then used as the Argument value.
	//
	// Default values skip both pre- and post-validation steps.
	//
	// If a custom unmarshaler is provided via WithUnmarshaler, that unmarshaler
	// will be used for this value.
	WithRawDefault(value string) TypedArgumentSpecBuilder[T]

	// WithDefault sets a default value for the ArgumentSpecBuilder.
	//
	// If the Argument is not provided a value at runtime, this default will be
	// used as the Argument value.
	//
	// Default values skip both pre- and post-validation steps.
	WithDefault(value T) TypedArgumentSpecBuilder[T]

	// WithDefaultProvider configures the ArgumentSpecBuilder to use the given
	// function to get the default value for the resulting Argument.
	//
	// If the Argument is not provided a value at runtime, this provider will be
	// called once to get a default value for the Argument.  If the Argument _is_
	// provided a value at runtime, this function will not be called.
	//
	// Default values skip both pre- and post-validation steps.
	WithDefaultProvider(provider func() T) TypedArgumentSpecBuilder[T]

	// WithPreValidator adds a PreParseArgumentValidator to this ArgumentSpecBuilder.
	//
	// Pre-validation is handled before Argonaut attempts to parse the raw input
	// value.
	//
	// Validators are called in the order they are added to the ArgumentSpecBuilder.
	WithPreValidator(validator PreParseArgumentValidator) TypedArgumentSpecBuilder[T]

	// WithPostValidator adds a PostParseArgumentValidator to this ArgumentSpecBuilder.
	//
	// Post-validation is handled after Argonaut successfully parses the raw input
	// value.
	//
	// If the raw input value could not be parsed, post-validation does not occur.
	//
	// Validators are called in the order they are added to the ArgumentSpecBuilder.
	WithPostValidator(validator PostParseArgumentValidator[T]) TypedArgumentSpecBuilder[T]

	// WithValueConsumer adds a consumer function that will be called with the
	// parsed value after it has been validated.
	WithValueConsumer(consumer ArgumentValueConsumer[T]) TypedArgumentSpecBuilder[T]

	// WithBinding adds a pointer binding to this ArgumentSpecBuilder which will
	// be filled with the parsed value after it has been validated.
	WithBinding(binding *T) TypedArgumentSpecBuilder[T]

	// WithDeepBinding adds a pointer of arbitrary depth to this
	// ArgumentSpecBuilder which will be filled with the parsed value after it has
	// been validated.
	//
	// As this is any typed, it falls to run-time checks to verify that the type
	// of the binding is valid.
	WithDeepBinding(binding any) TypedArgumentSpecBuilder[T]

	// WithUnmarshaler sets a custom ValueUnmarshaler to use when deserializing
	// the raw CLI input into a value of type T.
	WithUnmarshaler(unmarshaler ValueUnmarshaler[T]) TypedArgumentSpecBuilder[T]

	// Required marks this ArgumentSpecBuilder as being required.
	//
	// Required arguments MUST be provided a value at runtime for their parent CLI
	// component usage to be valid.
	//
	// Marking an ArgumentSpecBuilder as required in addition to providing a default
	// value is an error and will result in an error being return when the CLI
	// parser is built.
	Required() TypedArgumentSpecBuilder[T]
}

//

type ArgumentContainer interface {
	Arguments() []ArgumentSpec

	HasArguments() bool
}

//

type ArgumentSpecContainer interface {
	Arguments() []ArgumentSpec

	HasArguments() bool
}

//

type ArgumentValueConsumer[T any] interface {
	Accept(value T) error
}

type ArgumentValueConsumerFn[T any] func(value T) error

func (a ArgumentValueConsumerFn[T]) Accept(value T) error {
	return a(value)
}

func SimpleArgumentValueConsumerFn[T any](fn func(value T)) ArgumentValueConsumer[T] {
	return ArgumentValueConsumerFn[T](func(value T) error {
		fn(value)
		return nil
	})
}

//

type ValueUnmarshaler[T any] interface {
	Unmarshal(rawInput string) (T, error)
}

type ValueUnmarshalerFn[T any] func(rawInput string) (T, error)

func (v ValueUnmarshalerFn[T]) Unmarshal(rawInput string) (T, error) {
	return v(rawInput)
}

//

type PreParseArgumentValidator interface {
	Validate(rawInput string, prev error) error
}

type PreParseArgumentValidatorFn func(rawInput string, prev error) error

func (p PreParseArgumentValidatorFn) Validate(rawInput string, prev error) error {
	return p(rawInput, prev)
}

func SimplePreParseArgumentValidatorFn(fn func(rawInput string) error) PreParseArgumentValidator {
	return PreParseArgumentValidatorFn(func(rawInput string, _ error) error { return fn(rawInput) })
}

//

type PostParseArgumentValidator[T any] interface {
	Validate(parsedValue T, rawInput string, prev error) error
}

type PostParseArgumentValidatorFn[T any] func(parsedValue T, rawInput string, prev error) error

func (p PostParseArgumentValidatorFn[T]) Validate(parsedValue T, rawInput string, prev error) error {
	return p(parsedValue, rawInput, prev)
}

func SimplePostParseArgumentValidatorFn[T any](fn func(parsedValue T) error) PostParseArgumentValidator[T] {
	return PostParseArgumentValidatorFn[T](func(parsedValue T, _ string, _ error) error { return fn(parsedValue) })
}
