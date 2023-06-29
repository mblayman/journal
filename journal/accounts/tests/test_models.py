from journal.accounts.models import Account
from journal.accounts.tests.factories import AccountFactory, UserFactory


class TestAccount:
    def test_factory(self):
        """The factory produces a valid instance."""
        account = AccountFactory()

        assert account is not None
        assert account.user is not None
        assert account.status == account.Status.TRIALING

    def test_active(self):
        """The active manager method returns active accounts."""
        trialing = AccountFactory(status=Account.Status.TRIALING)
        active = AccountFactory(status=Account.Status.ACTIVE)
        exempt = AccountFactory(status=Account.Status.EXEMPT)
        AccountFactory(status=Account.Status.CANCELED)
        AccountFactory(status=Account.Status.TRIAL_EXPIRED)

        accounts = Account.objects.active()

        assert set(accounts) == {trialing, active, exempt}


class TestUser:
    def test_factory(self):
        """The factory produces a valid instance."""
        user = UserFactory()

        assert user is not None
        assert user.account is not None
