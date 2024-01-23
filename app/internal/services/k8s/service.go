package k8s

type K8S struct {
	clusters []Cluster
}

func NewK8SService() *K8S {
	cluster := NewCluster()
	return &K8S{
		clusters: []Cluster{
			cluster,
		},
	}
}
