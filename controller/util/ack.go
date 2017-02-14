package util

import (
    "time"
)

type Ack struct {
    Id          string `bson:"id"`
    Timestamp   int64 `bson:"timestamp"`
    LastUpdated time.Time `bson:"lastUpdated"`
}

func (ack *Ack) Acknowledge() error {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("ack")
    return c.Insert(ack)
}
