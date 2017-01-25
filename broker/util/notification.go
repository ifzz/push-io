package util

import (
    "github.com/parnurzeal/gorequest"
    "fmt"
    "time"
    "encoding/json"
    "net/http"
    "strings"
)

type Notification struct {
    AppId       string `json:"appId" bson:"appId"`
    AppKey      string `json:"-" bson:"appKey"`
    Id          string `json:"id" bson:"id"`
    Timestamp   int64 `json:"timestamp" bson:"timestamp"`
    LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
    Qos         int `json:"-" bson:"qos"`
    Retain      int `json:"-" bson:"retain"`
    Topic       string `json:"topic" bson:"topic"`
    Message     map[string]interface{} `json:"message" bson:"message"`
    Error       string `json:"error" bson:"error"`
    Ack         bool `json:"ack" bson:"ack"`
}

func (n *Notification) Save() error {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("notification")
    return c.Insert(n)
}

func Total() (int, error) {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("notification")

    return c.Find(nil).Count()
}

func List(rows []Notification, page int, pageSize int) error {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("notification")

    var results []Notification
    err := c.Find(nil).Skip((page - 1) * pageSize).Limit(pageSize).
        //Select(bson.M{"id": 1, "ack": 1, "message": 1, "appId": 1, "timestamp": 1, "error": 1, "topic": 1, "lastUpdated": 1}).
        All(&results)

    copy(rows, results)

    return err
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
        Topic: strings.ToUpper(n.AppId) + "/" + n.Topic,
        Message: string(jsonString),
    }

    _, _, errs := request.Post(config.PushServer).
        Type("form").
        Send(data).
        Retry(3, 2 * time.Minute, http.StatusBadRequest, http.StatusInternalServerError).
        End()

    if len(errs) > 0 {
        err := fmt.Sprintf("%+v", errs[0])
        fmt.Println(err)
        n.Error = err
    }
    n.Ack = false

    return n.Save()
}
