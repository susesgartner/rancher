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
	credentialsConfig := testConfig[awsCredentialsConfigKey]
	credentials := permutation.CreateRelationship(AWSName, []string{awsCredentialsConfigKey}, credentialsConfig, nil)

	machineConfigsConfig := testConfig[awsMachineConfigsKey]
	machineConfigs := permutation.CreateRelationship(AWSName, []string{awsMachineConfigsKey}, machineConfigsConfig, nil)

	nodeProvider := permutation.CreateRelationship(AWSName, []string{ClusterConfigKey, NodeProvidersKey}, awsNodeProvider, nil)

	return []permutation.Relationship{credentials, machineConfigs, nodeProvider}
}
