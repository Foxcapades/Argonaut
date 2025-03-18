package argo

// Branch represents a subcommand under a CommandTree that is an
// intermediate node between the tree root and an executable CommandLeaf.
//
// CommandBranches enable the organization of subcommands into categories.
//
// Example command tree:
//
//	docker
//	 |- compose
//	 |   |- build
//	 |   |- down
//	 |   |- ...
//	 |- container
//	 |   |- exec
//	 |   |- ls
//	 |   |- ...
//	 |- ...
type Branch interface {
	ParentNode[Branch]
	ChildNode[Branch]
}

type CommandBranchCallback = func(branch Branch)

// A BranchBuilder instance may be used to configure a new Branch
// instance to be built.
//
// CommandBranches are intermediate steps between the root of the CommandTree
// and the CommandLeaf instances.
//
// For example, given the following command example, the tree is "foo", the
// branch is "bar", and the leaf is "fizz":
//
//	./foo bar fizz
type BranchBuilder interface {
	ChildBuilder[BranchBuilder]
	ParentBuilder[BranchBuilder]

	Build(warnings *WarningContext) (Branch, error)
}
