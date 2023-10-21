package emit_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/internal/emit"
)

func TestEmitter_Next01(t *testing.T) {
	input := []string{"com", ""}
	emitter := emit.NewEmitter(input)

	next := emitter.Next()

	if next.Kind != emit.EventKindText {
		t.Error("expected event type to be text but it wasn't")
	}
	if next.Data != "" {
		t.Error("expected event data to be empty but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindBreak {
		t.Error("expected event type to be break but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindEnd {
		t.Error("expected event type to be end but it wasn't")
	}
}

func TestEmitter_Next02(t *testing.T) {
	input := []string{"com", "-abcd"}
	emitter := emit.NewEmitter(input)

	next := emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindText {
		t.Error("expected event type to be text but it wasn't")
	}
	if next.Data != "abcd" {
		t.Error("expected event data to be abcd but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindBreak {
		t.Error("expected event type to be break but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindEnd {
		t.Error("expected event type to be end but it wasn't")
	}
}

func TestEmitter_Next03(t *testing.T) {
	input := []string{"com", "---"}
	emitter := emit.NewEmitter(input)

	next := emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindBreak {
		t.Error("expected event type to be break but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindEnd {
		t.Error("expected event type to be end but it wasn't")
	}
}

func TestEmitter_Next04(t *testing.T) {
	input := []string{"com", "--foo=bar"}
	emitter := emit.NewEmitter(input)

	next := emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindDash {
		t.Error("expected event type to be dash but it wasn't")
	}
	if next.Data != "-" {
		t.Error("expected event data to be \"-\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindText {
		t.Error("expected event type to be text but it wasn't")
	}
	if next.Data != "foo" {
		t.Error("expected event data to be \"foo\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindEquals {
		t.Error("expected event type to be equals but it wasn't")
	}
	if next.Data != "=" {
		t.Error("expected event data to be \"=\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindText {
		t.Error("expected event type to be text but it wasn't")
	}
	if next.Data != "bar" {
		t.Error("expected event data to be \"bar\" but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindBreak {
		t.Error("expected event type to be break but it wasn't")
	}

	next = emitter.Next()

	if next.Kind != emit.EventKindEnd {
		t.Error("expected event type to be end but it wasn't")
	}
}

func TestEventKind_String(t *testing.T) {
	if emit.EventKindDash.String() != "dash" {
		t.Errorf("expected \"dash\" got \"%s\"", emit.EventKindDash.String())
	}

	if emit.EventKindText.String() != "text" {
		t.Errorf("expected \"text\" got \"%s\"", emit.EventKindText.String())
	}

	if emit.EventKindEquals.String() != "equals" {
		t.Errorf("expected \"equals\" got \"%s\"", emit.EventKindEquals.String())
	}

	if emit.EventKindBreak.String() != "break" {
		t.Errorf("expected \"break\" got \"%s\"", emit.EventKindBreak.String())
	}

	if emit.EventKindEnd.String() != "end" {
		t.Errorf("expected \"end\" got \"%s\"", emit.EventKindEnd.String())
	}

	dummy := emit.EventKind(55)

	if dummy.String() != "invalid" {
		t.Errorf("expected \"invalid\" got \"%s\"", dummy.String())
	}

}
