#encoding: utf-8
import uuid

from django.conf import settings

import pika

from base.management import RabbitMQWatcher


class Command(RabbitMQWatcher):
    help = 'Agent Register RPC Server'

    def get_config(self, *args, **options):
        return settings.RABBITMQ['QUEUE_RPC'], settings.RABBITMQ['EXCHANGE_RPC'], settings.RABBITMQ['ROUTINGKEY_RPC']

    def dispatch(self, channel, method, props, body):
        reply_to = getattr(props, 'reply_to', None)
        correlation_id = getattr(props, 'correlation_id', None)
        if reply_to is not None and correlation_id is not None:
            channel.basic_publish(exchange='',
                                     routing_key=reply_to,
                                     properties=pika.BasicProperties(correlation_id=correlation_id),
                                     body=str(uuid.uuid1()))

        channel.basic_ack(delivery_tag=method.delivery_tag)
