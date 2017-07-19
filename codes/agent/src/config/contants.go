package config


/*
RabbitMQ
*/
const RABBITMQ_PORT int = 5672
const RABBITMQ_VHOST string = "shadow"
const RABBITMQ_USER string = "shadow"
const RABBITMQ_PASSWORD string = "shadow@2017"

const RABBITMQ_EXCHANGE_CMD string = "cmd.direct"
const RABBITMQ_ROUTINGKEY_CMD string = "cmd.%s"
const RABBITMQ_QUEUE_CMD string = "cmd.%s"

const RABBITMQ_EXCHANGE_HEARTBEAT string = "heartbeat.direct"
const RABBITMQ_ROUTINGKEY_HEARTBEAT string = "heartbeat"
const RABBITMQ_QUEUE_HEARTBEAT string = "heartbeat"

const RABBITMQ_EXCHANGE_CMD_RESULT string = "cmd.result.direct"
const RABBITMQ_ROUTINGKEY_CMD_RESULT string = "cmd.result"
const RABBITMQ_QUEUE_CMD_RESULT string = "cmd.result"

const RABBITMQ_EXCHANGE_LOG string = "log.direct"
const RABBITMQ_ROUTINGKEY_LOG string = "log.%s"
const RABBITMQ_QUEUE_LOG string = "log.%s"

const RABBITMQ_EXCHANGE_RPC string = "rpc.direct"
const RABBITMQ_ROUTINGKEY_RPC string = "rpc"

/*
*/
const TASK_STATUS_RUNNING int = 1
const TASK_STAUS_ERROR int = 2
const TASK_STATUS_SUCCESS int = 3
