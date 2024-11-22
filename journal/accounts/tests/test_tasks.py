import datetime

from django.utils import timezone

from journal.accounts import constants
from journal.accounts.models import Account
from journal.accounts.tasks import expire_trials
from journal.accounts.tests.factories import AccountFactory


class TestExpireTrials:
    def test_expires_trials(self):
        """Old trials are marked as expired."""
        trialing = Account.Status.TRIALING
        account = AccountFactory(
            status=trialing,
            user__date_joined=timezone.now()
            - datetime.timedelta(days=constants.TRIAL_DAYS + 5),
        )

        expire_trials()

        account.refresh_from_db()
        assert account.status == Account.Status.TRIAL_EXPIRED

    def test_keep_active_trials(self):
        """Recent trials are not expired."""
        trialing = Account.Status.TRIALING
        account = AccountFactory(status=trialing, user__date_joined=timezone.now())

        expire_trials()

        account.refresh_from_db()
        assert account.status == Account.Status.TRIALING

    def test_other_statuses_not_expired(self):
        """Only TRIALING is marked as expired."""
        active = Account.Status.ACTIVE
        account = AccountFactory(
            status=active,
            user__date_joined=timezone.now()
            - datetime.timedelta(days=constants.TRIAL_DAYS + 5),
        )

        expire_trials()

        account.refresh_from_db()
        assert account.status == Account.Status.ACTIVE
