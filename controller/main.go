package main

import (
    "fmt"
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
            Id: msg.Payload(),
            LastUpdated: now,
            Timestamp: now.Unix(),
        },
        Do: func(action util.Action) {
            if err := action.Ack(); err != nil {
                fmt.Printf("fail to notify %+v, error %+v\n", action, err)
            }
        },
    }
    jobQueue <- job
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
    /*opts := MQTT.NewClientOptions().AddBroker("ws://54.223.124.84:80/mqtt")
    opts.SetClientID("monitor")
    opts.SetUsername("monitor")
    opts.SetPassword("7C7DC73CDFAB3838C5E2CE82E1BFC")*/

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
}