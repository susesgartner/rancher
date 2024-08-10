package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	azuremachineConfigsKey    = "azureMachineConfigs"
	azureCredentialsConfigKey = "azureCredentials"
)

func LoadAzureRelationships(s *suite.Suite, testConfig map[string]any) []permutation.Relationship {
	credentialsKeyPath := []string{azureCredentialsConfigKey}
	credentialsValue, err := permutation.GetKeyPathValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	credentials := permutation.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{azuremachineConfigsKey}
	machineConfigsValue, err := permutation.GetKeyPathValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	machineConfigs := permutation.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	return []permutation.Relationship{credentials, machineConfigs}
}
