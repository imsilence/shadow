package test

import (
    "fmt"
    "log"

    "config"
    "plugins"
    "entity"
)

const name string = "test"

type TestPlugin struct {
    cxt *config.Context
}

func (t *TestPlugin) Init(cxt *config.Context) (err error) {
    t.cxt = cxt
    t.cxt.Channels.Log <- fmt.Sprintf("Init Plugin %s", name)
    return nil
}

func (t *TestPlugin) Run() (err error) {
    return nil
}


func (t *TestPlugin) Call(command *entity.Command) (result *entity.CommandResult, err error) {
    if command.Type == 10001 {
        panic(fmt.Sprintf("Error Type:%d", command.Type))
    }
    log.Printf("%s Run Command:%x\n", name, *command)
    return &entity.CommandResult{}, nil
}

func (t *TestPlugin) Destory() (err error) {
    return nil
}

func init() {
    plugins.Register(name, new(TestPlugin))
}
