package main

import (
    "fmt"
    "github.com/kataras/iris"
    "github.com/parnurzeal/gorequest"
    "encoding/json"
    "./util"
)

type Resp struct {
    clients_count int `json:"clients/count"`
}

func main() {
    config := util.GetInstance()
    //fmt.Printf("%+v\n", *config)

    iris.Get("/api/v1/server", func(ctx *iris.Context) {

        request := gorequest.New().SetBasicAuth(config.Username, config.Password)
        //request.SetDebug(true)

        count := 0
        host := ""
        for _, value := range config.Cluster {
            //fmt.Println(value)
            _, body, errs := request.Get(fmt.Sprintf("http://%s/api/stats", value)).End()
            if (len(errs) > 0) {
                fmt.Println("error %+v\n", errs[0])
                continue
            }
            fmt.Printf("body: %+v\n", body);
            r := &Resp{}
            if err := json.Unmarshal([]byte(body), r); err != nil {
                fmt.Printf("error %+v\n", err)
                continue
            }
            if (r.clients_count < count) {
                host = value
                count = r.clients_count
            }
        }

        ctx.JSON(iris.StatusOK, iris.Map{
            "host": host,
        })
    })

    iris.Listen(":8080")
}
