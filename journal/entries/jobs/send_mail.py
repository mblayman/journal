from django.core import mail
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        accounts = Account.objects.active().select_related("user")
        for account in accounts:
            mail.send_mail(
                subject="Replace this subject",
                message="Replace this message",
                html_message="Replace this HTML message",
                from_email="who is this from email",
                recipient_list=[account.user.email],
            )
