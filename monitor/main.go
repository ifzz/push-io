package main

import (
    "fmt"
    "time"
    "./util"
)

func main() {
    // update the status of cluster per minute
    for t := range time.NewTicker(time.Minute).C {
        fmt.Println(t)
        util.Update()
    }
}
