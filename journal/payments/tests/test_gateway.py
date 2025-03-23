from django.conf import settings

from journal.payments.gateway import PaymentsGateway


def test_publishable_key():
    """The gateway can return a valid publishable secret.

    The publishable key is global conftest data.
    """
    gateway = PaymentsGateway()

    publishable_key = gateway.publishable_key

    assert publishable_key.startswith("pk_")


def test_price():
    """The gateway can return a valid price ID.

    The price is global conftest data.
    """
    gateway = PaymentsGateway()

    price = gateway.price

    assert price.lookup_key == settings.PRICE_LOOKUP_KEY
