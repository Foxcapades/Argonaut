package argo

type Flag interface {
	ShortForm() byte
	HasShortForm() bool

	LongForm() string
	HasLongForm() bool

	Description() string
	HasDescription() bool

	Argument() Argument
	HasArgument() bool

	IsRequired() bool
	RequiresArgument() bool

	Hit() error

	HitWithArg(rawArg string) error

	WasHit() bool

	HitCount() int
}
