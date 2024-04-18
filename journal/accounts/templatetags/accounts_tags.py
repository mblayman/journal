from django import template

from journal.accounts.constants import MAX_TRIALING_USERS
from journal.accounts.models import Account

register = template.Library()


@register.simple_tag
def available_trial_slots():
    """For the signup page, advertise the number of open trial slots."""
    # I don't feel like testing this as I view it as temporary.
    return MAX_TRIALING_USERS - Account.objects.trialing().count()  # pragma: no cover
