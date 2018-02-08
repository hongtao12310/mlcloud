package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
    "github.com/deepinsight/mlcloud/src/utils"
)

// SubmitConfigDataCenter is inner conf for mlcloud
type SubmitConfigDataCenter struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Usercert string `yaml:"usercert"`
	Userkey  string `yaml:"userkey"`
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
}

// SubmitConfig is configuration load from user config yaml files
type SubmitConfig struct {
	DC                []SubmitConfigDataCenter `yaml:"datacenters"`
	ActiveConfig      *SubmitConfigDataCenter
	CurrentDatacenter string `yaml:"current-datacenter"`
}

// DefaultConfigFile returns the path of paddlecloud default config file path
func DefaultConfigFile() string {
	return filepath.Join(utils.UserHomeDir(), ".mlcloud", "config")
}

// ParseDefaultConfig returns default parsed config struct in ~/.mlcloud/config
func ParseDefaultConfig() (*SubmitConfig, error) {
	return ParseConfig(DefaultConfigFile())
}

// ParseConfig parse paddlecloud config to a struct
func ParseConfig(configFile string) (*SubmitConfig, error) {
	// ------------------- load paddle config -------------------
	buf, err := ioutil.ReadFile(configFile)
	config := SubmitConfig{}
	if err == nil {
		yamlErr := yaml.Unmarshal(buf, &config)
		if yamlErr != nil {
			return nil, fmt.Errorf("load config %s error: %v", configFile, yamlErr)
		}
		// put active config
		config.ActiveConfig = nil
		for _, item := range config.DC {
			if item.Name == config.CurrentDatacenter {
				config.ActiveConfig = &item
				break
			}
		}

		if config.ActiveConfig == nil {
			return nil, errors.New("current data center is not defined!\n")
		}
		// TODO: check token cache
		// if len(token) == 0 {return err}
		return &config, nil
	}
	return nil, fmt.Errorf("config %s error: %v", configFile, err)
}
