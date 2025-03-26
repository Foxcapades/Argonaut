package argument

import (
	"iter"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewValueAppender(arguments []argo.Argument) ValueAppender {
	return ValueAppender{stream: func(yield func(argo.Argument) bool) {
		for _, arg := range arguments {
			yield(arg)
		}
	}}
}

type ValueAppender struct {
	stream iter.Seq[argo.Argument]
}

func (a ValueAppender) Append(rawValue string) (bool, error) {
	next, _ := iter.Pull(a.stream)

	for {
		arg, ok := next()

		if !ok {
			break
		}

		if arg.WasHit() {
			continue
		}

		return true, arg.SetValue(rawValue)
	}

	return false, nil
}

func (a ValueAppender) Remaining() iter.Seq[argo.Argument] {
	return a.stream
}
