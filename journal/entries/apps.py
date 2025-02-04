from anymail.signals import inbound
from django.apps import AppConfig


class EntriesConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "journal.entries"

    def ready(self):
        from .receivers import handle_inbound

        inbound.connect(handle_inbound)
