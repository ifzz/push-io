package util

import (
    "sync"
    "gopkg.in/mgo.v2"
)

var config *Config
var once sync.Once
var session *mgo.Session
var key *Key

func init() {
    once.Do(func() {
        config = InitConfig()
        session = NewMongoSession()
        key = InitKey()
    })
}
