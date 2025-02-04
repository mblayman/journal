import unicodedata
from datetime import date

from anymail.message import AnymailMessage
from django.template.loader import render_to_string
from django.utils.timesince import timesince

from journal.accounts.models import Account

from .models import Entry, Prompt


def send_prompt(account: Account, today: date) -> None:
    """Send a prompt to the account for the date provided.

    Returns True when sent, otherwise False.
    """
    if Prompt.objects.exists_for(account.user, today):
        print(f"Prompt already exists for {account.user.id} on {today}.")
        return None

    entry = Entry.objects.get_random_for(account.user)
    message = _send_message(account, entry, today)

    Prompt.objects.create(
        user=account.user,
        when=today,
        # message_id is not nullable, but during the tests, the in-memory
        # backend does not set a value. Accept empty string to avoid
        # nasty mocking hacks.
        message_id=message.anymail_status.message_id or "",
    )
    print(f"Prompt sent for {account.user.id}.")


def _send_message(account: Account, entry: Entry | None, today: date):
    """Send an individual message to an account."""
    context = {"entry": entry, "today": today}
    if entry:
        # We need to normalize timesince because it uses non-breakable space
        # (i.e., \xa0) and this is a character from Latin1 (ISO-8859-1).
        # Gmail expects all unicode and will add a "View entire message" link
        # when there are characters that it doesn't like.
        # By normalizing, this replaces the non-breakable space with a regular
        # space.
        delta = timesince(entry.when, today)
        context["delta"] = unicodedata.normalize("NFKD", delta)

    text_message = render_to_string("entries/email/prompt.txt", context)
    html_message = render_to_string("entries/email/prompt.html", context)
    from_email = f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>'
    message = AnymailMessage(
        subject=(f"It's {today:%A}, {today:%b}. {today:%-d}, {today:%Y}. How are you?"),
        body=text_message,
        from_email=from_email,
        to=[account.user.email],
    )
    message.attach_alternative(html_message, "text/html")
    message.metadata = {"entry_date": str(today)}
    message.send()
    return message
