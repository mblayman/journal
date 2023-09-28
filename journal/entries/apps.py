from anymail.signals import inbound
from django.apps import AppConfig

from .receivers import handle_inbound


class EntriesConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "journal.entries"

    def ready(self):
        inbound.connect(handle_inbound)
