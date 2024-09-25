package flag

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type Flag struct {
	isRequired bool
	argument   argo.Argument
}
