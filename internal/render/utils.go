package render

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
)

func Pad(size int, out *utils.BatchWriter) {
	for i := 0; i < size; i++ {
		out.WriteByte(text.SpaceByte)
	}
}

const (
	FlagDivider    = " | "
	ParagraphBreak = "\n\n"
)

var HeaderPadding = [...]string{
	"",
	"  ",
	"    ",
	"      ",
}

var SubLinePadding = [...]string{
	"  ",
	"    ",
	"      ",
	"        ",
}

var DescriptionPadding = [...]string{
	"    ",
	"      ",
	"        ",
	"          ",
}

func IsBreakChar(b byte) bool {
	// TODO: '-' needs to be handled differently than spaces
	//       spaces are removed from the output, however the
	//       dash should be maintained.
	return b == text.SpaceByte || b == text.TabByte
}

type Break = [2]int

// NewDescriptionFormatter returns a configured DescriptionFormatter instance.
//
// TODO: get rid of the constructor and just expose the format function.  There
//
//	is no need for the caller to know about the config struct.
func NewDescriptionFormatter(
	prefixPadding string,
	maxTotalWidth int,
	writer *utils.BatchWriter,
) *DescriptionFormatter {
	maxLineWidth := maxTotalWidth - len(prefixPadding)

	if maxLineWidth < 1 {
		panic(fmt.Errorf(
			"illegal argument combination: cannot break text into lines of width"+
				" %d.  Given width was %d minus prefix length %d.",
			maxLineWidth,
			maxTotalWidth,
			len(prefixPadding),
		))
	}

	return &DescriptionFormatter{
		prefixPadding: prefixPadding,
		maxLineWidth:  maxLineWidth,
		writer:        writer,
	}
}

// DescriptionFormatter is a type that may be used to format a CLI element
// description by breaking it down to fit into the allowed maximum line width.
type DescriptionFormatter struct {
	consecutiveLineFeeds int
	currentLineWidth     int
	maxLineWidth         int
	prefixPadding        string
	writer               *utils.BatchWriter
}

func (d *DescriptionFormatter) writeLineFeed() {
	d.currentLineWidth = 0
	d.writer.WriteByte(text.LineFeedByte)
}

func (d *DescriptionFormatter) writePadding() {
	d.writer.WriteString(d.prefixPadding)
}

func (d *DescriptionFormatter) breakLine() {
	d.writeLineFeed()
	d.writePadding()
}

// Format is the standard entrypoint for the description formatter.
//
// It takes a description string and formats it to fit within the configured
// maximum line length.
func (d *DescriptionFormatter) Format(text string) {
	if len(text) == 0 {
		return
	}

	scanner := breakScanner{text: text}
	lastSegment := segment{}

	for scanner.hasNext() {

		currentSegment := scanner.next()

		switch currentSegment.Type {
		case segmentTypeBreak:
			d.BreakFormatTypeBreak(&lastSegment, &currentSegment)

		case segmentTypeLineBreak:
			d.BreakFormatTypeLineBreak(&lastSegment, &currentSegment)

		case segmentTypeWord:
			d.BreakFormatTypeWord(&lastSegment, &currentSegment)

		default:
			panic("illegal state: unrecognized segment type")
		}
	}
}

func (d *DescriptionFormatter) BreakFormatTypeWord(last, current *segment) {
	switch last.Type {
	case segmentTypeBreak:
		lastAndCurWidth := len(last.Data) + len(current.Data)

		// If the break and the word can fit neatly onto the line
		if d.currentLineWidth+lastAndCurWidth <= d.maxLineWidth {
			d.writer.WriteString(last.Data)
			d.writer.WriteString(current.Data)
			*last = *current
			d.currentLineWidth += lastAndCurWidth
			return
		}

		if d.currentLineWidth+len(last.Data) <= d.maxLineWidth/3*2 {
			d.writer.WriteString(last.Data)
			*last = *current
			d.BreakFormatSplitWord(current.Data)
			return
		}

		d.breakLine()

		*last = *current
		d.BreakFormatSplitWord(current.Data)
		return

	case segmentTypeLineBreak:
		// If there were 2 or more line breaks, write out another one, but eat the
		// rest.
		if d.consecutiveLineFeeds > 1 {
			d.writeLineFeed()
		}

		// Write out the line break.
		d.breakLine()
		d.consecutiveLineFeeds = 0
		*last = *current
		d.BreakFormatSplitWord(current.Data)
		return

	case segmentTypeStart:
		d.writePadding()
		*last = *current
		d.BreakFormatSplitWord(current.Data)
		return

	default:
		panic(fmt.Errorf("illegal state: unexpected segment type %s", last.Type))
	}
}

func (d *DescriptionFormatter) BreakFormatTypeBreak(last, current *segment) {
	switch last.Type {
	// If the previous segment was the initial state, then we are going to
	// silently ignore this segment because we eat leading spaces.
	case segmentTypeStart:
		// Do nothing

	// If the previous segment was a word, then store this segment off as the last
	// segment.
	case segmentTypeWord:
		*last = *current

	// If the previous segment was a break, we are going to silently ignore this
	// segment because we eat leading spaces.
	case segmentTypeLineBreak:
		d.consecutiveLineFeeds = 0

	default:
		panic(fmt.Errorf("illegal state: unexpected segment type %s", last.Type))
	}
}

