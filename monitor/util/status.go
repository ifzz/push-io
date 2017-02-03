package util

import (
    "gopkg.in/alexcesaro/statsd.v2"
    "github.com/parnurzeal/gorequest"
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/spf13/viper"
    "bytes"
    "strconv"
)

type Node struct {
    Name string `json:"name"`
    Address string
    ClusterStatus string `json:"cluster_status"`
    Clients int
    TotalMemory string `json:"total_memory"`
    UsedMemory string `json:"used_memory"`
    Load1 string `json:"load1"`
    Load5 string `json:"load5"`
    Load15 string `json:"load15"`
}

const API_NODES = "http://%s:%d/api/nodes"

const API_STATS = "http://%s:%d/api/stats"

const BROKER_NODES = "io.gf.com.cn:nodes"

const BROKER_STATS = "io.gf.com.cn:stats:%s"

const FIELD_CLIENTS = "clients"

const FIELD_STATUS = "status"

const FIELD_ADDRESS = "address"

const FIELD_TOTAL_MEMORY = "total_memory"

const FIELD_USED_MEMORY = "used_memory"

const FIELD_LOAD_1 = "load1"

const FIELD_LOAD_5 = "load5"

const FIELD_LOAD_15 = "load15"

func Update() {
    nodes := queryNodes()
    for _, node := range nodes {
        tokens := strings.Split(node.Name, "@")
        node.Address = tokens[1]
        queryStats(&node)
        fmt.Println(node)
        save(&node)
    }
}

// get the status of node
func queryStats(node *Node) {
    result := query(fmt.Sprintf(API_STATS, node.Address, config.BrokerPort))
    if len(result) > 0 {
        //fmt.Println(result)
        viper.ReadConfig(bytes.NewBuffer([]byte(result)))
        node.Clients = viper.GetInt("clients/count")
    }
}

// get the list of nodes in cluster
func queryNodes() []Node {
    var nodes []Node
    result := query(fmt.Sprintf(API_NODES, config.BrokerHost, config.BrokerPort))
    if len(result) > 0 {
        //fmt.Println(result)
        json.Unmarshal([]byte(result), &nodes)
    }
    return nodes
}

func query(url string) string {
    request := gorequest.New().SetBasicAuth(config.Username, config.Password)
    if len(config.Proxy) > 0 {
        request.Proxy(config.Proxy)
    }

    resp, body, errs := request.Get(url).End()

    if resp.StatusCode != http.StatusOK || len(errs) > 0 {
        fmt.Printf("%d: %+v\n", resp.StatusCode, errs[0])
        return ""
    }

    return body
}

func save(node *Node) {
    conn := redisPool.Get()
    defer conn.Close()

    // save the list of node
    if _, err := conn.Do("sadd", BROKER_NODES, node.Address); err != nil {
        fmt.Printf("error %+v\n", err)
    }
    fmt.Printf("%+v\n", *node)

    // save the status of node
    host := getPublicAddr(node.Address)
    if _, err := conn.Do("HMSET", fmt.Sprintf(BROKER_STATS, node.Address),
        FIELD_CLIENTS, node.Clients,
        FIELD_STATUS, node.ClusterStatus,
        FIELD_ADDRESS, host,
        FIELD_TOTAL_MEMORY, node.TotalMemory,
        FIELD_USED_MEMORY, node.UsedMemory,
        FIELD_LOAD_1, node.Load1,
        FIELD_LOAD_5, node.Load5,
        FIELD_LOAD_15, node.Load15,
    ); err != nil {
        fmt.Printf("error %+v\n", err)
    }

    gauge(host, FIELD_CLIENTS, node.Clients)
    gauge(host, FIELD_STATUS, node.ClusterStatus)
    gauge(host, FIELD_LOAD_1, node.Load1)
    gauge(host, FIELD_LOAD_5, node.Load5)
    gauge(host, FIELD_LOAD_15, node.Load15)
    gauge(host, FIELD_TOTAL_MEMORY, node.TotalMemory)
    gauge(host, FIELD_USED_MEMORY, node.UsedMemory)
}

const IP_CONFIG_PATH = "."
const IP_CONFIG_NAME = "ip"

func getPublicAddr(ip string) string {
    viper.SetConfigName(IP_CONFIG_NAME)
    viper.AddConfigPath(IP_CONFIG_PATH)
    err := viper.ReadInConfig()

    if err != nil {
        panic("IP configuration file not found")
    }

    return viper.GetString(ip)
}

func gauge(host string, field string, value interface{}) {
    if config.Debug {
        return
    }
    c, err := statsd.New(statsd.Address(config.StatsdServer))
    if err != nil {
        fmt.Printf("fail to initialize statsd %+v\n", err)
    } else {
        // Increment a counter.
        c.Gauge("dolphin." + convert(host) + "." + field, toValue(field, value))
    }
    defer c.Close()
}

func convert(text string) string {
    return strings.Replace(text, ".", "_", -1)
}

func toValue(field string, text interface{}) int {
    if field == FIELD_STATUS {
        if text == "Running" {
            return 1
        } else {
            return 0
        }
    }

    if field == FIELD_LOAD_1 || field == FIELD_LOAD_5 || field == FIELD_LOAD_15 {
        value, err := strconv.ParseFloat(text.(string), 32)
        if err != nil {
            fmt.Printf("fail to parse field [load 1/5/15], %+v\n", err)
        } else {
            return int(value * 100)
        }
    }

    if field == FIELD_TOTAL_MEMORY || field == FIELD_USED_MEMORY {
        valueText := strings.Replace(text.(string), "M", "", -1)
        value, err := strconv.ParseFloat(valueText, 32)
        if err != nil {
            fmt.Printf("fail to parse field [total/used memory], %+v\n", err)
        } else {
            return int(value)
        }
    }

    if field == FIELD_CLIENTS {
        return text.(int)
    }

    return -1
}
