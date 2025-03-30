package utils_test

import (
	"errors"
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/utils"
)

func TestMust(t *testing.T) {
	defer func() { recover() }()

	utils.Must(errors.New("error"))

	t.Error("expected Must to panic but it didn't")
}

func TestMustReturn(t *testing.T) {
	defer func() { recover() }()

	utils.MustReturn(3, errors.New("error"))

	t.Error("expected MustReturn to panic but it didn't")
}

func TestMustReturn2(t *testing.T) {
	v := utils.MustReturn(3, nil)
	if v != 3 {
		t.Error("expected output to match input but it didn't")
	}
}
