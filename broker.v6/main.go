package main

import (
    "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/adaptors/httprouter"
    "gopkg.in/kataras/iris.v6/adaptors/view"
    "./util"
)

var key = util.InitKey()
var config = util.InitConfig()
var jobQueue chan util.Job

func init() {
}

func main() {
    app := iris.New(iris.Configuration{Gzip: true, Charset: "UTF-8"})

    app.Adapt(iris.DevLogger())
    app.Adapt(httprouter.New())

    app.Adapt(view.HTML("./templates", ".html"))
    app.StaticWeb("/scripts", "./templates/scripts/")
    app.StaticWeb("/styles",  "./templates/styles")

    // HTTP Method: GET
    // PATH: http://127.0.0.1/
    // Handler(s): index
    app.Get("/", index)

    app.Post("/api/v1/login", login)

    //iris.Post("/api/v1/notification", notification)

    //iris.Get("/api/v1/message/:page/:pageSize", message)

    //iris.Get("/api/v1/application", application)

    app.Listen(":8080")
}

func index(ctx *iris.Context) {
    ctx.Render("index.html", nil)
}

func login(ctx *iris.Context) {
    type Account struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    data := &Account{}
    if err := ctx.ReadJSON(data); err != nil {
        //ctx.Log("%+v\n", err)
        ctx.EmitError(iris.StatusInternalServerError)
        return
    }
    if (!isAuthorized(data.Username, data.Password)) {
        ctx.EmitError(iris.StatusUnauthorized)
        return
    }
    ctx.JSON(iris.StatusOK, iris.Map{
        "status": "success",
    })
}

func isAuthorized(appId string, appKey string) bool {
    for _, configApp := range key.Apps {
        if configApp.AppId == appId && configApp.AppKey == appKey {
            return true
        }
    }
    return false
}
