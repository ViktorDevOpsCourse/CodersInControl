package clusters

type Namespace struct {
	Name string
}

func (n *Namespace) GetName() string {
	return n.Name
}
