from django.apps import AppConfig
from django.contrib.auth.signals import user_logged_in

from .receivers import verify_account


class AccountsConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "journal.accounts"

    def ready(self):
        user_logged_in.connect(verify_account)
