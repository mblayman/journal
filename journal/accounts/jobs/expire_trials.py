import datetime

from django.utils import timezone
from django_extensions.management.jobs import DailyJob

from journal.accounts.models import Account

from .. import constants


class Job(DailyJob):
    help = "Expire trialing accounts that are past the trial window"

    def execute(self):
        print("Check for expired trials.")
        # Give an extra couple of days to be gracious and avoid customer complaints.
        cutoff_days = constants.TRIAL_DAYS + 2
        trial_cutoff = timezone.now() - datetime.timedelta(days=cutoff_days)
        expired_trials = Account.objects.filter(
            status=Account.Status.TRIALING, user__date_joined__lt=trial_cutoff
        )
        count = expired_trials.update(status=Account.Status.TRIAL_EXPIRED)
        print(f"Expired {count} trial(s)")
