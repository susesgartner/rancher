//go:build (validation || extended) && !infra.any && !infra.aks && !infra.eks && !infra.gke && !infra.rke2k3s && !cluster.any && !cluster.custom && !cluster.nodedriver && !sanity && !stress

package rke2

import (
	"os"
	"testing"

	"github.com/rancher/rancher/tests/v2/actions/clusters"
	"github.com/rancher/rancher/tests/v2/actions/config/permutationdata"
	"github.com/rancher/rancher/tests/v2/actions/machinepools"
	"github.com/rancher/rancher/tests/v2/actions/provisioning"
	"github.com/rancher/rancher/tests/v2/actions/provisioninginput"
	"github.com/rancher/rancher/tests/v2/actions/reports"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/shepherd/extensions/cloudcredentials"
	"github.com/rancher/shepherd/extensions/clusters/kubernetesversions"
	"github.com/rancher/shepherd/extensions/users"
	password "github.com/rancher/shepherd/extensions/users/passwordgenerator"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/config/operations"
	"github.com/rancher/shepherd/pkg/config/operations/permutations"
	"github.com/rancher/shepherd/pkg/environmentflag"
	namegen "github.com/rancher/shepherd/pkg/namegenerator"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RKE2NodeDriverProvisioningTestSuite struct {
	suite.Suite
	client             *rancher.Client
	session            *session.Session
	standardUserClient *rancher.Client
	provisioningConfig *provisioninginput.Config
	permutedConfigs    []map[string]any
}

func (r *RKE2NodeDriverProvisioningTestSuite) TearDownSuite() {
	r.session.Cleanup()
}

func (r *RKE2NodeDriverProvisioningTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)
	r.client = client

	cattleConfig := config.LoadConfigFromFile(os.Getenv(config.ConfigEnvironmentKey))

	providerPermutation, err := permutationdata.CreateProviderPermutation(cattleConfig)
	require.NoError(r.T(), err)

	k8sPermutation, err := permutationdata.CreateK8sPermutation(cattleConfig)
	require.NoError(r.T(), err)

	if len(k8sPermutation.KeyPathValues) == 0 || k8sPermutation.KeyPathValues == nil {
		rke2Versions, err := kubernetesversions.ListRKE2AllVersions(client)
		require.NoError(r.T(), err)

		logrus.Infof("Setting k8s versions to %v", rke2Versions)

		for _, v := range rke2Versions {
			k8sPermutation.KeyPathValues = append(k8sPermutation.KeyPathValues, v)
		}
	}

	cniPermutation, err := permutationdata.CreateCNIPermutation(cattleConfig)
	require.NoError(r.T(), err)

	permutedConfigs, err := permutations.Permute([]permutations.Permutation{*k8sPermutation, *providerPermutation, *cniPermutation}, cattleConfig)
	require.NoError(r.T(), err)

	r.permutedConfigs = append(r.permutedConfigs, permutedConfigs...)

	enabled := true
	var testuser = namegen.AppendRandomString("testuser-")
	var testpassword = password.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(r.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(r.T(), err)

	r.standardUserClient = standardUserClient
}

func (r *RKE2NodeDriverProvisioningTestSuite) TestProvisioningRKE2Cluster() {
	nodeRolesAll := []provisioninginput.MachinePools{provisioninginput.AllRolesMachinePool}
	nodeRolesShared := []provisioninginput.MachinePools{provisioninginput.EtcdControlPlaneMachinePool, provisioninginput.WorkerMachinePool}
	nodeRolesDedicated := []provisioninginput.MachinePools{provisioninginput.EtcdMachinePool, provisioninginput.ControlPlaneMachinePool, provisioninginput.WorkerMachinePool}
	nodeRolesWindows := []provisioninginput.MachinePools{provisioninginput.EtcdMachinePool, provisioninginput.ControlPlaneMachinePool, provisioninginput.WorkerMachinePool, provisioninginput.WindowsMachinePool}

	tests := []struct {
		name         string
		machinePools []provisioninginput.MachinePools
		client       *rancher.Client
		isWindows    bool
		runFlag      bool
	}{
		{"1 Node all roles " + provisioninginput.StandardClientName.String(), nodeRolesAll, r.standardUserClient, false, r.client.Flags.GetValue(environmentflag.Short) || r.client.Flags.GetValue(environmentflag.Long)},
		{"2 nodes - etcd|cp roles per 1 node " + provisioninginput.StandardClientName.String(), nodeRolesShared, r.standardUserClient, false, r.client.Flags.GetValue(environmentflag.Short) || r.client.Flags.GetValue(environmentflag.Long)},
		{"3 nodes - 1 role per node " + provisioninginput.StandardClientName.String(), nodeRolesDedicated, r.standardUserClient, false, r.client.Flags.GetValue(environmentflag.Long)},
		{"4 nodes - 1 role per node + 1 Windows worker " + provisioninginput.StandardClientName.String(), nodeRolesWindows, r.standardUserClient, true, r.client.Flags.GetValue(environmentflag.Long)},
	}

	for _, tt := range tests {
		if !tt.runFlag {
			r.T().Logf("SKIPPED")
			continue
		}

		r.Run(tt.name, func() {
			for _, permutedConfig := range r.permutedConfigs {
				clusterConfig := new(clusters.ClusterConfig)
				operations.LoadObjectFromMap(permutationdata.ClusterConfigKey, permutedConfig, clusterConfig)

				require.NotNil(r.T(), clusterConfig.Provider)
				if clusterConfig.Provider != "vsphere" && tt.isWindows {
					r.T().Skip("Windows test requires access to vsphere")
				}

				clusterConfig.MachinePools = tt.machinePools

				provider := provisioning.CreateProvider(clusterConfig.Provider)
				credentialSpec := cloudcredentials.LoadCloudCredential(string(provider.Name))
				machineConfigSpec := machinepools.LoadMachineConfigs(string(provider.Name))

				clusterObject, err := provisioning.CreateProvisioningCluster(tt.client, provider, credentialSpec, clusterConfig, machineConfigSpec, nil)
				reports.TimeoutClusterReport(clusterObject, err)
				require.NoError(r.T(), err)

				provisioning.VerifyCluster(r.T(), tt.client, clusterConfig, clusterObject)
			}
		})
	}
}

func (r *RKE2NodeDriverProvisioningTestSuite) TestProvisioningRKE2ClusterDynamicInput() {
	tests := []struct {
		name   string
		client *rancher.Client
	}{
		{provisioninginput.AdminClientName.String(), r.client},
		{provisioninginput.StandardClientName.String(), r.standardUserClient},
	}

	for _, tt := range tests {
		r.Run(tt.name, func() {
			for _, permutedConfig := range r.permutedConfigs {
				clusterConfig := new(clusters.ClusterConfig)
				operations.LoadObjectFromMap(permutationdata.ClusterConfigKey, permutedConfig, clusterConfig)
				if len(clusterConfig.MachinePools) == 0 {
					r.T().Skip()
				}

				provider := provisioning.CreateProvider(clusterConfig.Provider)
				credentialSpec := cloudcredentials.LoadCloudCredential(string(provider.Name))
				machineConfigSpec := machinepools.LoadMachineConfigs(string(provider.Name))

				clusterObject, err := provisioning.CreateProvisioningCluster(tt.client, provider, credentialSpec, clusterConfig, machineConfigSpec, nil)
				reports.TimeoutClusterReport(clusterObject, err)
				require.NoError(r.T(), err)

				provisioning.VerifyCluster(r.T(), tt.client, clusterConfig, clusterObject)

			}
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRKE2ProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(RKE2NodeDriverProvisioningTestSuite))
}
