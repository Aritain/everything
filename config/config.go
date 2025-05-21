package config

import (
	"sync"

	"everything/models"

	"github.com/spf13/viper"
)

var (
	serviceInstance *ConfigService
	once            sync.Once
)

type ConfigService struct {
	mu     sync.RWMutex
	config models.Config
}

func Initialize() error {
	var initErr error
	once.Do(func() {
		v := viper.New()
		v.AddConfigPath("config/")
		v.SetConfigName("config")
		v.SetConfigType("toml")

		if err := v.ReadInConfig(); err != nil {
			initErr = err
			return
		}

		var cfg models.Config
		if err := v.Unmarshal(&cfg); err != nil {
			initErr = err
			return
		}

		serviceInstance = &ConfigService{
			config: cfg,
		}
	})
	return initErr
}
func Get() *ConfigService {
	return serviceInstance
}

func (cs *ConfigService) Config() models.Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.config
}
