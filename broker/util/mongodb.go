package util

import (
    "gopkg.in/mgo.v2"
    "fmt"
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

    session.SetMode(mgo.Monotonic, true)

    return session
}
