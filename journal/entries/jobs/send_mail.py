from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account

from ..transmitter import send_prompt


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        print("Sending prompts to active accounts")
        accounts = Account.objects.promptable().select_related("user")
        today = timezone.localdate()
        for account in accounts:
            send_prompt(account, today)
