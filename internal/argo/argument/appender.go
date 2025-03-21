package argument

import (
	"iter"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Appender struct {
	stream iter.Seq[argo.Argument]
}

func (a Appender) Append(rawValue string) (bool, error) {
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

func (a Appender) Remaining() iter.Seq[argo.Argument] {
	return a.stream
}
