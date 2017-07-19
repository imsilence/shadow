#encoding: utf-8
from django.core.management import BaseCommand
from django.conf import settings
import pika

class RabbitMQWatcher(BaseCommand):
    help = 'RabbitMQ Watcher'

    def get_config(self):
        raise BaseException('function get_config is undefined')
        return 'queue', 'exchange', 'routing_key'


    def handle(self, *args, **options):
        queue, exchange, routing_key = self.get_config()

        credentials = pika.credentials.PlainCredentials(settings.RABBITMQ['USER'], settings.RABBITMQ['PASSWORD'])
        parameters = pika.ConnectionParameters(host=settings.RABBITMQ['HOST'], port=settings.RABBITMQ['PORT'], virtual_host=settings.RABBITMQ['VHOST'], credentials=credentials)
        connection = pika.BlockingConnection(parameters)
        channel = connection.channel()

        channel.queue_declare(queue, passive=False, durable=True, exclusive=False, auto_delete=False, arguments=None)
        channel.queue_bind(queue, exchange, routing_key, None)

        channel.basic_qos(prefetch_count=1)
        channel.basic_consume(self.dispatch, queue=queue)
        channel.start_consuming()

    def dispatch(self, channel, method, props, body):
        raise BaseException('function dispatch is undefined')
