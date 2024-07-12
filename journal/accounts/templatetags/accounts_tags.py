from django import template

from journal.accounts.constants import MAX_TRIALING_USERS
from journal.accounts.models import Account
from journal.payments.gateway import PaymentsGateway

register = template.Library()


@register.simple_tag
def available_trial_slots():
    """For the signup page, advertise the number of open trial slots."""
    # I don't feel like testing this as I view it as temporary.
    return MAX_TRIALING_USERS - Account.objects.trialing().count()  # pragma: no cover


@register.inclusion_tag("accounts/trial_banner.html", takes_context=True)
def trial_banner(context):
    user = context["user"]
    if not user.is_authenticated or user.account.status not in (
        user.account.Status.TRIALING,
        user.account.Status.TRIAL_EXPIRED,
    ):
        return {}

    payments_gateway = PaymentsGateway()
    return {
        "payments_publishable_key": payments_gateway.publishable_key,
        "price": payments_gateway.price,
        "user": user,
    }
