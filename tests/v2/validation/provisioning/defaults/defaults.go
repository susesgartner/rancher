package permutations

import (
	"errors"
	"os"
	"path"
	"runtime"

	"sigs.k8s.io/yaml"
)

const (
	aws           = "aws.yaml"
	clusterConfig = "cluster_config.yaml"
)

func getDefaultsPath(filename string) string {
	_, directory, _, _ := runtime.Caller(0)
	filePath := path.Join(path.Dir(directory), filename)

	return filePath
}

func LoadDefaults(defaultFile string, defaultNames []string) (map[string]any, error) {
	if defaultFile == "" {
		yaml.Unmarshal([]byte("{}"), defaultFile)
		err := errors.New("No default file found")
		return nil, err
	}

	allString, err := os.ReadFile(defaultFile)
	if err != nil {
		panic(err)
	}

	var all map[string]any
	err = yaml.Unmarshal(allString, &all)
	if err != nil {
		panic(err)
	}

	var loadedDefaults map[string]any
	for _, defaultName := range defaultNames {
		defaultConfig := all[defaultName]
		loadedDefaults[defaultName] = defaultConfig
	}

	return loadedDefaults, nil
}
