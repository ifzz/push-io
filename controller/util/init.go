package util

import (
    "sync"
    "gopkg.in/mgo.v2"
)

var config *Config
var once sync.Once
var session *mgo.Session

func init() {
    once.Do(func() {
        config = NewConfig()
        session = NewMongoSession()
    })
}
