package main

import (
    "fmt"
    "github.com/kataras/iris"
    "github.com/parnurzeal/gorequest"
    "./util"
)

func main() {
    config := util.GetInstance()
    //fmt.Printf("%+v\n", *config)

    iris.Get("/api/v1/server", func(ctx *iris.Context) {

        request := gorequest.New().SetBasicAuth(config.Username, config.Password)
        //request.SetDebug(true)

        for _, value := range config.Cluster {
            //fmt.Println(value)
            resp, body, errs := request.Get(fmt.Sprintf("http://%s/api/stats", value)).End()
            if (len(errs) > 0) {
                continue
            }
            fmt.Printf("resp: %+v\n", resp)
            fmt.Printf("body: %+v\n", body);
        }

        ctx.JSON(iris.StatusOK, iris.Map{
            "host": "",
        })
    })

    iris.Listen(":8080")
}
