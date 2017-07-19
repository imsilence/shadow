#encoding: utf-8

from functools import wraps

from django.http import JsonResponse
from django.shortcuts import render
from django.conf import settings

from . import constants

def login_required(func):

    @wraps(func)
    def wrapper(request, *args, **kwargs):
        if not request.user.is_authenticated():
            if request.is_ajax():
                return JsonResponse({'status' : constants.HTTP_STATUS_UNAUTHENTICATION})
            else:
                return render(settings.LOGIN_URL)
        return func(request, *args, **kwargs)

    return wrapper
