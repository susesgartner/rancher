package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	nodeProvider         = "ec2"
	machineConfigsKey    = "awsMachineConfigs"
	machineConfigKey     = "awsMachineConfig"
	credentialsConfigKey = "awsCredentials"
	provider             = "aws"
)

/*
func LoadAWSRelationships(s *suite.Suite, testConfig map[string]any) []permutation.Relationship {
	credentialsKeyPath := []string{credentialsConfigKey}
	credentialsValue, err := permutation.GetKeyPathValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	credentials := permutation.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{machineConfigsKey}
	machineConfigsValue, err := permutation.GetKeyPathValue(machineConfigsKeyPath, testConfig)
	require.NoError(s.T(), err)

	machineConfigs := permutation.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	return []permutation.Relationship{credentials, machineConfigs}
}*/

func createAMIPermutation(config map[string]any) (permutation.Permutation, error) {
	amiKeyPath := []string{machineConfigsKey, machineConfigKey, "ami"}
	amiKeyValue, err := permutation.GetKeyPathValue(amiKeyPath, config)
	amiPermutation := permutation.CreatePermutation(amiKeyPath, amiKeyValue.([]any), nil)

	return amiPermutation, err
}

func CreateAMIRelationship(config map[string]any) (permutation.Relationship, error) {
	amiPermutation, err := createAMIPermutation(config)
	amiRelationship := permutation.CreateRelationship(provider, nil, nil, []permutation.Permutation{amiPermutation})

	return amiRelationship, err
}
