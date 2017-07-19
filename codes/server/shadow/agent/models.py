#encoding: utf-8


import json
import datetime

from django.conf import settings
from django.db import models
from django.contrib.postgres.fields import ArrayField
from django.core.exceptions import ObjectDoesNotExist
from django.utils import timezone
from django.core.serializers.json import DjangoJSONEncoder

from base import constants

class ClientManager(models.Manager):

    def get_queryset(self):
        return super(ClientManager, self).get_queryset().filter(status=constants.STATUS_OK).order_by('-update_time')


class Client(models.Model):
    uuid = models.CharField(max_length=128, db_index=True)
    hostname = models.CharField(max_length=64)
    os = models.CharField(max_length=64)
    arch = models.CharField(max_length=64)
    pid = models.IntegerField()
    client_time = models.DateTimeField()
    create_time = models.DateTimeField(auto_now_add=True)
    update_time = models.DateTimeField(auto_now=True)
    status = models.BooleanField(default=constants.STATUS_OK)
    user = models.ForeignKey(settings.AUTH_USER_MODEL, null=True)

    objects = models.Manager()
    ok_objects = ClientManager()

    @classmethod
    def create_or_replace(cls, json_data):
        client = None
        uuid = json_data.get('uuid', '')
        try:
            client = cls.objects.get(uuid=uuid)
        except ObjectDoesNotExist as e:
            client = cls(uuid=uuid)

        client.status = constants.STATUS_OK

        keys = {'client_time' : 'time'}
        attrs = ['hostname', 'os', 'arch', 'pid', 'client_time']
        for attr in attrs:
            setattr(client, attr, json_data.get(keys.get(attr, attr), None))

        client.save()

        for mac, ips in json_data.get('interfaces').items():
            try:
                interface = client.interface_set.get(mac=mac)
                interface.ips = ips
            except ObjectDoesNotExist as e:
                interface = client.interface_set.create(mac=mac, ips=ips)
            interface.status = constants.STATUS_OK
            interface.save()

    @property
    def is_online(self):
        return self.update_time > timezone.now() - datetime.timedelta(seconds=constants.ONLINE_TIMEOUT)

    @property
    def interfaces(self):
        return self.interface_set.filter(status=constants.STATUS_OK)

    @classmethod
    def delete(cls, pk):
        try:
            obj = cls.objects.get(pk=pk)
            obj.status = constants.STATUS_DELETE
            obj.interface_set.update(status=constants.STATUS_DELETE)
            obj.save()
        except ObjectDoesNotExist as e:
            pass


    def __str__(self):
        return json.dumps(self.as_dict())

    def as_dict(self, json=False):
        return {
            'id' : self.id,
            'uuid' : self.uuid,
            'hostname' : self.hostname,
            'os' : self.os,
            'arch' : self.arch,
            'pid' : self.pid,
            'client_time' : timezone.localtime(self.client_time).strftime('%Y-%m-%d %H:%M:%S') if json else self.client_time,
            'create_time' : timezone.localtime(self.create_time).strftime('%Y-%m-%d %H:%M:%S') if json else self.create_time,
            'update_time' : timezone.localtime(self.update_time).strftime('%Y-%m-%d %H:%M:%S') if json else self.update_time,
            'is_online' : self.is_online,
            'status' : self.status,
            'user' : None,
            'interfaces' : [interface.as_dict(json) for interface in self.interfaces],
        }


class Interface(models.Model):
    client = models.ForeignKey(Client)
    mac = models.CharField(max_length=64)
    ips = ArrayField(models.GenericIPAddressField())
    create_time = models.DateTimeField(auto_now_add=True)
    update_time = models.DateTimeField(auto_now=True)
    status = models.BooleanField(default=constants.STATUS_OK)

    def __str__(self):
        return json.dumps(self.as_dict(True))

    def as_dict(self, json=False):
        return {
            'id' : self.id,
            'mac' : self.mac,
            'ips' : self.ips,
            'create_time' : timezone.localtime(self.create_time).strftime('%Y-%m-%d %H:%M:%S') if json else self.create_time,
            'update_time' : timezone.localtime(self.update_time).strftime('%Y-%m-%d %H:%M:%S') if json else self.update_time,
            'status' : self.status,
        }
