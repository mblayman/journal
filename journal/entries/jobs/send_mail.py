from django.core import mail
from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        accounts = Account.objects.active().select_related("user")
        today = timezone.localdate()
        for account in accounts:
            mail.send_mail(
                subject=f"It's {today:%A}, {today:%b}. {today:%-d}, how are you?",
                message="Replace this message",
                html_message="Replace this HTML message",
                from_email="who is this from email",
                recipient_list=[account.user.email],
            )
