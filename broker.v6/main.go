package main

import (
    "gopkg.in/alexcesaro/statsd.v2"
    "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/adaptors/httprouter"
    "gopkg.in/kataras/iris.v6/adaptors/view"
    "./util"
    "time"
    "github.com/satori/go.uuid"
    "fmt"
)

var key = util.InitKey()
var config = util.InitConfig()
var jobQueue chan util.Job

func init() {
    // initialize the job queue
    jobQueue = make(chan util.Job, config.MaxQueue)
    dispatcher := util.NewDispatcher(jobQueue)
    dispatcher.Run()
}

func main() {
    app := iris.New(iris.Configuration{Gzip: true, Charset: "UTF-8"})

    app.Adapt(iris.DevLogger())
    app.Adapt(httprouter.New())
    app.Adapt(view.HTML("./templates", ".html"))

    app.StaticWeb("/scripts", "./templates/scripts/")
    app.StaticWeb("/styles",  "./templates/styles")

    app.Get("/", index)
    app.Post("/api/v1/login", login)
    app.Post("/api/v1/notification", notification)
    app.Get("/api/v1/message/:page/:pageSize", message)
    app.Get("/api/v1/application", application)

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
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }
    if (!isAuthorized(data.Username, data.Password)) {
        ctx.JSON(iris.StatusUnauthorized, iris.Map{
            "error": "wrong username or password",
        })
        return
    }
    increment("dolphin.api.v1.login")
    ctx.JSON(iris.StatusOK, iris.Map{
        "status": "success",
    })
}

func notification(ctx *iris.Context) {
    type Data struct {
        Topic   string `json:"topic"`
        Message map[string]interface{} `json:"message"`
        AppId string `json:"appId"`
        AppKey string `json:"appKey"`
    }
    data := &Data{}
    if err := ctx.ReadJSON(data); err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }
    if (!isAuthorized(data.AppId, data.AppKey)) {
        ctx.JSON(iris.StatusUnauthorized, iris.Map{
            "error": "wrong username or password",
        })
        return
    }

    now := time.Now()
    notification := &util.Notification{
        AppId: data.AppId,
        AppKey: data.AppKey,
        Id: uuid.NewV4().String(),
        Qos: 2,
        Retain: 1,
        Topic: data.Topic,
        Message: data.Message,
        LastUpdated: now,
        Timestamp: now.Unix(),
    }
    fmt.Printf("%+v\n", *notification)

    job := util.Job{
        Payload: notification,
        Do: func(action util.Action) {
            if err := action.Notify(); err != nil {
                fmt.Errorf("fail to notify %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job

    increment("dolphin.api.v1.notification")
    ctx.JSON(iris.StatusOK, iris.Map{
        "status": "success",
    })
}

func message(ctx *iris.Context) {
    page, err := ctx.ParamInt("page")
    if err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }

    pageSize, err := ctx.ParamInt("pageSize")
    if err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }
    if pageSize >= 100 {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": "page size should be 0 < pageSize < 100",
        })
        return
    }

    total, err := util.Total()
    if err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }

    if (total/pageSize + 1) < page {
        ctx.JSON(iris.StatusOK, iris.Map{
            "total": total,
            "page": page,
            "pageSize": pageSize,
            "messages": nil,
        })
        return
    }

    size := total - (page-1) * pageSize
    if size >= pageSize {
        size = pageSize
    }
    rows := make([]util.Notification, size)
    if err := util.List(rows, page, pageSize); err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }

    increment("dolphin.api.v1.message")
    ctx.JSON(iris.StatusOK, iris.Map{
        "total": total,
        "page": page,
        "pageSize": pageSize,
        "messages": rows,
    })
}

func application(ctx *iris.Context) {
    var rows []string
    var err error

    if rows, err = util.Application(); err != nil {
        ctx.JSON(iris.StatusInternalServerError, iris.Map{
            "error": err,
        })
        return
    }

    increment("dolphin.api.v1.application")
    ctx.JSON(iris.StatusOK, iris.Map{
        "applications": rows,
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

func increment(text string) {
    if config.Debug {
        return
    }
    c, err := statsd.New(statsd.Address(config.StatsdServer))
    if err != nil {
        fmt.Errorf("fail to initialize statsd %+v\n", err)
    } else {
        // Increment a counter.
        c.Increment(text)
    }
    defer c.Close()
}
