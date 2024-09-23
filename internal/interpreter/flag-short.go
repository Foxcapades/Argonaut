package interpreter

import (
	"github.com/foxcapades/argonaut/internal/component"
	"github.com/foxcapades/argonaut/internal/token"
)

func handleShortFlag(tkn token.Token, facade component.Facade, queue tokenQueue) error {
	f := tkn.Value[1]

	if flag, ok := facade.ShortFlag(f); ok {

	}
}
