package argo

type event struct {
	Kind eventKind
	Data string
}

type eventKind uint8

const (
	eventKindDash eventKind = iota
	eventKindText
	eventKindEquals
	eventKindBreak
	eventKindEnd
)

func (k eventKind) String() string {
	switch k {
	case eventKindDash:
		return "dash"
	case eventKindText:
		return "text"
	case eventKindEquals:
		return "equals"
	case eventKindBreak:
		return "break"
	case eventKindEnd:
		return "end"
	default:
		return "invalid"
	}
}

func endEvent() event {
	return event{Kind: eventKindEnd, Data: strEmpty}
}

func breakEvent() event {
	return event{Kind: eventKindBreak, Data: strEmpty}
}

func dashEvent() event {
	return event{Kind: eventKindDash, Data: strDash}
}

func textEvent(txt string) event {
	return event{Kind: eventKindText, Data: txt}
}

func equalsEvent() event {
	return event{Kind: eventKindEquals, Data: strEquals}
}

func newEmitter(args []string) emitter {
	return emitter{arguments: args, next: newDeque[event](6), argumentIndex: 1}
}

type emitter struct {
	arguments     []string
	argumentIndex int
	next          deque[event]
}

func (e *emitter) Next() event {
	if e.next.IsEmpty() {
		e.queueNext()
	}

	return e.next.Poll()
}

func (e *emitter) queueNext() {
	if e.argumentIndex >= len(e.arguments) {
		e.next.Offer(endEvent())
	} else {
		e.scan(e.arguments[e.argumentIndex])
		e.argumentIndex++
	}
}

func (e *emitter) scan(arg string) {
	// If the argument is an empty string.
	if len(arg) == 0 {
		e.next.Offer(textEvent(arg))
		e.next.Offer(breakEvent())
		return
	}

	// Consume any leading dash characters and pass them up.
	i := 0
	for arg[i] == charDash {
		e.next.Offer(dashEvent())
		i++
	}

	// Substring our argument
	sub := arg[i:]

	// Find the next equals character
	ne := nextEquals(sub)

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
