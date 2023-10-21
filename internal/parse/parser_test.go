package parse_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/internal/emit"
	"github.com/Foxcapades/Argonaut/internal/parse"
)

func TestParser_Next01(t *testing.T) {
	input := []string{"com", ""}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypePlainText {
		t.Error("expected element type to be plain text but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected element data to have exactly 1 element but it didn't")
	} else if next.Data[0] != "" {
		t.Error("expected element data to be an empty string but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next02(t *testing.T) {
	input := []string{"com", "-t"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypeShortBlockSolo {
		t.Error("expected element type to be short block solo but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected element data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "t" {
		t.Error("expected element data to be 't' but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next03(t *testing.T) {
	input := []string{"com", "-t=1"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypeShortBlockPair {
		t.Error("expected element type to be short block pair but it wasn't")
	}

	if len(next.Data) != 2 {
		t.Error("expected element data length to be exactly 2 but it wasn't")
	} else {
		if next.Data[0] != "t" {
			t.Error("expected element data 0 to be 't' but it wasn't")
		} else if next.Data[1] != "1" {
			t.Error("expected element data 1 to be '1' but it wasn't")
		}
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next04(t *testing.T) {
	input := []string{"com", "--trumpet"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypeLongFlagSolo {
		t.Error("expected element type to be long flag solo but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected element data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "trumpet" {
		t.Error("expected element data to be 'trumpet' but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next05(t *testing.T) {
	input := []string{"com", "--bandannas=56"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypeLongFlagPair {
		t.Error("expected element type to be long flag pair but it wasn't")
	}

	if len(next.Data) != 2 {
		t.Error("expected element data length to be exactly 2 but it wasn't")
	} else {
		if next.Data[0] != "bandannas" {
			t.Error("expected element data 0 to be 'bandannas' but it wasn't")
		} else if next.Data[1] != "56" {
			t.Error("expected element data 1 to be '56' but it wasn't")
		}
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next06(t *testing.T) {
	input := []string{"com", "--whoa, partner"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypePlainText {
		t.Error("expected element type to be plain text but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected element data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "--whoa, partner" {
		t.Error("expected element data 0 to be '--whoa, partner' but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next07(t *testing.T) {
	input := []string{"com", "---"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypePlainText {
		t.Error("expected element type to be plain text but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected element data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "---" {
		t.Error("expected element data 0 to be '---' but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next08(t *testing.T) {
	input := []string{"com", "--", "--what"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypeBoundary {
		t.Error("expected element type to be boundary but it wasn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypePlainText {
		t.Error("expected element type to be plain text but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "--what" {
		t.Error("expected data to match input but it didn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestParser_Next09(t *testing.T) {
	input := []string{"com", "-ab cd"}
	parser := parse.NewParser(emit.NewEmitter(input))

	next := parser.Next()

	if next.Type != parse.ElementTypePlainText {
		t.Error("expected element type to be plain text but it wasn't")
	}

	if len(next.Data) != 1 {
		t.Error("expected data length to be exactly 1 but it wasn't")
	} else if next.Data[0] != "-ab cd" {
		t.Error("expected data to match input but it didn't")
	}

	next = parser.Next()

	if next.Type != parse.ElementTypeEnd {
		t.Error("expected element type to be end but it wasn't")
	}
}

func TestElement_String(t *testing.T) {
	{
		e := parse.Element{Type: parse.ElementTypeLongFlagPair, Data: []string{"foo", "bar"}}
		if e.String() != "--foo=bar" {
			t.Errorf("expected flag to be '--foo=bar' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypeShortBlockPair, Data: []string{"abc", "bar"}}
		if e.String() != "-abc=bar" {
			t.Errorf("expected flag to be '--abc=bar' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypeLongFlagSolo, Data: []string{"goop"}}
		if e.String() != "--goop" {
			t.Errorf("expected flag to be '--goop' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypeShortBlockSolo, Data: []string{"rump"}}
		if e.String() != "-rump" {
			t.Errorf("expected flag to be '-rump' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypeBoundary, Data: []string{"--"}}
		if e.String() != "--" {
			t.Errorf("expected flag to be '--' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypePlainText, Data: []string{"aspartame"}}
		if e.String() != "aspartame" {
			t.Errorf("expected flag to be 'aspartame' but it was '%s'", e.String())
		}
	}
	{
		e := parse.Element{Type: parse.ElementTypeEnd, Data: []string{"\000"}}
		if e.String() != "\000" {
			t.Errorf("expected flag to be '\000' but it was '%s'", e.String())
		}
	}
}
