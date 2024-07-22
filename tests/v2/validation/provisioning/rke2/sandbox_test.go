//go:build (validation || extended) && !infra.any && !infra.aks && !infra.eks && !infra.gke && !infra.rke2k3s && !cluster.any && !cluster.custom && !cluster.nodedriver && !sanity && !stress

package rke2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/rancher/tests/v2/validation/provisioning/constants"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/clusters/kubernetesversions"
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/rancher/shepherd/extensions/provisioning"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

	config, _ := permutation.LoadConfigFromFile(os.Getenv(constants.ConfigEnvironmentKey))
	logrus.Info(config)

	k8sVersions, err := kubernetesversions.ListRKE2AllVersions(k.client)
	require.NoError(k.T(), err)

	logrus.Info("--------------------------------------")
	logrus.Info("K8S VERSIONS:")
	logrus.Info(k8sVersions)
	logrus.Info("--------------------------------------")
	testPermutation1 := permutation.Permutation{
		KeyPath:                   []string{"provisioningInput", "kubernetesVersion"},
		KeyPathValues:             []any{k8sVersions[0], k8sVersions[1], k8sVersions[2]},
		KeyPathValueRelationships: []permutation.Relationship{},
	}

	permutedConfigs, permutationNames, err := permutation.Permute([]permutation.Permutation{testPermutation1}, config)
	if err != nil {
		fmt.Println(err)
	}

	/*
		for _, permutedConfig := range permutedConfigs {
			logrus.Info("---------------------------------------------")
			indented, _ := json.MarshalIndent(permutedConfig, "", "    ")
			converted := string(indented)
			fmt.Println(converted)
		}
	*/
	logrus.Info("------STATS------")
	logrus.Infof("Configs: %v", len(permutedConfigs))
	logrus.Info("---------------------------------------------")

	k.T().Run(name, func() {
		clusterObjects := []v1.SteveAPIObject
		for _, permutedConfig := range permutedConfigs {
			logrus.Info("---------------------------------------------")
			logrus.Info("CONVERTED TO CLUSTERCONFIG")
			k.provisioningConfig = new(clusters.ClusterConfig)
			permutation.LoadConfigFromMap("provisioningInput", permutedConfig, k.provisioningConfig)
			indented, _ := json.MarshalIndent(k.provisioningConfig, "", "    ")
			converted := string(indented)
			fmt.Println(converted)

			logrus.Info("Provisioning Clusters")
			providers := *k.provisioningConfig.Providers
			provider := provider[0]
			nodeProvider := provisioning.CreateProvider(provider)
			clusterObject, err := provisioning.CreateProvisioningCluster(k.client, nodeProvider, k.provisioningConfig, nil)
			require.NoError(s.T(), err)

			clusterObjects = append(clusterObjects, clusterObject)
		}

		for _, clusterObject := range clusterObjects {
			provisioning.VerifyCluster(s.T(), client, testClusterConfig, clusterObject)
		}
	})
}

func (k *Sandbox) TestSandbox() {
	logrus.Info("test")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSandboxTestSuite(t *testing.T) {
	suite.Run(t, new(Sandbox))
}
