import datetime

from anymail.message import AnymailMessage
from django.template.loader import render_to_string
from django.utils import timezone
from huey import crontab
from huey.contrib.djhuey import db_periodic_task, db_task
from sesame.utils import get_query_string

from journal.core.site import full_url_reverse

from . import constants
from .models import Account, User


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


def _generate_magic_link(user_id):
    """Generate magic link and send email."""
    user = User.objects.get(id=user_id)
    login_url = full_url_reverse("sesame-login") + get_query_string(user)
    context = {"login_url": login_url}
    text_message = render_to_string("accounts/email/login.txt", context)
    html_message = render_to_string("accounts/email/login.html", context)

    message = AnymailMessage(
        subject="Signin to JourneyInbox",
        body=text_message,
        from_email='"JourneyInbox" <noreply@email.journeyinbox.com>',
        to=[user.email],
    )
    message.attach_alternative(html_message, "text/html")
    message.send()


generate_magic_link = db_task()(_generate_magic_link)