func (d *DescriptionFormatter) BreakFormatTypeLineBreak(last, current *segment) {
	switch last.Type {
	case segmentTypeStart:
		// eat the line break

	case segmentTypeWord:
		*last = *current
		d.consecutiveLineFeeds++

	case segmentTypeBreak:
		*last = *current
		d.consecutiveLineFeeds++

	case segmentTypeLineBreak:
		*last = *current
		d.consecutiveLineFeeds++

	default:
		panic("illegal state: invalid segment type")
	}
}

// BreakFormatSplitWord splits the given word to fit onto the line if necessary.
//
// If the word can fit onto the line without breaking, then it will be written.
//
// If the word is long enough that it needs to be split to fit onto the line, it
// will be split.  It will be split as many times as necessary to write it out
// while keeping in the max line width.
//
// The `curWidth` variable will be updated with the current width of the line
// after writing is completed.
func (d *DescriptionFormatter) BreakFormatSplitWord(word string) {
	breakAfter := d.maxLineWidth - d.currentLineWidth - 1

	// If breakAfter is longer than 1/3 of the line length, then we will break the
	// word.  If breakAfter is less than 1/3 of the line length, then we will
	// insert a newline first.
	if breakAfter > d.maxLineWidth/3 {

		// If the whole word can fit onto the line without breaking it up, then we
		// can write it out in its entirety.
		if breakAfter+1 >= len(word) {
			d.writer.WriteString(word)

			// Update the current line width to reflect our addition.
			d.currentLineWidth += len(word)

			// Return because no further work is necessary.
			return
		}

		// So the word is long enough that it can't fit onto the line without being
		// broken up.

		// Write out as much of the word as we can fit.
		d.writer.WriteString(word[:breakAfter])

		// Write out our hyphen
		d.writer.WriteByte(text.DashByte)

		// Chop the word down by the amount that we've written so far.
		word = word[breakAfter:]
	} else if breakAfter == 0 {
		d.writer.WriteByte(word[0])
		word = word[1:]
	}

	if len(word) == 0 {
		return
	}

	breakAfter = d.maxLineWidth - 1

	for {
		// Write out the line feed
		d.writeLineFeed()

		// write out the prefix padding
		d.writer.WriteString(d.prefixPadding)

		// If the whole remaining word can fit into the new line, then we can write
		// it without breaking it any further.
		if breakAfter+1 >= len(word) {
			// Write out the whole word
			d.writer.WriteString(word)

			// update the line width
			d.currentLineWidth = len(word)

			// bail out here
			return
		}

		// So the word is still long enough to need to be broken further.

		if breakAfter > 0 {
			// Write out as much of the word as we can.
			d.writer.WriteString(word[:breakAfter])
			// Write out the hyphen
			d.writer.WriteByte(text.DashByte)

			// Chop down the word even further.
			word = word[breakAfter:]
		} else {
			d.writer.WriteString(word[:breakAfter+1])

			// Chop down the word even further.
			word = word[breakAfter+1:]
		}
	}
}

// +------------------------------------------------------------------------+ //
// |                                                                        | //
// |  Description Break Scanner.                                            | //
// |                                                                        | //
// +------------------------------------------------------------------------+ //

type segmentType uint8

const (
	segmentTypeStart segmentType = iota
	segmentTypeBreak
	segmentTypeLineBreak
	segmentTypeWord
)

func (s segmentType) String() string {
	switch s {
	case segmentTypeStart:
		return "start"
	case segmentTypeBreak:
		return "break"
	case segmentTypeLineBreak:
		return "line-break"
	case segmentTypeWord:
		return "word"
	default:
		return "invalid"
	}
}

type segment struct {
	Type segmentType
	Data string
}

type breakScanner struct {
	text string
	pos  int
}

func (b *breakScanner) hasNext() bool {
	return b.pos < len(b.text)
}

func (b *breakScanner) next() segment {
	// Record the current starting position in the source data.
	start := b.pos

	// If the current character is a breakable character...
	if IsBreakChar(b.text[b.pos]) {
		// move onto the next character
		b.pos++

		// while there are more characters available (may be false on first hit due
		// to the b.pos++ above)...
		for b.hasNext() {
			// If the character is not a breakable character, halt the iteration
			if !IsBreakChar(b.text[b.pos]) {
				break
			}

			// if the character _is_ a breakable character, move forward another char
			b.pos++
		}

		// The segment type is break as it contains nothing but breakable
		// characters.
		return segment{segmentTypeBreak, b.text[start:b.pos]}
	}

	// If the next character is a CR...
	if b.text[b.pos] == text.CarriageReturnByte {
		// and there is a following LF...
		if b.hasNext() && b.text[b.pos+1] == text.LineFeedByte {
			// bump up the pos by 2
			b.pos += 2
		} else {
			// bump up the pos
			b.pos++
		}

		return segment{segmentTypeBreak, text.LineFeedString}
	}

	// If the next character is an LF
	if b.text[b.pos] == text.LineFeedByte {
		b.pos++

		return segment{segmentTypeLineBreak, text.LineFeedString}
	}

	b.pos++

	if b.hasNext() {
		for b.hasNext() {
			if text.IsWhitespace(b.text[b.pos]) {
				break
			}
			b.pos++
		}
	}

	return segment{segmentTypeWord, b.text[start:b.pos]}
}
