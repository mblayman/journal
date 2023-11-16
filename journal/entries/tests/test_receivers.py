import datetime

from anymail.inbound import AnymailInboundMessage
from anymail.signals import AnymailInboundEvent

from journal.accounts.models import Account
from journal.accounts.tests.factories import AccountFactory
from journal.entries.models import Entry
from journal.entries.receivers import handle_inbound


def test_persists_entry():
    """The entry is persisted from the event data."""
    when = datetime.date(2023, 11, 15)
    account = AccountFactory()
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        from_email='"JourneyInbox Journal" '
        f"<journal.{account.id}@email.journeyinbox.com>",
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 1
    entry = Entry.objects.first()
    assert entry.body == "Text"
    assert entry.user == account.user
    assert entry.when == when


def test_rejects_bad_date():
    """A malformed date in the subject is ignored."""
    account = AccountFactory()
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's nothing. How are you?",
        text="Text",
        from_email='"JourneyInbox Journal" '
        f"<journal.{account.id}@email.journeyinbox.com>",
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0


def test_rejects_bad_account():
    """A malformed account Sqid is ignored."""
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        from_email='"JourneyInbox Journal" <journal.abcdefg@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0


def test_rejects_non_account_from_email():
    """When the from_email doesn't look like it holds an account ID, ignore it."""
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        from_email='"JourneyInbox Journal" <help@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0


def test_rejects_inactive_account():
    """An inactive account does not persist new entries"""
    account = AccountFactory(status=Account.Status.TRIAL_EXPIRED)
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        from_email='"JourneyInbox Journal" '
        f"<journal.{account.id}@email.journeyinbox.com>",
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0


def test_missing_message():
    """An event with a missing message is ignored."""
    sender = None
    event = AnymailInboundEvent(event_type="inbound")

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0
