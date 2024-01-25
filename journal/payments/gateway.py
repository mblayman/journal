from django.conf import settings
from djstripe.enums import APIKeyType
from djstripe.models import APIKey


class PaymentsGateway:
    """A gateway for interacting with the payments vendor."""

    @property
    def publishable_key(self) -> str:
        return APIKey.objects.get(
            type=APIKeyType.publishable,
            livemode=settings.STRIPE_LIVE_MODE,
        ).secret
