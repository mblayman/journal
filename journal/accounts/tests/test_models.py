from journal.accounts.tests.factories import UserFactory


class TestUser:
    def test_factory(self):
        """The factory produces a valid instance."""
        user = UserFactory()

        assert user is not None
