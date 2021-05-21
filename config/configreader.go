package config

import (
	"io/ioutil"
	"log"
	"net/http"

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
	logger       *zap.Logger
	err          error
	PathToURLMap URLMapper
)

type URLMapper map[string]string

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}
	PathToURLMap = make(map[string]string)
	PathToURLMap["mail"] = "https://mail.google.com/mail/u/0/#inbox"
}

// YamlHandler parses yaml file and unmarshalls
// yaml bytes into struct defined
func YamlHandler(yamlFilePath string) (*Config, error) {
	logger.Info("yamlFilePath is ", zap.String("FilePath", yamlFilePath))
	buf, err := ioutil.ReadFile(yamlFilePath)

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

// YamlMapBuilder builds a map from data
// present in yaml file as Path and related
// URL to be redirected to. The map created
// is used to fetch URL based on path
func YamlMapBuilder(yamlFilePath string) map[string]string {
	config, err := YamlHandler(yamlFilePath)
	configMap := make(map[string]string)

	if err != nil {
		logger.Fatal("couldnt read data from yaml file")
		return configMap
	}

	mapper := config.Map

	for _, mapping := range mapper {
		configMap[mapping.Path] = mapping.URL
	}
	return configMap

}

// URLRedirectMW middleware function which
// acts on request and redirects URL
// to another URL when it sees path matching
// another URL defined in yaml file
func URLRedirectMW(next http.Handler, yamlFilePath string) http.HandlerFunc {

	yamlMap := YamlMapBuilder(yamlFilePath)
	// combining 2 different map configs from map built from yaml file and local map config
	PathToURLMap := FuseMap(PathToURLMap, yamlMap)
	logger.Info("Map fused")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received a request")
		app := r.URL.Path
		if redirectURL, ok := PathToURLMap[app]; ok {
			logger.Info("URL is redirecting to %s since received app is %s \n", zap.String("redirectURL", redirectURL), zap.String("app", app))
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		} else {
			logger.Info("no mapping defined for \n", zap.String("app", app))
		}

		next.ServeHTTP(w, r)
	})
}

// FuseMap combines 2 maps. m1 & m2 into
// single map key values present in m2 will
// over write key value pair from m1 map
func FuseMap(m1, m2 URLMapper) URLMapper {
	for key, value := range m1 {
		if _, ok := m2[key]; !ok {
			m2[key] = value
		}
	}
	return m2
}
