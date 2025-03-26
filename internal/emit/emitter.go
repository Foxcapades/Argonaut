package emit

import (
	"strings"

	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
)

type Event struct {
	Kind EventKind
	Data string
}

type EventKind uint8

const (
	EventKindDash EventKind = iota
	EventKindText
	EventKindEquals
	EventKindBreak
	EventKindEnd
)

func (k EventKind) String() string {
	switch k {
	case EventKindDash:
		return "dash"
	case EventKindText:
		return "text"
	case EventKindEquals:
		return "equals"
	case EventKindBreak:
		return "break"
	case EventKindEnd:
		return "end"
	default:
		return "invalid"
	}
}

func endEvent() Event {
	return Event{Kind: EventKindEnd, Data: text.EmptyString}
}

func breakEvent() Event {
	return Event{Kind: EventKindBreak, Data: text.EmptyString}
}

func dashEvent() Event {
	return Event{Kind: EventKindDash, Data: text.DashString}
}

func textEvent(txt string) Event {
	return Event{Kind: EventKindText, Data: txt}
}

func equalsEvent() Event {
	return Event{Kind: EventKindEquals, Data: text.EqualsString}
}

func NewEmitter(args []string) Emitter {
	return Emitter{arguments: args, next: utils.NewDeque[Event](6), argumentIndex: 1}
}

type Emitter struct {
	arguments     []string
	argumentIndex int
	next          utils.Deque[Event]
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
	} else {
		e.scan(e.arguments[e.argumentIndex])
		e.argumentIndex++
	}
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
	for i < len(arg) && arg[i] == text.DashByte {
		e.next.Offer(dashEvent())
		i++
	}

	// If the whole thing was just dashes, then break here.
	if i >= len(arg) {
		e.next.Offer(breakEvent())
		return
	}

	// Substring our argument
	sub := arg[i:]

	// Find the next equals character
	ne := strings.IndexByte(sub, text.EqualsByte)

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
