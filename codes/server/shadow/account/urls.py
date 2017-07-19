#encoding: utf-8

from django.conf.urls import url, include
from django.views.generic import TemplateView
from django.contrib.auth import views as auth_views

from .forms import LoginForm

app_name = 'account'

urlpatterns = [
    url(r'^login/$', auth_views.LoginView.as_view(template_name='account/login.html', authentication_form=LoginForm), name='login'),
    url(r'^logout/$', auth_views.LogoutView.as_view(), name='logout'),
    url(r'^list/$', TemplateView.as_view(template_name='account/list.html'), name='list')
]
#
