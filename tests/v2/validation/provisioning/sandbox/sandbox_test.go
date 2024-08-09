package sandbox

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations/permutationdata"
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const (
	provisioningInputKey = "clusterConfig"
	k8sVersionKey        = "kubernetesVersion"
	awsMachineKey        = "awsMachineConfigs"
)

type Sandbox struct {
	suite.Suite
	client             *rancher.Client
	session            *session.Session
	standardUserClient *rancher.Client
	provisioningConfig *clusters.ClusterConfig
}

func (k *Sandbox) TearDownSuite() {
	k.session.Cleanup()
}

func (k *Sandbox) SetupSuite() {
	testSession := session.NewSession()
	k.session = testSession

	//client, err := rancher.NewClient("", testSession)
	//require.NoError(k.T(), err)
	//k.client = client

	config := config.LoadConfigFromFile(os.Getenv(config.ConfigEnvironmentKey))
	logrus.Info(config)
	/*
		k8sVersions, err := kubernetesversions.ListRKE2AllVersions(k.client)
		require.NoError(k.T(), err)*/
	k8sPermutation := permutation.Permutation{
		KeyPath:                   []string{provisioningInputKey, k8sVersionKey},
		KeyPathValues:             config[provisioningInputKey].(map[string]any)[k8sVersionKey].([]any),
		KeyPathValueRelationships: []permutation.Relationship{},
	}

	amiPermutation := permutation.Permutation{
		KeyPath:                   []string{"awsMachineConfigs", "awsMachineConfig", "ami"},
		KeyPathValues:             config["awsMachineConfigs"].(map[string]any)["awsMachineConfig"].([]any)[0].(map[string]any)["ami"].([]any),
		KeyPathValueRelationships: []permutation.Relationship{},
	}

	amiRelationship := permutation.Relationship{
		ParentValue:       "aws",
		ChildKeyPath:      []string{},
		ChildKeyPathValue: "",
		ChildPermutations: []permutation.Permutation{amiPermutation},
	}

	providerRelations := permutationdata.LoadProviderRelationships(config)
	providerRelations = append(providerRelations, amiRelationship)
	providerPermutation := permutation.Permutation{
		KeyPath:                   []string{permutationdata.ClusterConfigKey, permutationdata.ProviderKey},
		KeyPathValues:             config[permutationdata.ClusterConfigKey].(map[string]any)[permutationdata.ProviderKey].([]any),
		KeyPathValueRelationships: providerRelations,
	}

	permutedConfigs, _, err := permutation.Permute([]permutation.Permutation{k8sPermutation, providerPermutation}, config)
	if err != nil {
		fmt.Println(err)
	}

	for _, permutedConfig := range permutedConfigs {
		logrus.Info("---------------------------------------------")
		indented, _ := json.MarshalIndent(permutedConfig, "", "    ")
		converted := string(indented)
		fmt.Println(converted)
	}

	logrus.Info("------STATS------")
	logrus.Infof("Configs: %v", len(permutedConfigs))
	logrus.Info("---------------------------------------------")
	/*
		k.T().Run(name, func() {
			clusterObjects := []v1.SteveAPIObject
			for _, permutedConfig := range permutedConfigs {
				k.provisioningConfig = new(clusters.ClusterConfig)
				config.LoadObjectFromMap(provisioningInputKey, permutedConfig, k.provisioningConfig)

				logrus.Info("Provisioning Clusters")
				providers := *k.provisioningConfig.Providers
				provider := provider[0]
				nodeProvider := provisioning.CreateProvider(provider)

				clusterName := namegen.AppendRandomString(nodeProvider.Name.String())
				generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
				machinePoolConfigs := nodeProvider.MachinePoolFunc(permutedConfig, generatedPoolName, namespace)

				clusterObject, err := provisioning.CreateProvisioningCluster(k.client, nodeProvider, k.provisioningConfig, machinePoolConfigs, clusterName, nil)
				require.NoError(s.T(), err)

				clusterObjects = append(clusterObjects, clusterObject)
			}

			for _, clusterObject := range clusterObjects {
				provisioning.VerifyCluster(s.T(), client, testClusterConfig, clusterObject)
			}
		})
	*/
}

func (k *Sandbox) TestSandbox() {
	logrus.Info("test")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSandboxTestSuite(t *testing.T) {
	suite.Run(t, new(Sandbox))
}
