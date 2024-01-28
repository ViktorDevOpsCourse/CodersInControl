package clusters

// TODO move it to config file
var excludeNamespaces = map[string]struct{}{
	"flux-system":     struct{}{},
	"kube-system":     struct{}{},
	"kube-public":     struct{}{},
	"kube-node-lease": struct{}{},
}

type Namespace struct {
	Name string
}

func (n *Namespace) GetName() string {
	return n.Name
}
