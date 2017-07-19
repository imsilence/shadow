#encoding: utf-8

import json

from django.conf import settings

from base.management import RabbitMQWatcher
from agent.models import Client


class Command(RabbitMQWatcher):
    help = 'Agent Heartbeat'

    def get_config(self):
        return settings.RABBITMQ['QUEUE_HEARTBEAT'], settings.RABBITMQ['EXCHANGE_HEARTBEAT'], settings.RABBITMQ['ROUTINGKEY_HEARTBEAT']

    def dispatch(self, channel, method, props, body):
        json_body = json.loads(body.decode('utf-8'))
        Client.create_or_replace(json_body)
        channel.basic_ack(delivery_tag = method.delivery_tag)
