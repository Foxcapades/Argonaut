package argo

import (
	"reflect"
)

type Unmarshaler interface {
	Unmarshal(value string) (err error)
}

type InternalUnmarshaler func(string, interface{}, UnmarshalProps) error

var (
	defaultUnmarshaler = InternalUnmarshaler(Unmarshal)
	unmarshalerType    = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	unmarshaller       = defaultUnmarshaler
)
