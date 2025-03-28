package flag

import (
	"fmt"
	"iter"

	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewQueue() Queue {
	return Queue{
		ordered:  make([]utils.Pair[byte, string], 0, 4),
		distinct: make(map[utils.Pair[byte, string]]argo.Flag, 4),
	}
}

type Queue struct {
	ordered  []utils.Pair[byte, string]
	distinct map[utils.Pair[byte, string]]argo.Flag
}

func (q *Queue) Append(f argo.Flag) {
	var key = utils.Pair[byte, string]{L: f.ShortForm(), R: f.LongForm()}
	if _, ok := q.distinct[key]; !ok {
		q.distinct[key] = f
		q.ordered = append(q.ordered, key)
	}
}

func (q *Queue) Iterator() iter.Seq[argo.Flag] {
	return func(yield func(argo.Flag) bool) {
		for _, pair := range q.ordered {
			yield(q.distinct[pair])
		}
	}
}

func (q *Queue) Size() int {
	return len(q.ordered)
}

func (q *Queue) Clear() {
	*q = NewQueue()
}

//

func NewDefaultGroupBuilder() argo.FlagGroupBuilder {
	return NewGroupBuilder(DefaultFlagGroupName)
}

func IsDefaultGroup(group argo.FlagGroupBuilder) bool {
	return group.Name() == DefaultFlagGroupName
}

//

func PrintFlagNames(flag argo.Flag) string {
	if flag.HasLongForm() {
		if flag.HasShortForm() {
			return fmt.Sprintf("-%c | --%s", flag.ShortForm(), flag.LongForm())
		}

		return fmt.Sprintf("--%s", flag.LongForm())
	}

	return fmt.Sprintf("-%c", flag.ShortForm())
}

func UniqueFlagNames(groups []argo.FlagGroupBuilder, errs argo.MultiError) {
	longs := make(map[string]uint8, len(groups))
	shorts := make(map[byte]uint8, len(groups))

	for _, group := range groups {
		for _, flag := range group.Flags() {
			if flag.HasLongForm() {
				longs[flag.LongForm()]++
				if longs[flag.LongForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag longform name %s", flag.LongForm()))
				}
			}

			if flag.HasShortForm() {
				shorts[flag.ShortForm()]++
				if shorts[flag.ShortForm()] == 2 {
					errs.AppendError(fmt.Errorf("conflicting flag shortform character %c", flag.ShortForm()))
				}
			}
		}
	}
}
