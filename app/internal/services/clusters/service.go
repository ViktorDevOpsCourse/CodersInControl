package clusters

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"sync"
)

type K8S struct {
	clusters map[string]*Cluster // map[environment]Cluster
	sync.RWMutex
}

func NewK8SService(cfg Config,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository) *K8S {
	log := logger.FromDefaultContext()

	k8s := &K8S{
		clusters: make(map[string]*Cluster),
	}

	for envName, kubeConfig := range cfg.Clusters {
		client, err := NewClient(kubeConfig)
		if err != nil {
			log.Errorf("Failed connection to `%s` cluster. Err `%s`", envName, err)
			continue
		}
		cluster := NewCluster(client, appsStatesStorage, appsEventsStorage)
		cluster.Run()
		k8s.clusters[envName] = cluster
	}

	if len(k8s.clusters) == 0 {
		log.Fatal("Can't connect to no one cluster")
		return k8s
	}

	return k8s
}

func (k *K8S) GetClustersCopy() map[string]Cluster {

	k.RLock()
	defer k.RUnlock()

	copyClusters := make(map[string]Cluster)
	for env, cluster := range k.clusters {
		copyClusters[env] = Cluster{
			client:       cluster.client,
			Applications: cluster.Applications,
			Namespaces:   cluster.Namespaces,
			Controller:   cluster.Controller,
		}
	}

	return copyClusters
}

func (k *K8S) GetCluster(env string) (*Cluster, error) {
	if cluster, ok := k.clusters[env]; ok {
		return cluster, nil
	}
	return nil, fmt.Errorf("cluster for env `%s` not found", env)
}
