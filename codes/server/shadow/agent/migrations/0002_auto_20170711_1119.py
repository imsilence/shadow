# -*- coding: utf-8 -*-
# Generated by Django 1.11.3 on 2017-07-11 03:19
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('agent', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='client',
            name='uuid',
            field=models.CharField(db_index=True, max_length=128),
        ),
    ]