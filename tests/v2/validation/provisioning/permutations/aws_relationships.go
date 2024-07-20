package permutations

import (
	"runtime"

	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/sirupsen/logrus"
)

const (
	provider              = "aws"
	provisioningConfigKey = "clusterConfig"
	machineConfigKey      = "awsMachineConfigs"
	credentialConfigKey   = "awsCredentials"
)

func MachineConfig(machineConfigKey string, machineConfig map[string]any) permutation.Relationship {
	return permutation.Relationship{
		ParentValue:       provider,
		ChildKeyPath:      []string{machineConfigKey},
		ChildKeyPathValue: machineConfig,
	}
}

func Credentials(credentialKeyPath []string, credentialsConfig map[string]any) permutation.Relationship {
	return permutation.Relationship{
		ParentValue:       provider,
		ChildKeyPath:      credentialKeyPath,
		ChildKeyPathValue: credentialsConfig,
	}
}

func NodeProvider(credentialKeyPath []string) permutation.Relationship {
	return permutation.Relationship{
		ParentValue:       provider,
		ChildKeyPath:      []string{provisioningConfigKey, "nodeProviders"},
		ChildKeyPathValue: "ec2",
	}
}

func LoadProviderDefaults() {
	
}

func CreateAWSRelationships() []permutation.Relationship {
	_, awsDefaultsFilePath, _, _ := runtime.Caller(1)
	logrus.Info(awsDefaultsFilePath)
	permutation.LoadDefaults()
}
