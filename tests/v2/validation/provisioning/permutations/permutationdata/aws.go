package permutationdata

import (
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	awsNodeProvider         = "ec2"
	awsMachineConfigsKey    = "awsMachineConfigs"
	awsCredentialsConfigKey = "awsCredentials"
)

func LoadAWSRelationships(s *suite.Suite, testConfig map[string]any) []permutation.Relationship {
	credentialsKeyPath := []string{awsCredentialsConfigKey}
	credentialsValue, err := permutation.GetKeyPathValue(credentialsKeyPath, testConfig)
	require.NoError(s.T(), err)

	credentials := permutation.CreateRelationship(AWSName, credentialsKeyPath, credentialsValue, nil)

	machineConfigsKeyPath := []string{awsMachineConfigsKey}
	machineConfigsValue, err := permutation.GetKeyPathValue(machineConfigsKeyPath, testConfig)
	require.NoError(s.T(), err)

	machineConfigs := permutation.CreateRelationship(AWSName, machineConfigsKeyPath, machineConfigsValue, nil)

	nodeProvider := permutation.CreateRelationship(AWSName, []string{ClusterConfigKey, NodeProvidersKey}, awsNodeProvider, nil)

	return []permutation.Relationship{credentials, machineConfigs, nodeProvider}
}
