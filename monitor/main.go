package main

import (
    "gopkg.in/alexcesaro/statsd.v2"
    "fmt"
    "time"
    "./util"
)

var config *util.Config

func init() {
    config = util.NewConfig()
    fmt.Printf("%+v\n", config)
}

func main() {
    // update the status of cluster per minute
    for t := range time.NewTicker(time.Minute).C {
        fmt.Println(t)
        util.Update()
        increment("dolphin.monitor.update")
    }
}

func increment(text string) {
    c, err := statsd.New(statsd.Address(config.StatsdServer))
    if err != nil {
        fmt.Printf("fail to initialize statsd %+v\n", err)
    } else {
        // Increment a counter.
        c.Increment(text)
    }
    defer c.Close()
}
