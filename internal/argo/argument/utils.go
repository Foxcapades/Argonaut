package argument

import (
	"reflect"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func IsBoolean(arg argo.Argument) bool {
	return arg.HasBinding() && arg.BindingType().Kind() == reflect.Bool
}
