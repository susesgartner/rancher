package permutationdata

import (
	"github.com/rancher/shepherd/extensions/configoperations"
	"github.com/rancher/shepherd/extensions/configoperations/permutations"
)

const (
	clusterConfigKey = "clusterConfig"
	nodeProvidersKey = "nodeProviders"
	providerKey      = "providers"
	k8sVersionKey    = "kubernetesVersion"
	cniKey           = "cni"
)

func CreateK8sPermutation(config map[string]any) (permutations.Permutation, error) {
	k8sKeyPath := []string{clusterConfigKey, k8sVersionKey}
	k8sKeyValue, err := configoperations.GetValue(k8sKeyPath, config)
	k8sPermutation := permutations.CreatePermutation(k8sKeyPath, k8sKeyValue.([]any), nil)

	return k8sPermutation, err
}

func CreateProviderPermutation(config map[string]any) (permutations.Permutation, error) {
	providerKeyPath := []string{clusterConfigKey, providerKey}
	providerKeyValue, err := configoperations.GetValue(providerKeyPath, config)
	providerPermutation := permutations.CreatePermutation(providerKeyPath, providerKeyValue.([]any), nil)

	return providerPermutation, err
}

func CreateCNIPermutation(config map[string]any) (permutations.Permutation, error) {
	cniKeyPath := []string{clusterConfigKey, cniKey}
	cniKeyValue, err := configoperations.GetValue(cniKeyPath, config)
	cniPermutation := permutations.CreatePermutation(cniKeyPath, cniKeyValue.([]any), nil)

	return cniPermutation, err
}
