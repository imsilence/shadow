#encoding: utf-8

from django.contrib.auth.forms import AuthenticationForm
from django.utils.translation import gettext, gettext_lazy as _

class LoginForm(AuthenticationForm):

    error_messages = {
        'invalid_login': _(
            "Please enter a correct username and password."
        ),
        'inactive': _("This account is inactive."),
    }
