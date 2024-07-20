package rke2

import (
	"github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/extensions/clusters/kubernetesversions"
	"github.com/rancher/shepherd/extensions/permutation"
	"github.com/sirupsen/logrus"
)

func LoadRKE2Defaults(client *rancher.Client, filePath string, overrideConfig map[string]any) (map[string]any, error) {
	config, err := permutation.LoadDefaults(filePath, overrideConfig)
	logrus.Info(config)
	if err != nil {
		return nil, err
	}

	if _, ok := config["kubernetesVersion"].(string); ok {
		k8sVersions, err := kubernetesversions.ListRKE2AllVersions(client)
		if err != nil {
			return nil, err
		}

		config["kubernetesVersion"] = k8sVersions[0]
	}

	return config, nil
}
