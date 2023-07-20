import time_machine

from journal.accounts.tests.factories import UserFactory
from journal.entries.jobs.send_mail import Job as SendMailJob


class TestSendMailJob:
    @time_machine.travel("2023-07-19")
    def test_send_email(self, mailoutbox):
        """An active account receives an email prompt."""
        user = UserFactory()
        job = SendMailJob()

        job.execute()

        assert len(mailoutbox) == 1
        mail = mailoutbox[0]
        assert mail.to == [user.email]
        assert mail.subject == "It's Wednesday, Jul. 19, how are you?"
        # TODO: assert message
        # TODO: assert html_message
        # TODO: assert from_email
