package util

import (
    "gopkg.in/mgo.v2"
    "fmt"
)

func NewMongoSession() *mgo.Session {

    host := fmt.Sprintf("mongodb://%s", config.MongoServer)

    // Connect to our local mongo
    var err error
    session, err := mgo.Dial(host)

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }

    c := session.DB("dolphin").C("notification")
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
