from django.conf import settings
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
    def price(self) -> Price:
        return Price.objects.get(
            lookup_key=settings.PRICE_LOOKUP_KEY,
            livemode=settings.STRIPE_LIVE_MODE,
        )
