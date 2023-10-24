package util_test

import (
	"testing"

	"github.com/Foxcapades/Argonaut/internal/util"
)

func TestIfElse01(t *testing.T) {
	if util.IfElse(true, 1, 2) != 1 {
		t.Fail()
	}
}

func TestIfElse02(t *testing.T) {
	if util.IfElse(false, 1, 2) != 2 {
		t.Fail()
	}
}
