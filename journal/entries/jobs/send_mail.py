from django.conf import settings
from django.core import mail
from django.template.loader import render_to_string
from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account

from ..models import Entry


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        accounts = Account.objects.active().select_related("user")
        today = timezone.localdate()
        for account in accounts:
            context = {
                "entry": Entry.objects.get_random_for(account.user),
                "today": today,
            }
            text_message = render_to_string("entries/email/prompt.txt", context)
            html_message = render_to_string("entries/email/prompt.html", context)
            mail.send_mail(
                subject=f"It's {today:%A}, {today:%b}. {today:%-d}, how are you?",
                message=text_message,
                html_message=html_message,
                from_email=settings.EMAIL_SENDGRID_REPLY_TO,
                recipient_list=[account.user.email],
            )
