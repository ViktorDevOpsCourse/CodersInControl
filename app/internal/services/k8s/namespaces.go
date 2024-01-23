package k8s

type Namespace struct {
	Name        string
	Deployments []Deployment
}

type Deployment struct {
	Name string
}
