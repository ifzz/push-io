package util

import (
    "github.com/parnurzeal/gorequest"
    "fmt"
    "time"
    "encoding/json"
    "net/http"
    "strings"
    "gopkg.in/mgo.v2/bson"
    "bytes"
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

    var result []Notification
    err := c.Find(nil).Skip((page - 1) * pageSize).Limit(pageSize).
        //Select(bson.M{"id": 1, "ack": 1, "message": 1, "appId": 1, "timestamp": 1, "error": 1, "topic": 1, "lastUpdated": 1}).
        All(&result)

    copy(rows, result)

    return err
}

func Application() ([]string, error) {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("notification")
    var result []string

    err := c.Find(nil).Select(bson.M{"appId": 1}).Distinct("appId", &result)

    rows := make([]string, len(result))

    copy(rows, result)

    return rows, err
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
        e := Exception{
            AppId: n.AppId,
            AppKey: n.AppKey,
            Id: n.Id,
            Timestamp: n.Timestamp,
            LastUpdated: n.LastUpdated,
            Description: getErrorMessage([]error{err}),
            Message: n.Message,
        }
        return e.Save()
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
        e := Exception{
            AppId: n.AppId,
            AppKey: n.AppKey,
            Id: n.Id,
            Timestamp: n.Timestamp,
            LastUpdated: n.LastUpdated,
            Description: getErrorMessage(errs),
            Message: n.Message,
        }
        return e.Save()
    }

    return n.Save()
}

func getErrorMessage(errs []error) string {
    var buffer bytes.Buffer
    for _, e := range errs {
        buffer.WriteString(e.Error())
    }
    fmt.Println(buffer.String())
    return buffer.String()
}
