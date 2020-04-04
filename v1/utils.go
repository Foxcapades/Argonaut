package argo

import "reflect"

type props = uint8

const (
	hasArg  props = 1 << iota
	hasDef
	hasDesc
	hasHint
	isReq
)

func isUnmarshalable(val interface{}) bool {
	var tmp Unmarshaler
	tp := reflect.TypeOf(val)

	// Implements unmarshaler
	if un := reflect.TypeOf(tmp); tp.Implements(un) {
		return true
	}

	return false
}