package permutationdata

import (
	"github.com/rancher/shepherd/extensions/configoperations"
	"github.com/rancher/shepherd/extensions/configoperations/permutations"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	azuremachineConfigsKey    = "azureMachineConfigs"
	azureCredentialsConfigKey = "azureCredentials"
)

func LoadAzureRelationships(s *suite.Suite, testConfig map[string]any) []permutations.Relationship {
	credentialsKeyPath := []string{azureCredentialsConfigKey}
	credentialsValue, err := configoperations.GetValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	credentials := permutations.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{azuremachineConfigsKey}
	machineConfigsValue, err := configoperations.GetValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	machineConfigs := permutations.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	return []permutations.Relationship{credentials, machineConfigs}
}
