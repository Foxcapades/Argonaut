{{ $vars := (map "ImplType" "Tree" "OutputType" "argo.TreeCommand") -}}
package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Tree struct {
	{{ template "CommandBaseProps" $vars }}
	{{ template "ParentCommonProps" $vars }}
	options argo.TreeCommandOptions
}

func (_ *Tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (i *Tree) SelectedCommand() argo.LeafCommand {
	child := i.selectedChild

	for child != nil {
		if leaf, ok := child.(argo.LeafCommand); ok {
			return leaf
		}

		child = child.(argo.ParentNode).SelectedChild()
	}

	return nil
}

{{ template "ParentCommonFuncs" $vars }}

{{ template "CommandBaseFuncs" $vars }}

func (i *Tree) Options() argo.TreeCommandOptions {
	return i.options
}

func (i *Tree) FindShortFlag(b byte) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Tree) FindLongFlag(name string) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Tree) FindShortFlagRecursive(c byte) argo.Flag {
	if selected := i.SelectedCommand(); selected != nil {
		return selected.FindShortFlagRecursive(c)
	}

	return nil
}

func (i *Tree) FindLongFlagRecursive(name string) argo.Flag {
	if selected := i.SelectedCommand(); selected != nil {
		return selected.FindLongFlagRecursive(name)
	}

	return nil
}
