package config

import (
    "os"
    "net"
    "errors"
    "path/filepath"
    "strings"
    "log"
    "os/signal"
    "time"
    "fmt"

    "gopkg.in/ini.v1"
    "github.com/streadway/amqp"

    "entity"
)

type AgentConfig struct {
    UUID string
    IP string
    MAC string
    Hostname string
}


type ServerConfig struct {
    IP string
    Rabbitmq_port int
    Rabbitmq_vhost string
    Rabbitmq_user string
    Rabbitmq_password string

    Rabbitmq_exchange_cmd string
    Rabbitmq_routingkey_cmd string
    Rabbitmq_queue_cmd string

    Rabbitmq_exchange_heartbeat string
    Rabbitmq_routingkey_heartbeat string
    Rabbitmq_queue_heartbeat string

    Rabbitmq_exchange_cmd_result string
    Rabbitmq_routingkey_cmd_result string
    Rabbitmq_queue_cmd_result string

    Rabbitmq_exchange_log string
    Rabbitmq_routingkey_log string
    Rabbitmq_queue_log string

    Rabbitmq_exchange_rpc string
    Rabbitmq_routingkey_rpc string
}


type PluginConfig struct {
    Name string `ini:"name"`
    Enable bool `ini:"enable"`
    Trigger string `ini:"trigger"`
    Cmd_Types []int `ini:"cmd_types"`
}


type Context struct {
    Home string
    AConfig *AgentConfig
    SConfig *ServerConfig
    PConfigList []*PluginConfig
    Channels struct {
        Terminate chan string
        Log chan string
        Message chan entity.Message
        Command chan entity.Command
    }
    MQ struct {
        Conn *amqp.Connection
        Channel *amqp.Channel
        CommandConsume <-chan amqp.Delivery
    }
}


func Init(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt

    cxt, err = initChannels(cxt)
    if nil != err {
        return
    }
    go func() {
        for e := range cxt.Channels.Log {
            log.Println(e)
        }
    }()

    cxt, err = initHome(cxt)
    if nil != err {
        return
    }

    cxt, err = initAgentConfig(cxt)
    if nil != err {
        return
    }

    cxt, err = initServerConfig(cxt)
    if nil != err {
        return
    }

    cxt, err = initMQConnect(cxt)
    if nil != err {
        return
    }

    cxt, err = initAgentUUID(cxt)
    if nil != err {
        return
    }

    cxt, err = initMQDeclare(cxt)
    if nil != err {
        return
    }

    cxt, err = iniPluginsConfig(cxt)
    if nil != err {
        return
    }

    var ch chan os.Signal = make(chan os.Signal)
    signal.Notify(ch, os.Interrupt, os.Kill)

    go func() {
        signal := <- ch
        cxt.Channels.Terminate <- signal.String()
    }()

    return
}

func initChannels(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    cxt.Channels.Terminate = make(chan string, 10)
    cxt.Channels.Log = make(chan string, 100)
    cxt.Channels.Message = make(chan entity.Message, 1000)
    cxt.Channels.Command = make(chan entity.Command, 10)
    cxt.Channels.Log <- "Success Init Channels"
    return
}

func initHome(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    dir, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
    if nil != err {
        return
    }
    cxt.Home = strings.Replace(dir, "\\", "/", -1)
    cxt.Channels.Log <- "Success Init Home"
    return
}


func getIpMac() (string, string, error) {
    inters, err := net.Interfaces()
    if nil != err {
        return "", "", err
    }

    for _, inter := range inters {
        if "" != inter.HardwareAddr.String() {
            addrs, _ := inter.Addrs()
            for _, addr := range addrs {
                if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && nil != ipnet.IP.To4() {
                    return ipnet.IP.String(), inter.HardwareAddr.String(), nil
                }
            }
        }
    }
    return "", "", errors.New("Error Get IP and Mac")
}

