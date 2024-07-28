//go:build (validation || extended) && !infra.any && !infra.aks && !infra.eks && !infra.gke && !infra.rke2k3s && !cluster.any && !cluster.custom && !cluster.nodedriver && !sanity && !stress

package rke2

import (
	"fmt"
	"os"
	"testing"

	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/clusters/kubernetesversions"
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/rancher/shepherd/extensions/provisioning"
	"github.com/rancher/shepherd/pkg/config"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations/permutationdata"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	provisioningInputKey = "clusterConfig"
	k8sVersionKey        = "kubernetesVersion"
	awsMachineKey = "awsMachineConfigs"
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

	client, err := rancher.NewClient("", testSession)
	require.NoError(k.T(), err)
	k.client = client

	config, _ := config.LoadConfigFromFile(os.Getenv(constants.ConfigEnvironmentKey))
	logrus.Info(config)
/*
	k8sVersions, err := kubernetesversions.ListRKE2AllVersions(k.client)
	require.NoError(k.T(), err)*/

	k8sPermutation := permutation.Permutation{
		KeyPath:                   []string{provisioningInputKey, k8sVersionKey},
		KeyPathValues:             []any{config[provisioningInputKey].(map[string]any)[k8sVersionKey]...},
		KeyPathValueRelationships: []permutation.Relationship{},
	}

	amiPermutation := permutation.Permutation{
		KeyPath:                   []string{"awsMachineConfigs", "awsMachineConfig", "ami"},
		KeyPathValues:             []any{config["awsMachineConfigs"].(map[string]any)["awsMachineConfig"].([]map[string]any)[0]["ami"]...},
		KeyPathValueRelationships: []permutation.Relationship{},
	}

	amiRelationship := permutation.Relationship {
		ParentValue:       "aws",
		ChildKeyPath:      []string{},
		ChildKeyPathValue: "",
		ChildPermutations: []permutation.Permutation{amiPermutation},
	}

	providerPermutation := permutation.Permutation{
		KeyPath:                   []string{permutationdata.ClusterConfigKey, permutationdata.ProviderKey},
		KeyPathValues:             []any{config[permutationdata.ClusterConfigKey].(map[string]any)[permutationdata.k8sVersionKey]...},
		KeyPathValueRelationships: []permutation.Relationship{permutationsdata.LoadProviderRelationships(config)..., amiRelationship},
	}

	permutedConfigs, permutationNames, err := permutation.Permute([]permutation.Permutation{k8sPermutation, providerPermutation}, config)
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
	})*/
}

func (k *Sandbox) TestSandbox() {
	logrus.Info("test")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSandboxTestSuite(t *testing.T) {
	suite.Run(t, new(Sandbox))
}
