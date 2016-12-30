package util

import (
    "strings"
    "fmt"
    "strconv"
    "github.com/parnurzeal/gorequest"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func (n *Notification) GetStatus() {
    request := gorequest.New().SetBasicAuth(config.Username, config.Password)
    request.SetDebug(config.Debug)

    for _, value := range config.Cluster {
        fmt.Println(value)
        _, body, errs := request.Get(fmt.Sprintf("http://%s/api/stats", value)).End()
        if (len(errs) > 0) {
            fmt.Printf("error %+v\n", errs[0])
            continue
        }
        fmt.Printf("%s:%s\n", value, body)

        result := strings.Split(body, ",")
        clients_count := 0
        for _, stat := range result {
            pair := strings.Split(stat, ":")
            if strings.Contains(pair[0], "clients") && strings.Contains(pair[0], "count") {
                //fmt.Println(stat)
                clients_count, _ = strconv.Atoi(pair[1])
            }
        }

        cluster[value] = clients_count
    }
    fmt.Printf("%+v\n", cluster)
}

func GetServer() (string, int) {
    count := MaxInt
    host := ""

    for address, value := range cluster {
        if value < count {
            host = (strings.Split(address, ":"))[0]
            count = value
        }
    }
    fmt.Printf("%+v\n", cluster)
    return host, count
}
