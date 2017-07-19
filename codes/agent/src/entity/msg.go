package entity

type Message struct {
    Exchange string
    RoutingKey string
    Body []byte
}
