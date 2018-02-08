package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"github.com/deepinsight/mlcloud/src/utils/log"
	"github.com/deepinsight/mlcloud/src/pkg/common/syscfg/store"
)

var (
	// the default path of configuration file
	defaultPath = os.Getenv("HOME") + "/mlcloud/config/config.json"
)

type cfgStore struct {
	path string // the path of cfg file
	sync.RWMutex
}

// NewCfgStore returns an instance of cfgStore that stores the configurations
// in a json file. The file will be created if it does not exist.
func NewCfgStore(path ...string) (store.Driver, error) {
	p := defaultPath
	if len(path) > 0 && len(path[0]) > 0 {
		p = path[0]
	}

	log.Debugf("path of configuration file: %s", p)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		log.Infof("the configuration file %s does not exist, creating it...", p)
		if err = os.MkdirAll(filepath.Dir(p), 0600); err != nil {
			return nil, err
		}

        if _, err := os.Create(p); err != nil {
            return nil, err
        }

		if err = ioutil.WriteFile(p, []byte{}, 0600); err != nil {
			return nil, err
		}
	}

	return &cfgStore{
		path: p,
	}, nil
}

// Name ...
func (c *cfgStore) Name() string {
	return "JSON"
}

// Read ...
func (c *cfgStore) Read() (map[string]interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	return read(c.path)
}

func read(path string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// empty file
	if len(b) == 0 {
		return nil, nil
	}

	config := map[string]interface{}{}
	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// Write ...
func (c *cfgStore) Write(config map[string]interface{}) error {
	c.Lock()
	defer c.Unlock()

	cfg, err := read(c.path)
	if err != nil {
		return err
	}

	if cfg == nil {
		cfg = config
	} else {
		for k, v := range config {
			cfg[k] = v
		}
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(c.path, b, 0600); err != nil {
		return err
	}

	return nil
}
