import time_machine
from django.utils import timezone

from journal.accounts.tests.factories import UserFactory
from journal.entries.jobs.send_mail import Job as SendMailJob
from journal.entries.models import Prompt
from journal.entries.tests.factories import EntryFactory, PromptFactory


class TestSendMailJob:
    @time_machine.travel("2023-07-19")
    def test_send_email(self, mailoutbox):
        """An active account receives an email prompt."""
        user = UserFactory()
        body = "This is the entry.\n\nIt has newlines."
        entry = EntryFactory(user=user, body=body)
        job = SendMailJob()

        job.execute()

        assert len(mailoutbox) == 1
        mail = mailoutbox[0]
        assert mail.from_email == (
            '"JourneyInbox Journal" '
            f"<journal.{user.account.id}@email.journeyinbox.com>"
        )
        assert mail.to == [user.email]
        assert mail.subject == "It's Wednesday, Jul. 19, 2023. How are you?"
        assert entry.body in mail.body  # Test the text email.
        html_message = mail.alternatives[0][0]
        assert "<p>This is the entry.</p>\n\n<p>It has newlines.</p>" in html_message
        assert Prompt.objects.filter(user=user).count() == 1

    def test_no_available_entries(self, mailoutbox):
        """The message indicates that a previous entry will appear once it exists."""
        UserFactory()
        job = SendMailJob()

        job.execute()

        assert len(mailoutbox) == 1
        mail = mailoutbox[0]
        assert "You have no entries yet!" in mail.body  # Test the text email.
        html_message = mail.alternatives[0][0]
        assert "<p>You have no entries yet!" in html_message

    @time_machine.travel("2023-07-19")
    def test_send_email_idempotent(self, mailoutbox):
        """A user will not receive a prompt twice."""
        user = UserFactory()
        body = "This is the entry.\n\nIt has newlines."
        EntryFactory(user=user, body=body)
        PromptFactory(user=user, when=timezone.localdate())
        job = SendMailJob()

        job.execute()

        assert len(mailoutbox) == 0
        assert Prompt.objects.filter(user=user).count() == 1
