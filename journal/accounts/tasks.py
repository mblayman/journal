import datetime

from django.utils import timezone
from huey import crontab
from huey.contrib.djhuey import db_periodic_task

from . import constants
from .models import Account


@db_periodic_task(crontab(minute="0", hour="0"))
def expire_trials():
    """Expire any accounts that are TRIALING beyond the trial days limit"""
    print("Check for expired trials.")
    # Give an extra couple of days to be gracious and avoid customer complaints.
    cutoff_days = constants.TRIAL_DAYS + 2
    trial_cutoff = timezone.now() - datetime.timedelta(days=cutoff_days)
    expired_trials = Account.objects.filter(
        status=Account.Status.TRIALING, user__date_joined__lt=trial_cutoff
    )
    count = expired_trials.update(status=Account.Status.TRIAL_EXPIRED)
    print(f"Expired {count} trial(s)")
