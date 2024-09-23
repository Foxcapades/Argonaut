package argo

type MultiError interface {
	error

	Size() int

	Get(idx int) error

	IsEmpty() bool
}