func initAgentConfig(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    path := filepath.Join(cxt.Home, "etc", "agent.ini")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return cxt, err
    }

    cfg, err := ini.Load(path)
    if nil != err {
        return
    }

    section, err := cfg.GetSection("AGENT")
    if nil != err {
        return
    }

    uuid := section.Key("uuid").String()

    ip, mac, err := getIpMac()
    if nil != err {
        return
    }
    hostname, _ := os.Hostname()
    cxt.AConfig = &AgentConfig{uuid, ip, mac, hostname}
    cxt.Channels.Log <- "Success Init Agent Config"
    return
}


func initServerConfig(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    path := filepath.Join(cxt.Home, "etc", "server.ini")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return cxt, err
    }

    cfg, err := ini.Load(path)
    if nil != err {
        return
    }

    section, err := cfg.GetSection("SERVER")
    if nil != err {
        return
    }

    ip := section.Key("ip").String()

    cxt.SConfig = &ServerConfig{ip, RABBITMQ_PORT, RABBITMQ_VHOST, RABBITMQ_USER, RABBITMQ_PASSWORD,
                            RABBITMQ_EXCHANGE_CMD, RABBITMQ_ROUTINGKEY_CMD, RABBITMQ_QUEUE_CMD,
                            RABBITMQ_EXCHANGE_HEARTBEAT, RABBITMQ_ROUTINGKEY_HEARTBEAT, RABBITMQ_QUEUE_HEARTBEAT,
                            RABBITMQ_EXCHANGE_CMD_RESULT, RABBITMQ_ROUTINGKEY_CMD_RESULT, RABBITMQ_QUEUE_CMD_RESULT,
                            RABBITMQ_EXCHANGE_LOG, RABBITMQ_ROUTINGKEY_LOG, RABBITMQ_QUEUE_LOG,
                            RABBITMQ_EXCHANGE_RPC, RABBITMQ_ROUTINGKEY_RPC,
                        }

    cxt.Channels.Log <- "Success Init Server Config"
    return
}

func initMQConnect(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cxt.SConfig.Rabbitmq_user, cxt.SConfig.Rabbitmq_password,
                        cxt.SConfig.IP, cxt.SConfig.Rabbitmq_port, cxt.SConfig.Rabbitmq_vhost)
    cxt.MQ.Conn, err = amqp.Dial(url)
    if nil != err {
        return
    }

    cxt.MQ.Channel, err = cxt.MQ.Conn.Channel()
    if nil != err {
        return
    }

    err = cxt.MQ.Channel.Qos(1, 0, false)
    if nil != err {
        return
    }

    cxt.Channels.Log <- "Success Init MQ Connection"
    return
}

func initAgentUUID(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    if "" == cxt.AConfig.UUID {
        uuid, err := register(&cxt)
        if nil != err {
            return cxt, err
        }
        cxt.AConfig.UUID = uuid
        writeAgentUUID(&cxt, uuid)
    }
    cxt.SConfig.Rabbitmq_routingkey_cmd = fmt.Sprintf(cxt.SConfig.Rabbitmq_routingkey_cmd, cxt.AConfig.UUID)
    cxt.SConfig.Rabbitmq_queue_cmd = fmt.Sprintf(cxt.SConfig.Rabbitmq_queue_cmd, cxt.AConfig.UUID)
    cxt.SConfig.Rabbitmq_routingkey_log = fmt.Sprintf(cxt.SConfig.Rabbitmq_routingkey_log, cxt.AConfig.UUID)
    cxt.SConfig.Rabbitmq_queue_log = fmt.Sprintf(cxt.SConfig.Rabbitmq_queue_log, cxt.AConfig.UUID)
    cxt.Channels.Log <- fmt.Sprintf("Success Init AgentID:%s", cxt.AConfig.UUID)
    return cxt, nil

}

