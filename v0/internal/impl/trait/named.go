package trait

type Named struct {
	NameTxt string
}

func (n *Named) Name() string {
	return n.NameTxt
}

func (n *Named) HasName() bool {
	return len(n.NameTxt) > 0
}
