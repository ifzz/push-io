package util

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
)

type App struct {
    AppId     string `json:"appId"`
    AppKey    string `json:"appKey"`
    AppSecret string `json:"appSecret"`
}

type Key struct {
    Apps []App `json:"apps"`
}

const KEY_PATH = "./key.json"

func InitKey() *Key {
    file, err := ioutil.ReadFile(KEY_PATH)
    if err != nil {
        panic(fmt.Sprintf("fail to read file %s error %+v", KEY_PATH, err))
    }

    instance := &Key{}
    json.Unmarshal(file, instance)

    return instance
}
