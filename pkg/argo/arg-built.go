package argo

type Argument interface {
	IsRequired() bool

	HasValue() bool

	Value() any
}
