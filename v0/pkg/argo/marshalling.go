package argo

type Unmarshaler interface {
	Unmarshal(value string) (err error)
}

type ValueUnmarshaler interface {
	Unmarshal(string, interface{}) error
}
