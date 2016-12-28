package main

import (
    "fmt"
    "github.com/kataras/iris"
    "github.com/parnurzeal/gorequest"
    "github.com/iris-contrib/middleware/cors"
    "github.com/iris-contrib/middleware/logger"
    "strings"
    "strconv"
    "time"
    "github.com/satori/go.uuid"
    "./util"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

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

    iris.Listen(":8080")
}

func index(ctx *iris.Context) {
    ctx.MustRender("index.html", struct {}{})
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

    found := false
    for _, app := range key.Apps {
        if (app.AppId == data.AppId && app.AppKey == data.AppKey) {
            found = true
        }
    }
    if (!found) {
        ctx.EmitError(iris.StatusUnauthorized)
        return
    }

    now := time.Now()
    notification := &util.Notification{
        AppId: data.AppId,
        AppKey: data.AppKey,
        Id: uuid.NewV4().String(),
        Timestamp: now.Unix(),
        Created: now,
        Qos: 2,
        Retain: 1,
        Topic: data.Topic,
        Message: data.Message,
        Success: false,
    }
    ctx.Log("%+v\n", *notification)

    job := util.Job{
        Payload: notification,
        Do: func(action util.Action) {
            if err := action.Save(); err != nil {
                ctx.Log("fail to save %+v, error %+v\n", action, err)
            }

            if err := action.Notify(); err != nil {
                ctx.Log("fail to notify %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job

    ctx.Text(iris.StatusOK, "ok")
}

func server(ctx *iris.Context) {
    request := gorequest.New().SetBasicAuth(config.Username, config.Password)
    request.SetDebug(config.Debug)

    count := MaxInt
    host := ""
    for _, value := range config.Cluster {
        //fmt.Println(value)
        _, body, errs := request.Get(fmt.Sprintf("http://%s/api/stats", value)).End()
        if (len(errs) > 0) {
            fmt.Printf("error %+v\n", errs[0])
            continue
        }
        //fmt.Println(body);
        ctx.Log("%s:%s\n", value, body)

        result := strings.Split(body, ",")
        clients_count := 0
        for _, stat := range result {
            pair := strings.Split(stat, ":")
            if strings.Contains(pair[0], "clients") && strings.Contains(pair[0], "count") {
                //fmt.Println(stat)
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
}
