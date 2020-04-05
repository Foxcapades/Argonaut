package argo

type Unmarshaler interface {
	Unmarshal(value string) (err error)
}

type ValueUnmarshaler func(string, interface{}, UnmarshalProps) error
