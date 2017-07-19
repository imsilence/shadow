package heartbeat

import (
    "os"
    "runtime"
    "time"
    "fmt"
    "encoding/json"

    "config"
    "plugins"
    "entity"
    "utils"
)

const name string = "heartbeat"
const interval_heartbeat int = 10

type HeartbeatPlugin struct {
    cxt *config.Context
}

func (h *HeartbeatPlugin) Init(cxt *config.Context) (err error) {
    h.cxt = cxt
    h.cxt.Channels.Log <- fmt.Sprintf("Init Plugin %s", name)
    return nil
}

func (h *HeartbeatPlugin) Run() (err error) {
    ticker := time.NewTicker(time.Duration(interval_heartbeat) * time.Second)
    for _ = range ticker.C {
        interfaces, _ := utils.GetInterfaces()
        agent := entity.Agent{
            UUID: h.cxt.AConfig.UUID,
            Hostname: h.cxt.AConfig.Hostname,
            PID: os.Getpid(),
            OS: runtime.GOOS,
            Arch: runtime.GOARCH,
            Interfaces: interfaces,
            Time: time.Now(),
        }
        body, err := json.Marshal(agent)
        if nil != err {
            h.cxt.Channels.Log <- fmt.Sprintf("Error Json Marshal Heartbeat: %s", err)
            continue
        }
        h.cxt.Channels.Message <- entity.Message{h.cxt.SConfig.Rabbitmq_exchange_heartbeat,
                                                    h.cxt.SConfig.Rabbitmq_routingkey_heartbeat,
                                                    body}
    }
    return nil
}


func (h *HeartbeatPlugin) Call(command *entity.Command) (result *entity.CommandResult, err error) {
    return nil, nil
}

func (h *HeartbeatPlugin) Destory() (err error) {
    return nil
}

func init() {
    plugins.Register(name, new(HeartbeatPlugin))
}
