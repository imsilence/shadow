#encoding: utf-8

from django.conf.urls import url

from django.views.generic import TemplateView

app_name = 'dashboard'

urlpatterns = [
    url(r'^$', TemplateView.as_view(template_name='dashboard/index.html'), name='index'),
]
