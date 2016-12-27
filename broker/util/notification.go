package util

import (
    "github.com/parnurzeal/gorequest"
    "fmt"
    "time"
)

type Action interface {
    Save() error
    Notify() error
}

type Notification struct {
    Id        string `bson:"id"`
    Timestamp int64 `bson:"timestamp"`
    Created   time.Time `bson:"Created"`
    Qos       int `bson:"qos"`
    Retain    int `bson:"retain"`
    Topic     string `bson:"topic"`
    Message   string `bson:"message"`
}

func (n *Notification) Save() error {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("notification")
    if err := c.Insert(n); err != nil {
        return err
    }
    return nil
}

func (n *Notification) Notify() error {
    request := gorequest.New().SetBasicAuth("gftrader", "A98D8B1134D34F6E161463F757139")
    request.SetDebug(config.Debug)

    data := fmt.Sprintf(`{"qos":%d, "retain":%d, "topic":"%s", "message":"%s#%d#%s"}`,
        n.Qos, n.Retain, n.Topic, n.Id, n.Timestamp, n.Message)

    _, body, errs := request.Post(config.PushServer).
        Type("form").
        Send(data).
        End()
    fmt.Println(body)

    if len(errs) > 0 {
        return errs[0]
    }
    return nil
}
