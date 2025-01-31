from django.utils import timezone

from journal.entries.transmitter import send_prompt


def handle_user_logged_in(sender, request, user, **kwargs):
    """Mark an account as verified."""
    # Must be first log in if not verified. Send a prompt immediately.
    if not user.account.verified:
        today = timezone.localdate()
        send_prompt(user.account, today)

    user.account.verified = True
    user.account.save()
