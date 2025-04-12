from django.utils import timezone
from huey import crontab
from huey.contrib.djhuey import db_periodic_task

from journal.accounts.models import Account

from .transmitter import send_prompt


@db_periodic_task(crontab(minute="0", hour="13"))
def send_mail():
    return
    print("Sending prompts to active accounts")
    accounts = Account.objects.promptable().select_related("user")
    today = timezone.localdate()
    for account in accounts:
        send_prompt(account, today)
