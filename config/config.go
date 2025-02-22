package config

import (
    m "everything/models"

    "github.com/spf13/viper"
)

func LoadConfig() (config m.Config, err error) {
    v := viper.New()
    v.AddConfigPath("config/")
    v.SetConfigName("config")
    v.SetConfigType("toml")
    err = v.ReadInConfig()
    if err != nil {
        return
    }
    
    err = v.Unmarshal(&config)
    return
}
