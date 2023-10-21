package util_test

import (
	"errors"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/util"
)

func TestMust(t *testing.T) {
	defer func() { recover() }()

	util.Must(errors.New("error"))

	t.Error("expected Must to panic but it didn't")
}

func TestMustReturn(t *testing.T) {
	defer func() { recover() }()

	util.MustReturn(3, errors.New("error"))

	t.Error("expected MustReturn to panic but it didn't")
}

func TestMustReturn2(t *testing.T) {
	v := util.MustReturn(3, nil)
	if v != 3 {
		t.Error("expected output to match input but it didn't")
	}
}
