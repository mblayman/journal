from typing import Any

from anymail.signals import AnymailInboundEvent
from dateutil.parser import ParserError, parse
from django.http import HttpRequest
from django.utils import timezone

from journal.accounts.models import Account
from journal.entries.models import Entry
from journal.entries.transmitter import send_prompt


# FIXME: orphaned. Need to send on first sign in
def handle_email_confirmed(
    sender: Any, request: HttpRequest, email_address: str, **kwargs: Any
) -> None:
    """Send a prompt to a user as soon as they confirm their email."""
    account = Account.objects.from_email(email_address)
    today = timezone.localdate()
    send_prompt(account, today)


def handle_inbound(
    sender: Any, event: AnymailInboundEvent, esp_name: str, **kwargs: Any
) -> None:
    message = event.message
    if message is None:
        print("No message")
        return None

    body = parse_body(message.text)

    try:
        entry_datetime = parse(message.subject, fuzzy_with_tokens=True)[0]  # type: ignore
    except ParserError:
        print("Bad date parse")
        return None

    username = message.to[0].username
    if "." not in username:
        print("Username missing period.")
        return None

    account_id = username.split(".")[1]
    try:
        account = Account.objects.active().get(id=account_id)
    except Account.DoesNotExist:
        print(f"No active account ID: {account_id}")
        return None

    Entry.objects.update_or_create(
        when=entry_datetime.date(),
        user=account.user,
        defaults={"body": body},
    )


def parse_body(message_text: str) -> str:
    """Parse the body out of the message text and strip off the prompt."""
    lines = message_text.splitlines()
    marker_line_index = 0
    marker_found = False
    for index, line in enumerate(lines):
        if "JourneyInbox Journal" not in line:
            continue
        marker_line_index = index
        marker_found = True
        break

    if marker_found:
        lines = lines[:marker_line_index]

    body = "\n".join(lines)
    return body.rstrip("\n")
