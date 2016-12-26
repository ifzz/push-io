package util

import "sync"

var config *Config
var once sync.Once

func init() {
    GetInstance()
}
