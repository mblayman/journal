from journal.accounts.models import Account
from journal.accounts.receivers import verify_account
from journal.accounts.tests.factories import UserFactory


class TestVerifyAccount:
    def test_verifies(self):
        """An account is verified."""
        user = UserFactory()

        verify_account(None, None, user)

        assert Account.objects.get(user=user).verified
