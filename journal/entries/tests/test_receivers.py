import datetime

from anymail.inbound import AnymailInboundMessage
from anymail.signals import AnymailInboundEvent
from django.http import HttpRequest

from journal.accounts.models import Account
from journal.accounts.tests.factories import AccountFactory, UserFactory
from journal.entries.models import Entry, Prompt
from journal.entries.receivers import handle_email_confirmed, handle_inbound


def test_email_confirmed_sends_prompt():
    """The email_confirmed signal handler sends a prompt."""
    sender = None
    request = HttpRequest()
    user = UserFactory()

    handle_email_confirmed(sender, request, user.email)

    assert Prompt.objects.filter(user=user).exists()


def test_persists_entry():
    """The entry is persisted from the event data."""
    when = datetime.date(2023, 11, 15)
    account = AccountFactory()
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 1
    entry = Entry.objects.first()
    assert entry.body == "Text"
    assert entry.user == account.user
    assert entry.when == when


def test_persists_one_entry():
    """The entry is persisted one time and updated."""
    account = AccountFactory()
    sender = None

    # Initial creation
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)
    handle_inbound(sender, event, "SendGrid")

    # Updated entry
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text Updated",
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)
    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 1
    entry = Entry.objects.first()
    assert entry.body == "Text Updated"


def test_rejects_bad_date():
    """A malformed date in the subject is ignored."""
    account = AccountFactory()
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's nothing. How are you?",
        text="Text",
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
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
        to='"JourneyInbox Journal" <journal.abcdefg@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 0


def test_rejects_non_account_to_email():
    """When the to_email doesn't look like it holds an account ID, ignore it."""
    sender = None
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text="Text",
        to='"JourneyInbox Journal" <help@email.journeyinbox.com>',
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
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
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


def test_parse_message_text():
    """The message text excludes the prompt information."""
    account = AccountFactory()
    sender = None
    text = "Did we make it?\r\n\r\nThis is a second line.\r\n\r\nOn Wed, Nov 15, 2023, 11:11 PM JourneyInbox Journal <\r\njournal.abcdef@email.journeyinbox.com> wrote:\r\n\r\n> How are you? Reply to this prompt to update your journal.\r\n>\r\n> You have no entries yet! As soon as you do, a random previous entry will\r\n> appear in your prompt.\r\n>\r\n>\r\n"  # noqa
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text=text,
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 1
    entry = Entry.objects.first()
    assert entry.body == "Did we make it?\n\nThis is a second line."


def test_parse_message_text_missing_marker():
    """When the prompt boundary marker is missing, all text is included."""
    account = AccountFactory()
    sender = None
    text = "Did we make it?\r\n\r\nThis is a second line.\r\n\r\nOn Wed, Nov 15, 2023, 11:11 PM JourneyInbox\r\nJournal <\r\njournal.abcdef@email.journeyinbox.com> wrote:\r\n\r\n> How are you? Reply to this prompt to update your journal.\r\n>\r\n> You have no entries yet! As soon as you do, a random previous entry will\r\n> appear in your prompt.\r\n>\r\n>\r\n"  # noqa
    message = AnymailInboundMessage.construct(
        subject="RE: It's Wednesday, Nov. 15, 2023. How are you?",
        text=text,
        to=f'"JourneyInbox Journal" <journal.{account.id}@email.journeyinbox.com>',
    )
    event = AnymailInboundEvent(event_type="inbound", message=message)

    handle_inbound(sender, event, "SendGrid")

    assert Entry.objects.count() == 1
    entry = Entry.objects.first()
    assert (
        entry.body
        == """Did we make it?\n\nThis is a second line.\n\nOn Wed, Nov 15, 2023, 11:11 PM JourneyInbox\nJournal <\njournal.abcdef@email.journeyinbox.com> wrote:\n\n> How are you? Reply to this prompt to update your journal.\n>\n> You have no entries yet! As soon as you do, a random previous entry will\n> appear in your prompt.\n>\n>"""  # noqa
    )
