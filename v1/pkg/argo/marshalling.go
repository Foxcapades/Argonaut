package argo

type Unmarshaler interface {
	Unmarshal(value string) (err error)
}

type InternalUnmarshaler func(string, interface{}, UnmarshalProps) error
