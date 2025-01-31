from django.apps import AppConfig
from django.contrib.auth.signals import user_logged_in


class AccountsConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "journal.accounts"

    def ready(self):
        from .receivers import handle_user_logged_in

        user_logged_in.connect(handle_user_logged_in)
