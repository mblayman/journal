from allauth.account.signals import email_confirmed
from anymail.signals import inbound
from django.apps import AppConfig


class EntriesConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "journal.entries"

    def ready(self):
        from .receivers import handle_email_confirmed, handle_inbound

        email_confirmed.connect(handle_email_confirmed)
        inbound.connect(handle_inbound)
