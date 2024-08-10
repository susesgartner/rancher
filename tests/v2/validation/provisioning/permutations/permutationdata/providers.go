package permutationdata

import (
	"fmt"

	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/stretchr/testify/suite"
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

func LoadProviderRelationships(s *suite.Suite, testConfig map[string]any) []permutation.Relationship {
	providers, _ := permutation.GetKeyPathValue([]string{ClusterConfigKey, ProviderKey}, testConfig)

	var providerRelationships []permutation.Relationship
	for _, provider := range providers.([]any) {
		switch {

		case provider == AWSName:
			providerRelationships = append(providerRelationships, LoadAWSRelationships(s, testConfig)...)

		case provider == AzureName:
			providerRelationships = append(providerRelationships, LoadAzureRelationships(s, testConfig)...)

			//TODO ADD SUPPORT FOR OTHER PROVIDERS
		/*
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
