package plugins

import (
    "fmt"
    "encoding/json"
    "config"
    "entity"
)

type Plugin interface {
    Init(cxt *config.Context) error
    Run() error
    Call(command *entity.Command) (*entity.CommandResult, error)
    Destory() error
}

var all_plugins map[string]Plugin = make(map[string]Plugin)
var enabled_plugins map[string]Plugin = make(map[string]Plugin)
var commands map[int]Plugin = make(map[int]Plugin)

func Register(name string, plugin Plugin) error {
    if _, ok := all_plugins[name]; ok {
        return fmt.Errorf("Plugin Already Exists:%s", name)
    }
    all_plugins[name] = plugin
    return nil
}

func Init(cxt *config.Context) error {
    for _, pconfig := range cxt.PConfigList {
        if pconfig.Enable {
            plugin, ok := all_plugins[pconfig.Name]
            if !ok {
                cxt.Channels.Log <- fmt.Sprintf("Plugin %s is Not Found", pconfig.Name)
                continue
            }
            plugin.Init(cxt)
            enabled_plugins[pconfig.Name] = plugin
            for _, cmd_type := range pconfig.Cmd_Types {
                commands[cmd_type] = plugin
            }
            cxt.Channels.Log <- fmt.Sprintf("Plugin %s is Enabled", pconfig.Name)
        }
    }
    return nil
}

func Run() error {
    for _, plugin := range enabled_plugins {
        go plugin.Run()
    }
    return nil
}

func Call(cxt *config.Context) error {
    go func() {
        for command := range cxt.Channels.Command {
            plugin, ok := commands[command.Type]
            var result *entity.CommandResult
            if !ok {
                reason := fmt.Sprintf("Cmd Type %d is Not Found", command.Type)
                result = &entity.CommandResult{ID: command.ID, Reason: reason}
                cxt.Channels.Log <- reason
            } else {
                func() {
                    defer func() {
                        if p := recover(); nil != p {
                            reason := fmt.Sprintf("Error: %v", p)
                            result = &entity.CommandResult{ID: command.ID, Reason: reason}
                            cxt.Channels.Log <- reason
                        }
                    }()
                    result, _ = plugin.Call(&command)
                }()
            }
            jresult, _ := json.Marshal(*result)
            msg := entity.Message{cxt.SConfig.Rabbitmq_exchange_cmd_result, cxt.SConfig.Rabbitmq_routingkey_cmd_result, jresult}
            cxt.Channels.Message <- msg
        }
    }()
    return nil
}


func Destory() error {
    for _, plugin := range enabled_plugins {
        plugin.Destory()
    }
    return nil
}
