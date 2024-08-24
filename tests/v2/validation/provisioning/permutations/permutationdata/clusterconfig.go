package permutationdata

import "github.com/rancher/shepherd/extensions/permutation"

const (
	clusterConfigKey = "clusterConfig"
	nodeProvidersKey = "nodeProviders"
	providerKey      = "providers"
	k8sVersionKey    = "kubernetesVersion"
	cniKey           = "cni"
)

func CreateK8sPermutation(config map[string]any) (permutation.Permutation, error) {
	k8sKeyPath := []string{clusterConfigKey, k8sVersionKey}
	k8sKeyValue, err := permutation.GetKeyPathValue(k8sKeyPath, config)
	k8sPermutation := permutation.CreatePermutation(k8sKeyPath, k8sKeyValue.([]any), nil)

	return k8sPermutation, err
}

func CreateProviderPermutation(config map[string]any) (permutation.Permutation, error) {
	providerKeyPath := []string{clusterConfigKey, providerKey}
	providerKeyValue, err := permutation.GetKeyPathValue(providerKeyPath, config)
	providerPermutation := permutation.CreatePermutation(providerKeyPath, providerKeyValue.([]any), nil)

	return providerPermutation, err
}

func CreateCNIPermutation(config map[string]any) (permutation.Permutation, error) {
	cniKeyPath := []string{clusterConfigKey, cniKey}
	cniKeyValue, err := permutation.GetKeyPathValue(cniKeyPath, config)
	cniPermutation := permutation.CreatePermutation(cniKeyPath, cniKeyValue.([]any), nil)

	return cniPermutation, err
}
