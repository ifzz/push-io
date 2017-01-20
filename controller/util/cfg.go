package util

import "github.com/spf13/viper"

const CONFIG_PATH = "."
const CONFIG_NAME = "config"

type Config struct {
    ClientId string
    Username string
    Password string
    MongoServer string
    //MongoPort int
    Debug bool
    MqttServer string
    MaxWorkers int
    MaxQueue   int
}

func NewConfig() *Config {
    viper.SetConfigName(CONFIG_NAME)
    viper.AddConfigPath(CONFIG_PATH)

    err := viper.ReadInConfig()
    if err != nil {
        panic("Config not found")
    }

    return &Config{
        ClientId: viper.GetString("clientId"),
        Username: viper.GetString("username"),
        Password: viper.GetString("password"),
        Debug: viper.GetBool("debug"),
        MqttServer: viper.GetString("mqttServer"),
        MongoServer: viper.GetString("mongoServer"),
        //MongoPort: viper.GetInt("mongoPort"),
        MaxWorkers: viper.GetInt("maxWorkers"),
        MaxQueue: viper.GetInt("maxQueue"),
    }
}
