import datetime

import stripe
from django.conf import settings
from django.contrib.auth.models import User
from django.contrib.sites.models import Site
from django.urls import reverse
from django.utils import timezone
from djstripe.enums import APIKeyType
from djstripe.models import APIKey, Customer, Price

from journal.accounts.constants import TRIAL_DAYS


class PaymentsGateway:
    """A gateway for interacting with the payments vendor."""

    @property
    def publishable_key(self) -> str:
        return APIKey.objects.get(
            type=APIKeyType.publishable,
            livemode=settings.STRIPE_LIVE_MODE,
        ).secret

    @property
    def secret_key(self) -> str:
        return APIKey.objects.get(
            type=APIKeyType.secret,
            livemode=settings.STRIPE_LIVE_MODE,
        ).secret

    @property
    def price(self) -> Price:
        return Price.objects.get(
            lookup_key=settings.PRICE_LOOKUP_KEY,
            livemode=settings.STRIPE_LIVE_MODE,
        )

    def create_checkout_session(self, price_id: str, user: User) -> str:
        """Create a Stripe checkout session."""
        site = Site.objects.get_current()
        success = reverse("success")

        session_parameters = {
            "customer_email": user.email,
            "success_url": f"https://{site}{success}",
            "cancel_url": f"https://{site}/",
            # TODO: Should we accept other payment methods? Issue #73
            "payment_method_types": ["card"],
            "mode": "subscription",
            "line_items": [{"price": price_id, "quantity": 1}],
            "client_reference_id": str(user.id),
        }

        if self._is_trial_eligible(user):
            # Be generous and include an extra two days.
            # This also makes Stripe display nicer so if someone signs up
            # on the same day with a credit card, it will show the full number
            # of days on the trial.
            trial_end = self._trial_end(user) + datetime.timedelta(days=2)
            session_parameters["subscription_data"] = {
                # Stripe expects a Unix timestamp in whole seconds.
                "trial_end": int(trial_end.timestamp())
            }

        checkout_session = stripe.checkout.Session.create(
            api_key=self.secret_key, **session_parameters
        )
        return checkout_session["id"]

    def _is_trial_eligible(self, user: User) -> bool:
        """Check if the account is eligible for Stripe's trial data.

        The trial must end at least 48 hours in the future. See:
        https://stripe.com/docs/api/checkout/sessions/create#create_checkout_session-subscription_data-trial_end
        """
        cutoff = timezone.now() + datetime.timedelta(days=2)
        return self._trial_end(user) > cutoff

    def _trial_end(self, user: User) -> datetime.datetime:
        return user.date_joined + datetime.timedelta(days=TRIAL_DAYS)

    def create_billing_portal_session(self, user: User) -> str:
        """Create a billing portal session at Stripe.

        This method assumes that there is an existing Stripe customer
        for the account.
        """
        site = Site.objects.get_current()
        return_url = reverse("settings")
        customer = Customer.objects.get(email=user.email)
        session = stripe.billing_portal.Session.create(
            api_key=self.secret_key,
            customer=customer.id,
            return_url=f"https://{site}{return_url}",
        )
        return session.url
