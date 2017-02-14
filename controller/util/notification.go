package util

/*
import (
    "time"
    "gopkg.in/mgo.v2/bson"
)

type Notification struct {
    AppId       string `bson:"appId"`
    AppKey      string `bson:"appKey"`
    Id          string `bson:"id"`
    Timestamp   int64 `bson:"timestamp"`
    LastUpdated time.Time `bson:"lastUpdated"`
    Qos         int `bson:"qos"`
    Retain      int `bson:"retain"`
    Topic       string `bson:"topic"`
    Message     map[string]interface{} `bson:"message"`
    Error       string `bson:"error"`
    Ack         bool `bson:"ack"`
}

func (n *Notification) PushAck() error {
    s := session.Copy()
    defer s.Close()
    c := s.DB("dolphin").C("notification")

    now := time.Now()

    result := Notification{}
    err := c.Find(bson.M{"id": n.Id}).One(&result)
    if err != nil {
        return err
    }

    query := bson.M{"id": n.Id}
    change := bson.M{"$set": bson.M{"ack": true, "timestamp": now.Unix(), "LastUpdated": now}}

    return c.Update(query, change)
}
*/
