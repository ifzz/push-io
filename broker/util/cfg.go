package util

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
)

const CONFIG_PATH = "./config.json"

type Config struct {
    Username   string
    Password   string
    //Cluster    []string
    Debug      bool
    PushServer string
    MaxWorkers int
    MaxQueue   int
    MongoServer string
    MongoPort int
}

func InitConfig() *Config {
    file, err := ioutil.ReadFile(CONFIG_PATH)
    if err != nil {
        panic(fmt.Sprintf("fail to read file %s error %+v", CONFIG_PATH, err))
    }

    config := &Config{}
    json.Unmarshal(file, config)

    return config
}
