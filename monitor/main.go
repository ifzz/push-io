package main

import (
    "fmt"
    "time"
    "./util"
)

func main() {
    for t := range time.NewTicker(time.Minute).C {
        fmt.Println(t)
        util.Update()
    }
}
