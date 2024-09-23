package util

type Property[T any] struct {
	Value T
	IsSet bool
}
