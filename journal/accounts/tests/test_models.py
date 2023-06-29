from journal.accounts.tests.factories import AccountFactory, UserFactory


class TestAccount:
    def test_factory(self):
        account = AccountFactory()

        assert account is not None
        assert account.user is not None
        assert account.status == account.Status.TRIALING


class TestUser:
    def test_factory(self):
        """The factory produces a valid instance."""
        user = UserFactory()

        assert user is not None
        assert user.account is not None
