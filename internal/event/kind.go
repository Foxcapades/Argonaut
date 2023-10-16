package event

type Kind uint8

const (
	KindDash Kind = iota
	KindText
	KindEquals
	KindBreak
	KindEnd
)

func (k Kind) String() string {
	switch k {
	case KindDash:
		return "dash"
	case KindText:
		return "text"
	case KindEquals:
		return "equals"
	case KindBreak:
		return "break"
	case KindEnd:
		return "end"
	default:
		return "invalid"
	}
}
