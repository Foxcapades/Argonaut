package event

type Kind uint8

const (
	KindDash Kind = iota
	KindText
	KindEquals
	KindBreak
	KindEnd
)
