package argo

type Unmarshaler interface {
	Unmarshal(raw string) error
}

type ValueUnmarshaler interface {
	Unmarshal(raw string, val any) error
}
