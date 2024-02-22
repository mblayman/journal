import stripe
from django.conf import settings
from django.contrib.auth.models import User
from django.contrib.sites.models import Site
from django.urls import reverse
from djstripe.enums import APIKeyType
from djstripe.models import APIKey, Price


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

        checkout_session = stripe.checkout.Session.create(
            api_key=self.secret_key, **session_parameters
        )
        return checkout_session["id"]
