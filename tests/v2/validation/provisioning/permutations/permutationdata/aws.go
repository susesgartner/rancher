package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	awsNodeProvider      = "ec2"
	machineConfigsKey    = "awsMachineConfigs"
	credentialsConfigKey = "awsCredentials"
)

func LoadAWSRelationships(testConfig map[string]any) []permutation.Relationship {
	credentialsConfig := testConfig[credentialsConfigKey]
	credentials := permutation.CreateRelationship(AWSName, []string{credentialsConfigKey}, credentialsConfig, nil)

	machineConfigsConfig := testConfig[machineConfigsKey]
	machineConfigs := permutation.CreateRelationship(AWSName, []string{machineConfigsKey}, machineConfigsConfig, nil)

	nodeProvider := permutation.CreateRelationship(AWSName, []string{ClusterConfigKey, NodeProvidersKey}, awsNodeProvider, nil)

	return []permutation.Relationship{credentials, machineConfigs, nodeProvider}
}
