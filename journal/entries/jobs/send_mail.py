from anymail.message import AnymailMessage
from django.conf import settings
from django.template.loader import render_to_string
from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account

from ..models import Entry


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        print("Sending prompts to active accounts")
        accounts = Account.objects.active().select_related("user")
        today = timezone.localdate()
        for account in accounts:
            context = {
                "entry": Entry.objects.get_random_for(account.user),
                "today": today,
            }
            text_message = render_to_string("entries/email/prompt.txt", context)
            html_message = render_to_string("entries/email/prompt.html", context)
            message = AnymailMessage(
                subject=f"It's {today:%A}, {today:%b}. {today:%-d}. How are you?",
                body=text_message,
                from_email=settings.EMAIL_SENDGRID_REPLY_TO,
                to=[account.user.email],
            )
            message.attach_alternative(html_message, "text/html")
            message.metadata = {
                "metadata_key": "metadata_value",
                "entry_date": str(today),
            }
            print(message.metadata)
            message.send()
            print("Message ID was:")
            print(message.anymail_status.message_id)
