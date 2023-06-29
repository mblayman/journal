from journal.accounts.tests.factories import UserFactory
from journal.entries.jobs.send_mail import Job as SendMailJob


class TestSendMailJob:
    def test_send_email(self, mailoutbox):
        """An active account receives an email prompt."""
        user = UserFactory()
        job = SendMailJob()

        job.execute()

        assert len(mailoutbox) == 1
        mail = mailoutbox[0]
        assert mail.to == [user.email]
        # TODO: assert subject
        # TODO: assert message
        # TODO: assert html_message
        # TODO: assert from_email
