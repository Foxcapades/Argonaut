package interpreter

import (
	"github.com/foxcapades/argonaut/internal/component"
	"github.com/foxcapades/argonaut/internal/errs"
	"github.com/foxcapades/argonaut/internal/token"
	"github.com/foxcapades/argonaut/internal/util/coll"
	"github.com/foxcapades/argonaut/pkg/argo"
)

type tokenQueue struct {
	queue  coll.Deque[token.Token]
	stream token.Stream
}

func (t tokenQueue) PopFront() token.Token {
	return t.queue.Poll()
}

func (t tokenQueue) PushFront(tkn token.Token) {
	t.queue.PushFront(tkn)
}

func (t tokenQueue) PushBack(tkn token.Token) {
	t.queue.PushBack(tkn)
}

func Run(input []string, facade component.Facade, options *argo.Config) error {
	queue := tokenQueue{
		queue:  coll.NewDeque[token.Token](2),
		stream: token.NewStream(input, options),
	}

	errors := errs.NewMultiError()

	for {
		tkn := queue.PopFront()

		switch tkn.Type {
		case token.TypeShort:
			errors.AppendIfNotNil(handleShortFlag(tkn, facade, queue))

		case token.TypeLong:

		case token.TypeArgument:

		case token.TypeEndOfOptions:
			eatPassthroughs(facade, &stream)
			break

		case token.TypeEndOfInput:
			break
		}
	}
}
