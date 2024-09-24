package argo

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
