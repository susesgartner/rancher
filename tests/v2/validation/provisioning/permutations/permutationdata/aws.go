package permutationdata

import (
	"github.com/rancher/shepherd/extensions/configoperations"
	"github.com/rancher/shepherd/extensions/configoperations/permutations"
)

const (
	nodeProvider         = "ec2"
	machineConfigsKey    = "awsMachineConfigs"
	machineConfigKey     = "awsMachineConfig"
	credentialsConfigKey = "awsCredentials"
	provider             = "aws"
)

/*
func LoadAWSRelationships(s *suite.Suite, testConfig map[string]any) []permutations.Relationship {
	credentialsKeyPath := []string{credentialsConfigKey}
	credentialsValue, err := permutations.GetKeyPathValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	credentials := permutations.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{machineConfigsKey}
	machineConfigsValue, err := permutations.GetKeyPathValue(machineConfigsKeyPath, testConfig)
	require.NoError(s.T(), err)

	machineConfigs := permutations.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	return []permutations.Relationship{credentials, machineConfigs}
}*/

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
