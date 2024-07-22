package providers

import (
	"fmt"

	"github.com/rancher/rancher/tests/v2/validation/provisioning/provisioningconfig/clusterconfig"
	"github.com/rancher/shepherd/extensions/permutation"
)

const (
	AWSName       = "aws"
	AzureName     = "azure"
	DOName        = "do"
	HarvesterName = "harvester"
	LinodeName    = "linode"
	GoogleName    = "google"
	VsphereName   = "vsphere"
)

func LoadProviderRelationships(testConfig map[string]any) []permutation.Relationship {
	providers := testConfig[clusterconfig.ClusterConfigKey].(map[string]any)[clusterconfig.ProviderKey].(*[]string)

	var providerRelationships []permutation.Relationship
	for _, provider := range *providers {
		switch {

		case provider == AWSName:
			providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			//TODO ADD SUPPORT FOR OTHER PROVIDERS
		/*
			case provider == providerNames.AzureName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			case provider == providerNames.VsphereName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			case provider == providerNames.LinodeName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			case provider == providerNames.HarvesterName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			case provider == providerNames.DOName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...)

			case provider == providerNames.GoogleName:
				providerRelationships = append(providerRelationships, LoadAWSRelationships(testConfig)...) */

		default:
			panic(fmt.Sprintf("Provider:%v not found", provider))
		}

	}

	return providerRelationships
}
