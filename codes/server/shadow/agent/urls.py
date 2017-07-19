#encoding: utf-8

from django.conf.urls import url

app_name = 'agent'
from . import views

urlpatterns = [
    url(r'^list/', views.lists, name='list'),
    url(r'^delete/', views.delete, name='delete'),
]
