package k8s

import (
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

type K8S struct {
	clusters map[string]Cluster
}

func NewK8SService(cfg Config) *K8S {
	log := logger.FromDefaultContext()

	k8s := &K8S{
		clusters: make(map[string]Cluster),
	}

	for envName, kubeConfig := range cfg.Clusters {
		client, err := NewClient(kubeConfig)
		if err != nil {
			log.Errorf("Failed connection to `%s` cluster. Err `%s`", envName, err)
			continue
		}
		cluster := NewCluster(client)
		k8s.clusters[envName] = cluster
	}

	if len(k8s.clusters) == 0 {
		log.Fatal("Can't connect to no one cluster")
		return k8s
	}

	return k8s
}
