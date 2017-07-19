package main

import (
    "log"
    "fmt"

    "config"
    "plugins"

    "communicate"
)

func main() {

    var cxt config.Context
    cxt, err := config.Init(cxt)
    if nil != err {
        log.Fatal(err)
    }
    cxt.Channels.Log <- fmt.Sprintf("Start Agent On: %s", cxt.Home)

    go communicate.PublishMessage(&cxt)

    plugins.Init(&cxt)
    plugins.Run()
    plugins.Call(&cxt)

    go communicate.GetCommands(&cxt)

    e := <- cxt.Channels.Terminate

    plugins.Destory()
    config.Destory(cxt)

    if "interrupt" != e {
        cxt.Channels.Log <- e
        log.Println(e)
    }
    log.Println("Exist Agent ")
}
