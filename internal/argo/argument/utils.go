package argument

import (
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/xreflect"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func IsBoolean(arg argo.Argument) bool {
	return arg.HasBinding() && xreflect.RootType(reflect.TypeOf(arg.Binding().BoundTo())).Kind() == reflect.Bool
}
