package argo

type CommandCallback[T any] = func(node T)

// IncompleteCommandHandler defines a function type that may be used as a callback
// for when a command leaf is not reached when parsing a command tree structure.
type IncompleteCommandHandler[T any] = func(command T)
