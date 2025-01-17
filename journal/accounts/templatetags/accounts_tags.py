from django import template

from journal.payments.gateway import PaymentsGateway

register = template.Library()


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
