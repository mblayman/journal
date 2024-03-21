from anymail.message import AnymailMessage
from django.template.loader import render_to_string
from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account

from ..models import Entry, Prompt


class Job(DailyJob):
    help = "Sent mail to active accounts"

    def execute(self):
        print("Sending prompts to active accounts")
        accounts = Account.objects.promptable().select_related("user")
        today = timezone.localdate()
        for account in accounts:
            if Prompt.objects.exists_for(account.user, today):
                print(f"Prompt already exists for {account.user.id} on {today}.")
                continue

            context = {
                "entry": Entry.objects.get_random_for(account.user),
                "today": today,
            }
            text_message = render_to_string("entries/email/prompt.txt", context)
            html_message = render_to_string("entries/email/prompt.html", context)
            from_email = (
                '"JourneyInbox Journal" '
                f"<journal.{account.id}@email.journeyinbox.com>"
            )
            message = AnymailMessage(
                subject=(
                    f"It's {today:%A}, {today:%b}. {today:%-d}, {today:%Y}. "
                    "How are you?"
                ),
                body=text_message,
                from_email=from_email,
                to=[account.user.email],
            )
            message.attach_alternative(html_message, "text/html")
            message.metadata = {"entry_date": str(today)}
            message.send()
            Prompt.objects.create(
                user=account.user,
                when=today,
                # message_id is not nullable, but during the tests, the in-memory
                # backend does not set a value. Accept empty string to avoid
                # nasty mocking hacks.
                message_id=message.anymail_status.message_id or "",
            )
            print(f"Prompt sent for {account.user.id}.")
