package util

import (
    "github.com/parnurzeal/gorequest"
    "fmt"
    "time"
    "encoding/json"
)

type Action interface {
    Save() error
    Notify() error
}

type Notification struct {
    AppId     string `bson:"appId"`
    AppKey    string `bson:"appKey"`
    Id        string `bson:"id"`
    Timestamp int64 `bson:"timestamp"`
    Created   time.Time `bson:"Created"`
    Qos       int `bson:"qos"`
    Retain    int `bson:"retain"`
    Topic     string `bson:"topic"`
    Message   map[string]interface{} `bson:"message"`
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
    request := gorequest.New().SetBasicAuth(n.AppId, n.AppKey)
    request.SetDebug(config.Debug)

    /*data := fmt.Sprintf(`{"qos":%d, "retain":%d, "topic":"%s", "message":"%s#%d#%s"}`,
        n.Qos, n.Retain, n.Topic, n.Id, n.Timestamp, n.Message)*/

    type Content struct {
        Id string `json:"id"`
        Timestamp int64 `json:"timestamp"`
        Payload map[string]interface{} `json:"payload"`
    }

    type Data struct {
        Qos int `json:"qos"`
        Retain int `json:"retain"`
        Topic string `json:"topic"`
        Message string `json:"message"`
    }

    content := Content{
        Id: n.Id,
        Timestamp: n.Timestamp,
        Payload: n.Message,
    }

    jsonString, err := json.Marshal(content)
    if err != nil {
        return err
    }
    fmt.Println(string(jsonString))

    data := Data{
        Qos: n.Qos,
        Retain: n.Retain,
        Topic: n.Topic,
        Message: string(jsonString),
    }

    //fmt.Println(strconv.Quote(string(jsonString)))
    _, _, errs := request.Post(config.PushServer).
        Type("form").
        Send(data).
        End()
    //fmt.Println(body)

    if len(errs) > 0 {
        return errs[0]
    }
    return nil
}
