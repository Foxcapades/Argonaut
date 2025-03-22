{{ $vars := (types "LeafCommand" "LeafCommandBuilder" "Leaf") -}}
package argo

// A LeafCommand is the final node in a CommandTree branch.
//
// Command leaves may be children of either a CommandTree directly, or of a
// BranchCommand.
type LeafCommand interface {
	ChildNode
  {{ template "CommandCommon" $vars }}
	{{- template "CommandEnd" $vars }}
}

// LeafCommandBuilder defines a builder type that is used to construct
// LeafCommand instances.
type LeafCommandBuilder interface {
	ChildNodeBuilder
	{{ template "CommandBuilderCommon" $vars }}
	{{- template "ChildNodeBuilderCommon" $vars }}
	{{- template "CommandEnd" $vars }}
}
