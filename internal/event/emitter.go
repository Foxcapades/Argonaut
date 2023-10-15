package event

import (
	"github.com/Foxcapades/Argonaut/v1/internal/chars"
	"github.com/Foxcapades/Argonaut/v1/internal/util"
)

func NewEmitter(args []string) Emitter {
	return Emitter{arguments: args, next: util.NewDeque[Event](6)}
}

type Emitter struct {
	arguments     []string
	argumentIndex int
	next          util.Deque[Event]
}

func (e *Emitter) Next() Event {
	if e.next.IsEmpty() {
		e.queueNext()
	}

	return e.next.Poll()
}

func (e *Emitter) queueNext() {
	if e.argumentIndex >= len(e.arguments) {
		e.next.Offer(endEvent())
	}

	e.scan(e.arguments[e.argumentIndex])
	e.argumentIndex++
}

func (e *Emitter) scan(arg string) {
	// If the argument is an empty string.
	if len(arg) == 0 {
		e.next.Offer(textEvent(arg))
		e.next.Offer(breakEvent())
		return
	}

	// Consume any leading dash characters and pass them up.
	i := 0
	for arg[i] == chars.CharDash {
		e.next.Offer(dashEvent())
		i++
	}

	// Substring our argument
	sub := arg[i:]

	// Find the next equals character
	ne := chars.NextEquals(sub)

	// If there is no next equals character, pass up the whole substring
	if ne == -1 {
		e.next.Offer(textEvent(sub))
	} else

	// If there is an equals character, divide up the substring.
	{
		e.next.Offer(textEvent(sub[:ne]))
		e.next.Offer(equalsEvent())
		e.next.Offer(textEvent(sub[ne+1:]))
	}

	// End of argument
	e.next.Offer(breakEvent())
}
