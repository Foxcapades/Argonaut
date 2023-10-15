package argo

type FlagRef interface {
	LongForm() string
	HasLongForm() bool

	ShortForm() byte
	HasShortForm() bool

	Value() string
	HasValue() bool
}
