package interpreter

import (
	"github.com/foxcapades/argonaut/internal/component"
	"github.com/foxcapades/argonaut/internal/token"
)

func eatPassthroughs(facade component.Facade, queue tokenQueue) {
	tkn := queue.PopFront()

	for tkn.Type != token.TypeEndOfInput {
		facade.AppendPassthrough(tkn.Value)
		tkn = queue.PopFront()
	}
}
