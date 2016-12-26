package main

import (
    "fmt"
    "github.com/kataras/iris"
    "github.com/parnurzeal/gorequest"
    "./util"
    "strings"
    "strconv"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

type Resp struct {
    clients_count uint `json:"clients\/count"`
}

func main() {
    config := util.GetInstance()
    //fmt.Printf("%+v\n", *config)

    iris.Get("/api/v1/server", func(ctx *iris.Context) {

        request := gorequest.New().SetBasicAuth(config.Username, config.Password)
        //request.SetDebug(true)

        count := MaxInt
        host := ""
        for _, value := range config.Cluster {
            fmt.Println(value)
            _, body, errs := request.Get(fmt.Sprintf("http://%s/api/stats", value)).End()
            if (len(errs) > 0) {
                fmt.Printf("error %+v\n", errs[0])
                continue
            }
            fmt.Println(body);

            result := strings.Split(body, ",")
            clients_count := 0
            for _, stat := range result {
                //fmt.Println(stat)
                pair := strings.Split(stat, ":")
                if strings.Contains(pair[0], "clients") && strings.Contains(pair[0], "count") {
                    fmt.Println(stat)
                    clients_count, _ = strconv.Atoi(pair[1])
                }
            }

            if (clients_count < count) {
                address := strings.Split(value, ":")
                host = address[0]
                count = clients_count
            }
        }

        ctx.JSON(iris.StatusOK, iris.Map{
            "host": host,
            "count": count,
        })
    })

    iris.Listen(":8080")
}
