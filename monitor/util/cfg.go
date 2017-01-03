package util

import (
    "github.com/spf13/viper"
)

const CONFIG_PATH = "."
const CONFIG_NAME = "config"

type Config struct {
    RedisHost string
    RedisPort int
    BrokerHost string
    BrokerPort int
    Username string
    Password string
    Proxy string
}

func NewConfig() *Config {
    viper.SetConfigName(CONFIG_NAME)
    viper.AddConfigPath(CONFIG_PATH)

    err := viper.ReadInConfig()
    if err != nil {
        panic("Config not found")
    }

    return &Config{
        RedisHost: viper.GetString("redisHost"),
        RedisPort: viper.GetInt("redisPort"),
        BrokerHost: viper.GetString("brokerHost"),
        BrokerPort: viper.GetInt("brokerPort"),
        Username: viper.GetString("username"),
        Password: viper.GetString("password"),
        Proxy: viper.GetString("proxy"),
    }
}
