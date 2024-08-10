package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	awsNodeProvider         = "ec2"
	awsMachineConfigsKey    = "awsMachineConfigs"
	awsCredentialsConfigKey = "awsCredentials"
)

func LoadAWSRelationships(testConfig map[string]any) []permutation.Relationship {
	credentialsKeyPath := []string{awsCredentialsConfigKey}
	credentialsValue, _ := permutation.GetKeyPathValue(credentialsKeyPath, testConfig)
	credentials := permutation.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{awsMachineConfigsKey}
	machineConfigsValue, _ := permutation.GetKeyPathValue(machineConfigsKeyPath, testConfig)
	machineConfigs := permutation.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	nodeProvider := permutation.CreateRelationship(AWSName, []string{ClusterConfigKey, NodeProvidersKey}, awsNodeProvider, nil)

	return []permutation.Relationship{credentials, machineConfigs, nodeProvider}
}
