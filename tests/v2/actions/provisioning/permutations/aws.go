package permutations

import (
	"github.com/rancher/shepherd/extensions/configoperations"
	"github.com/rancher/shepherd/pkg/config/operations/permutations"
)

const (
	nodeProvider         = "ec2"
	machineConfigsKey    = "awsMachineConfigs"
	machineConfigKey     = "awsMachineConfig"
	credentialsConfigKey = "awsCredentials"
	provider             = "aws"
)

func createAMIPermutation(config map[string]any) (permutations.Permutation, error) {
	amiKeyPath := []string{machineConfigsKey, machineConfigKey, "ami"}
	amiKeyValue, err := configoperations.GetValue(amiKeyPath, config)
	amiPermutation := permutations.CreatePermutation(amiKeyPath, amiKeyValue.([]any), nil)

	return amiPermutation, err
}

func CreateAMIRelationship(config map[string]any) (permutations.Relationship, error) {
	amiPermutation, err := createAMIPermutation(config)
	amiRelationship := permutations.CreateRelationship(provider, nil, nil, []permutations.Permutation{amiPermutation})

	return amiRelationship, err
}
