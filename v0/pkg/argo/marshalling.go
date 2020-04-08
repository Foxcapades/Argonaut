package argo

type Unmarshaler interface {
	Unmarshal(value string) (err error)
}

type SpecializedUnmarshaler interface {
	Unmarshaler

	ConsumesArguments() bool
}

type ValueUnmarshaler interface {
	Unmarshal(string, interface{}) error
}
