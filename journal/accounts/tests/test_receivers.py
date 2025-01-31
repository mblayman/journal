from journal.accounts.models import Account
from journal.accounts.receivers import handle_user_logged_in
from journal.accounts.tests.factories import UserFactory
from journal.entries.models import Prompt


class TestHandleUserLoggedIn:
    def test_verifies(self):
        """An account is verified."""
        sender = None
        request = None
        user = UserFactory()

        handle_user_logged_in(sender, request, user)

        assert Account.objects.get(user=user).verified

    def test_user_logged_in_sends_prompt(self):
        """The user_logged_in signal handler sends a prompt."""
        sender = None
        request = None
        user = UserFactory()

        handle_user_logged_in(sender, request, user)

        assert Prompt.objects.filter(user=user).exists()
