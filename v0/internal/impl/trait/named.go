package trait

type Named struct {
	NameValue string
}

func (n *Named) Name() string {
	return n.NameValue
}

func (n *Named) HasName() bool {
	return len(n.NameValue) > 0
}
