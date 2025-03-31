package flag

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type FormGroup struct {
	Node  GroupContainer
	Forms []Forms
}

func GroupInheritance(node GroupContainer) []FormGroup {
	shorts := make(map[byte]bool, 16)
	longs := make(map[string]bool, 16)

	currentNode := node

	result := make([]FormGroup, 0, 4)

	for {
		currentForms := FormGroup{
			Node:  currentNode,
			Forms: makeFormsSlice(currentNode, currentNode == node),
		}

		for _, flag := range Stream(currentNode) {
			forms := Forms{Flag: flag}

			if flag.HasLongForm() {
				if _, ok := longs[flag.LongForm()]; !ok {
					longs[flag.LongForm()] = true
					forms.Long = true
				}
			}

			if flag.HasShortForm() {
				if _, ok := shorts[flag.ShortForm()]; !ok {
					shorts[flag.ShortForm()] = true
					forms.Short = true
				}
			}

			// If the flag has an available form, and the current node is not the
			// original (we want them in the filter maps, but not in the output).
			if (forms.Short || forms.Long) && currentNode != node {
				currentForms.Forms = append(currentForms.Forms, forms)
			}
		}

		if currentNode != node && len(currentForms.Forms) > 0 {
			result = append(result, currentForms)
		}

		if t, ok := currentNode.(argo.ChildNode); ok {
			currentNode = t.Parent().(GroupContainer)
		} else {
			break
		}
	}

	return result
}

type named interface{ Name() string }

func RenderGroupedInheritance(formGroup *FormGroup, options opts.Options, padding uint8, sb *utils.BatchWriter) {
	sb.WriteByte(text.LineFeedByte)

	var name string
	if _, ok := formGroup.Node.(argo.TreeCommand); ok {
		name = "Command Root"
	} else {
		name = `Parent Command "` + formGroup.Node.(named).Name() + `"`
	}

	sb.WriteString("\nFlags Inherited from ")
	sb.WriteString(name)

	for i := range formGroup.Forms {
		if i > 0 && !formGroup.Forms[i-1].Flag.HasDescription() {
			sb.WriteByte(text.LineFeedByte)
		}

		sb.WriteByte(text.LineFeedByte)

		RenderInheritedForms(&formGroup.Forms[i], options, padding, sb)
	}
}

func makeFormsSlice(node GroupContainer, isRoot bool) []Forms {
	if isRoot {
		return nil
	}

	size := 0

	for _, g := range node.FlagGroups() {
		size += g.Size()
	}

	return make([]Forms, 0, size)
}
