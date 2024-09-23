package token

type Type uint8

const (
	TypeArgument Type = iota
	TypeShort
	TypeLong
	TypeEndOfOptions
	TypeEndOfInput
)

func (t Type) String() string {
	switch t {

	case TypeArgument:
		return "argument"

	case TypeShort:
		return "short"

	case TypeLong:
		return "long"

	case TypeEndOfOptions:
		return "end-of-options"

	case TypeEndOfInput:
		return "end-of-input"

	}

	return "invalid"
}
