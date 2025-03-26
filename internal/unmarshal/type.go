package unmarshal

import (
	"reflect"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

var UnmarshalerType = reflect.TypeOf((*argo.Unmarshaler)(nil)).Elem()
