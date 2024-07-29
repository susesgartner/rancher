package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	azuremachineConfigsKey    = "azureMachineConfigs"
	azureCredentialsConfigKey = "azureCredentials"
)

func LoadAzureRelationships(testConfig map[string]any) []permutation.Relationship {
	credentialsConfig := testConfig[azureCredentialsConfigKey]
	credentials := permutation.CreateRelationship(AWSName, []string{azureCredentialsConfigKey}, credentialsConfig, nil)

	machineConfigsConfig := testConfig[azuremachineConfigsKey]
	machineConfigs := permutation.CreateRelationship(AWSName, []string{azuremachineConfigsKey}, machineConfigsConfig, nil)

	return []permutation.Relationship{credentials, machineConfigs}
}
