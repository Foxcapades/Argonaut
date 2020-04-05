package argo

type Flag interface {

	// Short returns this flag's short notation character.
	//
	// If this flag does not have a short notation set, this
	// method returns a null character (\0).  Test whether or
	// not this flag has a short notation by using HasShort().
	Short() byte

	// HasShort returns whether or not this flag has a short
	// notation set.
	HasShort() bool

	// Long returns this flag's long notation name.
	//
	// If this flag does not have a long notation set, this
	// method returns an empty string.  Test whether or not
	// this flag has a long notation by using HasLong().
	Long() string

	// HasLong returns whether or not this flag has a long
	// notation set.
	HasLong() bool

	// Required returns whether or not this flag has been
	// marked as required.
	Required() bool

	// Argument returns this flag's argument definition.
	//
	// If this flag does not have an argument set, this method
	// returns nil.  Test whether or not this flag has an
	// argument set using HasArgument()
	Argument() Argument

	// Description returns this flag's description text.
	//
	// If this flag does not have description text set, this
	// method returns an empty string.  Test whether or not
	// this flag has a description by using HasDescription().
	Description() string

	// HasDescription returns whether or not this flag has a
	// description set.
	HasDescription() bool

	Hits() int
}
