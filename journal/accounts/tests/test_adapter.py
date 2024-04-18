from journal.accounts.adapter import AccountAdapter
from journal.accounts.constants import MAX_TRIALING_USERS
from journal.accounts.models import Account
from journal.accounts.tests.factories import AccountFactory


def test_signup_open():
    """When trials are below limit, signup is allowed."""
    adapter = AccountAdapter()

    assert adapter.is_open_for_signup(request=None)


def test_signup_closed():
    """When trials are below limit, signup is allowed."""
    for _ in range(MAX_TRIALING_USERS):
        AccountFactory(status=Account.Status.TRIALING)

    adapter = AccountAdapter()

    assert not adapter.is_open_for_signup(request=None)
