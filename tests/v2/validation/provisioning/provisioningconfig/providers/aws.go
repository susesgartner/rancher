package providers

import (
	"github.com/rancher/rancher/tests/v2/validation/provisioning/provisioningconfig"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/provisioningconfig/clusterconfig"
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	awsNodeProvider      = "ec2"
	machineConfigsKey    = "awsMachineConfigs"
	credentialsConfigKey = "awsCredentials"
)

func LoadAWSRelationships(testConfig map[string]any) []permutation.Relationship {
	credentialsConfig := testConfig[credentialsConfigKey]
	credentials := permutation.CreateRelationship(AWSName, []string{credentialsConfigKey}, credentialsConfig)

	machineConfigsConfig := testConfig[machineConfigsKey]
	machineConfigs := permutation.CreateRelationship(AWSName, []string{machineConfigsKey}, machineConfigsConfig)

	nodeProvider := permutation.CreateRelationship(AWSName, []string{provisioningconfig.ConfigEnvironmentKey, clusterconfig.NodeProvidersKey}, awsNodeProvider)

	return []permutation.Relationship{credentials, machineConfigs, nodeProvider}
}
