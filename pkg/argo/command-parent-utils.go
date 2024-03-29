package argo

import (
	"os"

	"github.com/Foxcapades/Argonaut/internal/util"
)

func defaultOnIncompleteHandler(parent CommandParent) {
	if tree, ok := parent.(CommandTree); ok {
		util.Must(comTreeRenderer{}.RenderHelp(tree, os.Stderr))
	} else if branch, ok := parent.(CommandBranch); ok {
		util.Must(comBranchRenderer{}.RenderHelp(branch, os.Stderr))
	} else {
		panic("illegal state: unrecognized command parent implementation")
	}

	os.Exit(1)
}
