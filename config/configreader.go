package config

import (
	"io/ioutil"
	"log"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Map []Mapper `yaml:"mapper"`
}

type Mapper struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

var (
	logger *zap.Logger
	err    error
)

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}
}

// YamlHandler reads yaml file and retuns struct
func YamlHandler() (*Config, error) {
	buf, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Fatal("error while reading file")
	}

	config := &Config{}

	logger.Info("getting ready for unmarshalling")

	err = yaml.Unmarshal(buf, config)

	if err != nil {
		log.Fatal("error occured while trying to unmarshal yaml file: ", err)
	}

	return config, nil
}

// YamlMapBuilder builds a map with path and URL mapping from yaml file
func YamlMapBuilder() map[string]string {
	config, err := YamlHandler()
	configMap := make(map[string]string)

	if err != nil {
		logger.Warn("couldnt read data from yaml file")
		return configMap
	}

	mapper := config.Map

	for _, mapping := range mapper {
		configMap[mapping.Path] = mapping.URL
	}
	return configMap

}
