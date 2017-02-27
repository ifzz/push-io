package util

import "time"

type Exception struct {
    AppId       string `json:"appId" bson:"appId"`
    AppKey      string `json:"-" bson:"appKey"`
    Id          string `json:"id" bson:"id"`
    Timestamp   int64 `json:"timestamp" bson:"timestamp"`
    LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
    Description string `json:"description" bson:"description"`
    Message     map[string]interface{} `json:"message" bson:"message"`
}

func (e *Exception) Save() error {
    s := session.Copy()
    defer s.Close()

    c := s.DB("dolphin").C("exception")
    return c.Insert(e)
}
