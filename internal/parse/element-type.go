package parse

type ElementType uint8

const (
	// ElementTypeLongFlagPair represents a longform flag with a value directly
	// attached by use of an equals ('=') character.
	ElementTypeLongFlagPair ElementType = iota

	// ElementTypeLongFlagSolo represents a longform flag that did not have a value
	// directly attached by use of an equals ('=') character.
	ElementTypeLongFlagSolo

	// ElementTypeShortBlockSolo represents a group of one or more characters
	// following a single dash ('-') character with no value directly attached via
	// an equals ('=') character.
	ElementTypeShortBlockSolo

	// ElementTypeShortBlockPair represents a group of one or more characters
	// following a single dash ('-') character with a value directly attached via
	// an equals ('=') character.
	ElementTypeShortBlockPair

	// ElementTypePlainText represents a plain-text argument that has no flag
	// indicators
	ElementTypePlainText

	ElementTypeBoundary

	ElementTypeEnd
)