func register(cxt *Context) (result string, err error) {
    channel, err := cxt.MQ.Conn.Channel()
    if nil != err {
        return
    }
    defer channel.Close()
    queue, err := channel.QueueDeclare("", false, false, true, false, nil)
    if nil != err {
        return
    }

    consume, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
    if nil != err {
        return
    }

    err = channel.Publish(cxt.SConfig.Rabbitmq_exchange_rpc,
                            cxt.SConfig.Rabbitmq_routingkey_rpc,
                            false,
                            false,
                            amqp.Publishing {
                                ContentType: "application/json",
                                CorrelationId: "0",
                                ReplyTo: queue.Name,
                                Body: []byte("register"),
                            })

    if nil != err {
        return
    }

    for reply := range consume {
        reply.Ack(false)
        return string(reply.Body), nil
    }

    return "", errors.New("Unknow")
}

func writeAgentUUID(cxt *Context, uuid string) (err error) {
    path := filepath.Join(cxt.Home, "etc", "agent.ini")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return err
    }

    cfg, err := ini.Load(path)
    if nil != err {
        return
    }

    section, err := cfg.GetSection("AGENT")
    if nil != err {
        return
    }

    section.Key("uuid").SetValue(uuid)

    cfg.SaveTo(path)

    cxt.Channels.Log <- "Success Write Agent UUID to Config"
    return nil
}

func initMQDeclare(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt

    cxt.MQ.Channel.ExchangeDeclare(cxt.SConfig.Rabbitmq_exchange_heartbeat, "direct", true, false, false, false, nil)
    cxt.MQ.Channel.ExchangeDeclare(cxt.SConfig.Rabbitmq_exchange_cmd, "direct", true, false, false, false, nil)
    cxt.MQ.Channel.ExchangeDeclare(cxt.SConfig.Rabbitmq_exchange_cmd_result, "direct", true, false, false, false, nil)
    cxt.MQ.Channel.ExchangeDeclare(cxt.SConfig.Rabbitmq_exchange_log, "direct", true, false, false, false, nil)
    cxt.MQ.Channel.ExchangeDeclare(cxt.SConfig.Rabbitmq_exchange_rpc, "direct", true, false, false, false, nil)

    queue, err := cxt.MQ.Channel.QueueDeclare(cxt.SConfig.Rabbitmq_queue_cmd, true, false, false, false, nil)
    if nil != err {
        return
    }

    err = cxt.MQ.Channel.QueueBind(queue.Name, cxt.SConfig.Rabbitmq_routingkey_cmd, cxt.SConfig.Rabbitmq_exchange_cmd, false, nil)
    if nil != err {
        return
    }

    cxt.MQ.CommandConsume, err = cxt.MQ.Channel.Consume(queue.Name, cxt.AConfig.UUID, false, false, false, false, nil)

    if nil != err {
        return
    }

    cxt.Channels.Log <- "Success Init MQ Connection"
    return
}


func iniPluginsConfig(orig_cxt Context) (cxt Context, err error) {
    cxt = orig_cxt
    path := filepath.Join(cxt.Home, "etc", "plugins.ini")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return cxt, err
    }

    cfg, err := ini.Load(path)
    if nil != err {
        return
    }

    var list []*PluginConfig= make([]*PluginConfig, 0)

    sections := cfg.Sections()
    for _, section := range sections {
        name := section.Key("name").String()
        if "" == name {
            continue
        }
        plugin := &PluginConfig{}
        err = section.MapTo(plugin)
        if nil != err {
            cxt.Channels.Log <- fmt.Sprintf("Error Parse Plugin: %v")
            continue
        }
        list = append(list, plugin)
    }
    cxt.PConfigList = list
    cxt.Channels.Log <- "Success Init Plugins Config"
    return
}

func Destory(cxt Context) (err error) {
    time.Sleep(1 * time.Second)
    close(cxt.Channels.Command)
    close(cxt.Channels.Message)

    close(cxt.Channels.Terminate)
    close(cxt.Channels.Log)

    if nil != cxt.MQ.Channel {
        cxt.MQ.Channel.Close()
    }

    if nil != cxt.MQ.Conn {
        cxt.MQ.Conn.Close()
    }
    return err
}
