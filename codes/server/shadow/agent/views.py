#encoding: utf-8
from django.shortcuts import render
from django.http import JsonResponse

from .models import Client
from base import constants

from base.decorators import login_required

@login_required
def lists(request):
    if request.is_ajax():
        result = [client.as_dict(True) for client in Client.ok_objects.all()]
        return JsonResponse({'status' : constants.HTTP_STATUS_OK, 'result' : result})
    else:
        return render(request, 'agent/list.html')

@login_required
def delete(request):
    id = request.POST.get('id')
    Client.delete(id)
    return JsonResponse({'status' : constants.HTTP_STATUS_OK})
