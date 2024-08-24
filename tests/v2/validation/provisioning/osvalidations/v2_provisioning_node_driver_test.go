package osvalidations

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations/permutationdata"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/rancher/shepherd/extensions/provisioning"
	"github.com/rancher/shepherd/pkg/config"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	provisioningInputKey = "clusterConfig"
	k8sVersionKey        = "kubernetesVersion"
	awsMachineKey        = "awsMachineConfigs"
	namespace            = "fleet-default"
)

type V2ProvisioningNodeDriver struct {
	suite.Suite
	client             *rancher.Client
	session            *session.Session
	standardUserClient *rancher.Client
	permutedConfigs    []map[string]any
	clusters           []v1.SteveAPIObject
}

func (k *V2ProvisioningNodeDriver) TearDownSuite() {
	k.session.Cleanup()
}

func (k *V2ProvisioningNodeDriver) SetupSuite() {
	testSession := session.NewSession()
	k.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(k.T(), err)
	k.client = client

	cattleConfig := config.LoadConfigFromFile(os.Getenv(config.ConfigEnvironmentKey))

	k8sKeyPath := []string{provisioningInputKey, k8sVersionKey}
	k8sKeyValue, err := permutation.GetKeyPathValue(k8sKeyPath, cattleConfig)
	require.NoError(k.T(), err)

	k8sPermutation := permutation.CreatePermutation(k8sKeyPath, k8sKeyValue.([]any), nil)

	amiKeyPath := []string{"awsMachineConfigs", "awsMachineConfig", "ami"}
	amiKeyValue, err := permutation.GetKeyPathValue(amiKeyPath, cattleConfig)
	require.NoError(k.T(), err)

	amiPermutation := permutation.CreatePermutation(amiKeyPath, amiKeyValue.([]any), nil)

	amiRelationship := permutation.CreateRelationship("aws", nil, nil, []permutation.Permutation{amiPermutation})
	providerRelations := permutationdata.LoadProviderRelationships(&k.Suite, cattleConfig)
	providerRelations = append(providerRelations, amiRelationship)

	providerKeyPath := []string{permutationdata.ClusterConfigKey, permutationdata.ProviderKey}
	providerKeyValue, err := permutation.GetKeyPathValue(providerKeyPath, cattleConfig)
	require.NoError(k.T(), err)

	providerPermutation := permutation.CreatePermutation(providerKeyPath, providerKeyValue.([]any), providerRelations)

	permutedConfigs, _, err := permutation.Permute([]permutation.Permutation{k8sPermutation, providerPermutation}, cattleConfig)
	if err != nil {
		fmt.Println(err)
	}

	k.permutedConfigs = append(k.permutedConfigs, permutedConfigs...)
	for _, permutedConfig := range permutedConfigs {
		logrus.Info("---------------------------------------------")
		indented, _ := json.MarshalIndent(permutedConfig, "", "    ")
		converted := string(indented)
		fmt.Println(converted)
	}

	logrus.Info("------STATS------")
	logrus.Infof("Configs: %v", len(permutedConfigs))
	logrus.Info("---------------------------------------------")

}

func (k *V2ProvisioningNodeDriver) TestV2ProvisioningNodeDriver() {
	k.Run("test", func() {
		for _, permutedConfig := range k.permutedConfigs {
			var clusterConfig clusters.ClusterConfig
			config.LoadObjectFromMap(provisioningInputKey, permutedConfig, &clusterConfig)

			provider := clusterConfig.Provider
			nodeProvider := provisioning.CreateProvider(provider)

			clusterName := namegen.AppendRandomString(nodeProvider.Name.String())
			generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
			machinePoolConfigs := nodeProvider.MachinePoolFunc(permutedConfig, generatedPoolName, namespace)

			clusterObject, err := provisioning.CreateProvisioningCluster(k.client, nodeProvider, &clusterConfig, machinePoolConfigs, clusterName, nil)
			require.NoError(k.T(), err)

			k.clusters = append(k.clusters, *clusterObject)
		}

		for _, cluster := range k.clusters {
			provisioning.VerifyCluster(k.T(), k.client, nil, &cluster)
		}
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestV2ProvisioningNodeDriverTestSuite(t *testing.T) {
	suite.Run(t, new(V2ProvisioningNodeDriver))
}
