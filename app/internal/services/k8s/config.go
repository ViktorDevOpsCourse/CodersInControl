package k8s

type Config struct {
	Clusters map[string]string // map[cluster_name]config_file_path
}
