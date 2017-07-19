package communicate

import (
    "fmt"
    "errors"
    "encoding/json"

    "github.com/streadway/amqp"

    "config"
    "entity"
    "log"
)

func PublishMessage(cxt *config.Context) (err error) {
    // channel, err := cxt.MQ.Conn.Channel()
    // defer channel.Close()
    channel := cxt.MQ.Channel
    for msg := range cxt.Channels.Message {
        err = channel.Publish(msg.Exchange,
                                        msg.RoutingKey,
                                        false,
                                        false,
                                        amqp.Publishing{
                                            ContentType: "application/json",
                                            Body: []byte(msg.Body),
                                        })

        cxt.Channels.Log <- fmt.Sprintf("Publish Message:%s", string(msg.Body))
        if nil != err {
            cxt.Channels.Log <- fmt.Sprintf("Error PublishMessage: %s", err)
        }
    }
    return
}

func GetCommands(cxt *config.Context) (err error) {
    for msg := range cxt.MQ.CommandConsume {
        msg.Ack(false)
        var cmd entity.Command
        err := json.Unmarshal(msg.Body, &cmd)
        if nil == err {
            cxt.Channels.Command <- cmd
            cxt.Channels.Log <- fmt.Sprintf("Recive Command:%s", string(msg.Body))
        } else {
            cxt.Channels.Log <- fmt.Sprintf("Error Parse Command:%s, error:%s", string(msg.Body), err)
        }
    }
    return
}


func Rpc(cxt *config.Context, msg []byte) (result []byte, err error) {
    // channel, err := cxt.MQ.Conn.Channel()
    // defer channel.Close()
    channel := cxt.MQ.Channel
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
                                Body: []byte(msg),
                            })

    if nil != err {
        return
    }

    for reply := range consume {
        log.Println("RPC Replay:", string(reply.Body))
        return reply.Body, nil
    }

    return nil, errors.New("Unknow")
}
