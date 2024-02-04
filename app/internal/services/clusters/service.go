package clusters

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"sync"
)

type Clusters map[string]*Cluster

type K8S struct {
	clusters Clusters // map[clusterName]Cluster
	sync.RWMutex
}

func NewK8SService(cfg Config,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository) *K8S {
	log := logger.FromDefaultContext()

	k8s := &K8S{
		clusters: make(Clusters),
	}

	for clusterName, kubeConfig := range cfg.Clusters {
		client, err := NewClient(kubeConfig)
		if err != nil {
			log.Errorf("Failed connection to `%s` cluster. Err `%s`", clusterName, err)
			continue
		}

		cluster := NewCluster(clusterName, client, appsStatesStorage, appsEventsStorage)

		log.Infof("Connected to k8s cluster `%s`", clusterName)

		cluster.Run()
		k8s.clusters[clusterName] = cluster
	}

	if len(k8s.clusters) == 0 {
		log.Fatal("Can't connect to no one cluster")
		return k8s
	}

	return k8s
}

func (k *K8S) GetClustersCopy() ClustersCopy {

	k.RLock()
	defer k.RUnlock()

	copyClusters := make(ClustersCopy)
	for clusterName, cluster := range k.clusters {
		copyClusters[clusterName] = Cluster{
			Applications: cluster.Applications,
			Namespaces:   cluster.Namespaces,
			ClusterName:  cluster.ClusterName,
		}
	}

	return copyClusters
}

func (k *K8S) GetCluster(clusterName string) (*Cluster, error) {
	if cluster, ok := k.clusters[clusterName]; ok {
		return cluster, nil
	}
	return nil, fmt.Errorf("cluster for env `%s` not found", clusterName)
}

type ClustersCopy map[string]Cluster

func (c ClustersCopy) GetCluster(clusterName string) Cluster {
	return c[clusterName]
}
