package utils_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/utils"
)

func TestIfElse01(t *testing.T) {
	if utils.IfElse(true, 1, 2) != 1 {
		t.Fail()
	}
}

func TestIfElse02(t *testing.T) {
	if utils.IfElse(false, 1, 2) != 2 {
		t.Fail()
	}
}
