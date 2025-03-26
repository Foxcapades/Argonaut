package flag_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
)

// no flag forms set
func TestFlagBuilder_Build01(t *testing.T) {
	_, err := flag.Build(flag.NewBuilder())

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid short flag character
func TestFlagBuilder_Build02(t *testing.T) {
	_, err := flag.Build(flag.NewBuilder().WithShortForm('-'))

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid long flag name 1
func TestFlagBuilder_Build03(t *testing.T) {
	_, err := flag.Build(flag.NewBuilder().WithLongForm("@@@@"))

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// invalid long flag name 1
func TestFlagBuilder_Build04(t *testing.T) {
	_, err := flag.Build(flag.NewBuilder().WithLongForm("a@@@"))

	if err == nil {
		t.Error("expected err to not be nil but it was")
	}
}

// busted-ass argument
func TestFlagBuilder_Build05(t *testing.T) {
	_, err := flag.Build(flag.NewBuilder().
		WithLongForm("test").
		WithBindingAndDefault(3, 4, true))

	if err == nil {
		t.Error("expected err to not have been nil, but it was")
	}

	t.Log(err)
}
