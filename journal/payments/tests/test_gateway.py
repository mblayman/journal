from django.conf import settings
from djstripe.enums import APIKeyType
from djstripe.models import APIKey

from journal.payments.gateway import PaymentsGateway
from journal.payments.tests.factories import PriceFactory


def test_publishable_key():
    """The gateway can return a valid publishable secret."""
    APIKey.objects.create(
        type=APIKeyType.publishable,
        livemode=settings.STRIPE_LIVE_MODE,
        secret="pk_some_value",  # noqa
    )
    gateway = PaymentsGateway()

    publishable_key = gateway.publishable_key

    assert publishable_key.startswith("pk_")


def test_price_id():
    """The gateway can return a valid price ID."""
    PriceFactory(lookup_key=settings.PRICE_LOOKUP_KEY)
    gateway = PaymentsGateway()

    price = gateway.price

    assert price.lookup_key == settings.PRICE_LOOKUP_KEY
