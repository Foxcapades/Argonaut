package parse

type Element struct {
	Type ElementType
	Data []string
}

func (e Element) String() string {
	switch e.Type {
	case ElementTypeLongFlagPair:
		return "--" + e.Data[0] + "=" + e.Data[1]
	case ElementTypeShortBlockPair:
		return "-" + e.Data[0] + "=" + e.Data[1]
	case ElementTypeLongFlagSolo:
		return "--" + e.Data[0]
	case ElementTypeShortBlockSolo:
		return "-" + e.Data[0]
	case ElementTypeBoundary:
		return "--"
	case ElementTypePlainText:
		return e.Data[0]
	case ElementTypeEnd:
		return string(byte(0))
	default:
		panic("illegal state")
	}
}

func longPairElement(flag, value string) Element {
	return Element{Type: ElementTypeLongFlagPair, Data: []string{flag, value}}
}

func longSoloElement(flag string) Element {
	return Element{Type: ElementTypeLongFlagSolo, Data: []string{flag}}
}

func shortPairElement(flags, value string) Element {
	return Element{Type: ElementTypeShortBlockPair, Data: []string{flags, value}}
}

func shortSoloElement(flags string) Element {
	return Element{Type: ElementTypeShortBlockSolo, Data: []string{flags}}
}

func textElement(text string) Element {
	return Element{Type: ElementTypePlainText, Data: []string{text}}
}

func boundaryElement() Element {
	return Element{Type: ElementTypeBoundary}
}

func endElement() Element {
	return Element{Type: ElementTypeEnd}
}
