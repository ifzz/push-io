package main

import (
    "fmt"
    "gopkg.in/alexcesaro/statsd.v2"
    //import the Paho Go MQTT library
    MQTT "github.com/eclipse/paho.mqtt.golang"
    "os"
    "./util"
    "time"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
    fmt.Printf("TOPIC: %s\n", msg.Topic())
    fmt.Printf("MSG: %s\n", msg.Payload())
    now := time.Now()
    job := util.Job{
        Payload: &util.Notification{
            Id: string(msg.Payload()),
            LastUpdated: now,
            Timestamp: now.Unix(),
        },
        Do: func(action util.Action) {
            if err := action.PushAck(); err != nil {
                fmt.Printf("fail to handle push ack event %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job

    increment("dolphin.topic.ack")
}

var config = util.NewConfig()
var jobQueue chan util.Job

func init() {
    jobQueue = make(chan util.Job, config.MaxQueue)
    dispatcher := util.NewDispatcher(jobQueue)
    dispatcher.Run()
}

func main() {
    //create a ClientOptions struct setting the broker address, clientid, turn
    //off trace output and set the default message handler
    opts := MQTT.NewClientOptions()
    opts.AddBroker(config.MqttServer)
    opts.SetClientID(config.ClientId)
    opts.SetUsername(config.Username)
    opts.SetPassword(config.Password)
    opts.SetCleanSession(false)
    opts.SetDefaultPublishHandler(f)

    //create and start a client using the above ClientOptions
    c := MQTT.NewClient(opts)
    if token := c.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    //subscribe to the topic /go-mqtt/sample and request messages to be delivered
    //at a maximum qos of zero, wait for the receipt to confirm the subscription
    if token := c.Subscribe("ack/#", 2, nil); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }
    fmt.Println(time.Now(), "running...")

    select {}
}

func increment(text string) {
    if config.Debug {
        return
    }
    c, err := statsd.New(statsd.Address(config.StatsdServer))
    if err != nil {
        fmt.Printf("fail to initialize statsd %+v\n", err)
    } else {
        // Increment a counter.
        c.Increment(text)
    }
    defer c.Close()
}
