package util

import (
    "gopkg.in/mgo.v2"
    "fmt"
    "time"
)

func NewMongoSession() *mgo.Session {

    host := fmt.Sprintf("mongodb://%s:%d", config.MongoServer, config.MongoPort)

    // Connect to our local mongo
    var err error
    session, err := mgo.Dial(host)

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }

    c := session.DB("dolphin").C("notification")

    /*index := mgo.Index{
        Key: []string{"lastUpdated"},
        ExpireAfter: time.Duration(config.TTL) * time.Hour,
    }
    if err := c.EnsureIndex(index); err != nil {
        panic(err)
    }*/

    // Index
    index := mgo.Index{
        Key:        []string{"id"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    err = c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }

    session.SetMode(mgo.Monotonic, true)

    return session
}
