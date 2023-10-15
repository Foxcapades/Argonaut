package event

import "github.com/Foxcapades/Argonaut/v1/internal/chars"

type Consumer = func(Event)

type Event struct {
	Kind Kind
	Data string
}

func endEvent() Event {
	return Event{Kind: KindEnd, Data: chars.StrEmpty}
}

func breakEvent() Event {
	return Event{Kind: KindBreak, Data: chars.StrEmpty}
}

func dashEvent() Event {
	return Event{Kind: KindDash, Data: chars.StrDash}
}

func textEvent(txt string) Event {
	return Event{Kind: KindText, Data: txt}
}

func equalsEvent() Event {
	return Event{Kind: KindEquals, Data: chars.StrEquals}
}
