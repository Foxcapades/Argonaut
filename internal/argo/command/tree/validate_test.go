package tree_test

import (
	"testing"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/tree"
)

func TestValidateCommandNodeName(t *testing.T) {
	if err := tree.ValidateNodeName(""); err == nil {
		t.Error("expected an error but didn't get one")
	}

	if err := tree.ValidateNodeName("-who"); err == nil {
		t.Error("expected an error but didn't get one")
	}

	if err := tree.ValidateNodeName("_who"); err != nil {
		t.Error("expected no error but got one")
	}

	if err := tree.ValidateNodeName("who\n"); err == nil {
		t.Error("expected an error but didn't get one")
	}
}
