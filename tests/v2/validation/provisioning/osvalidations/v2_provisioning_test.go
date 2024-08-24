package osvalidations

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/rancher/tests/v2/actions/clusters"
	"github.com/rancher/rancher/tests/v2/actions/provisioning"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations/permutationdata"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/configoperations"
	"github.com/rancher/shepherd/extensions/configoperations/permutations"
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

type V2Provisioning struct {
	suite.Suite
	client             *rancher.Client
	session            *session.Session
	standardUserClient *rancher.Client
	permutedConfigs    []map[string]any
	clusters           []v1.SteveAPIObject
}

func (k *V2Provisioning) TearDownSuite() {
	k.session.Cleanup()
}

func (k *V2Provisioning) SetupSuite() {
	testSession := session.NewSession()
	k.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(k.T(), err)
	k.client = client

	cattleConfig := config.LoadConfigFromFile(os.Getenv(config.ConfigEnvironmentKey))

	k8sPermutation, err := permutationdata.CreateK8sPermutation(cattleConfig)
	require.NoError(k.T(), err)

	amiRelationship, err := permutationdata.CreateAMIRelationship(cattleConfig)

	providerPermutation, err := permutationdata.CreateProviderPermutation(cattleConfig)
	require.NoError(k.T(), err)
	providerPermutation.KeyPathValueRelationships = append(providerPermutation.KeyPathValueRelationships, amiRelationship)

	permutedConfigs, _, err := permutations.Permute([]permutations.Permutation{k8sPermutation, providerPermutation}, cattleConfig)
	require.NoError(k.T(), err)

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

func (k *V2Provisioning) TestV2ProvisioningNodeDriver() {
	k.Run("test", func() {
		for _, permutedConfig := range k.permutedConfigs {
			var clusterConfig clusters.ClusterConfig
			configoperations.LoadObjectFromMap(provisioningInputKey, permutedConfig, &clusterConfig)

			provider := clusterConfig.Provider
			nodeProvider := provisioning.CreateProvider(provider)

			clusterName := namegen.AppendRandomString(nodeProvider.Name.String())
			generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
			machinePoolConfigs := nodeProvider.MachinePoolFunc(permutedConfig, generatedPoolName, namespace)

			machineRoles := nodeProvider.MachineRolesFunc(permutedConfig)

			clusterObject, err := provisioning.CreateProvisioningCluster(k.client, nodeProvider, &clusterConfig, machinePoolConfigs, machineRoles, clusterName, nil)
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
func TestV2OSValidationTestSuite(t *testing.T) {
	suite.Run(t, new(V2Provisioning))
}
