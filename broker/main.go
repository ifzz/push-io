package main

import (
    "github.com/kataras/iris"
    "github.com/iris-contrib/middleware/cors"
    "github.com/iris-contrib/middleware/logger"
    "time"
    "github.com/satori/go.uuid"
    "./util"
)

var key = util.InitKey()
var config = util.InitConfig()
var jobQueue chan util.Job

func init() {
    jobQueue = make(chan util.Job, config.MaxQueue)
    dispatcher := util.NewDispatcher(jobQueue)
    dispatcher.Run()

    irisLogger := logger.New(logger.Config{
        // Status displays status code
        Status: true,
        // IP displays request's remote address
        IP: true,
        // Method displays the http method
        Method: true,
        // Path displays the request path
        Path: true,
    })
    iris.Use(irisLogger)

    iris.Config.IsDevelopment = config.Debug // reloads the templates on each request, defaults to false
    iris.Config.Gzip  = true // compressed gzip contents to the client, the same for Serializers also, defaults to false
    iris.Config.Charset = "UTF-8" // defaults to "UTF-8", the same for Serializers also

    iris.Static("/scripts", "./templates/scripts/", 1)
    iris.Static("/styles", "./templates/styles/", 1)
    iris.Static("/fonts", "./templates/fonts/", 1)
    iris.Static("/images", "./templates/images/", 1)

    iris.Use(cors.Default())
}

func main() {
    iris.Get("/", index)

    iris.Get("/api/v1/server", server)

    iris.Post("/api/v1/notification", notification)

    iris.Post("/api/v1/ack", ack)

    iris.Listen(":8080")
}

func index(ctx *iris.Context) {
    ctx.MustRender("index.html", struct {}{})
}

func ack(ctx *iris.Context) {
    type Data struct {
        Id     string `json:"id"`
        AppId  string `json:"appId"`
        AppKey string `json:"appKey"`
    }
    data := &Data{}
    if err := ctx.ReadJSON(data); err != nil {
        ctx.Log("%+v\n", err)
        ctx.EmitError(iris.StatusInternalServerError)
        return
    }

    if (!IsAuthorized(data.AppId, data.AppKey)) {
        ctx.EmitError(iris.StatusUnauthorized)
        return
    }

    job := util.Job{
        Payload: &util.Notification{
            Id: data.Id,
            AppId: data.AppId,
            AppKey: data.AppKey,
        },
        Do: func(action util.Action) {
            if err := action.Update(); err != nil {
                ctx.Log("fail to update %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job

    ctx.Text(iris.StatusOK, "ok")
}

func IsAuthorized(appId string, appKey string) bool {
    found := false
    for _, app := range key.Apps {
        if (app.AppId == appId && app.AppKey == appKey) {
            found = true
        }
    }
    return found
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
        ctx.Log("%+v\n", err)
        ctx.EmitError(iris.StatusInternalServerError)
        return
    }

    if (!IsAuthorized(data.AppId, data.AppKey)) {
        ctx.EmitError(iris.StatusUnauthorized)
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
    ctx.Log("%+v\n", *notification)

    job := util.Job{
        Payload: notification,
        Do: func(action util.Action) {
            if err := action.Notify(); err != nil {
                ctx.Log("fail to notify %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job

    ctx.Text(iris.StatusOK, "ok")
}

func server(ctx *iris.Context) {
    host, count := util.GetServer()

    job := util.Job{
        Payload: &util.Notification{},
        Do: func(action util.Action) {
            action.GetStatus()
        },
    }
    jobQueue <- job

    ctx.JSON(iris.StatusOK, iris.Map{
        "host": host,
        "count": count,
    })
}
