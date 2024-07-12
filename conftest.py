import pytest
from django.conf import settings
from djstripe.enums import APIKeyType
from djstripe.models import APIKey

from journal.payments.tests.factories import PriceFactory


@pytest.fixture(autouse=True)
def aaa_db(db):
    pass


@pytest.fixture(autouse=True)
def publishable_key():
    yield APIKey.objects.create(
        secret="pk_test_1234",  # noqa: S106 This is test data.
        livemode=False,
        type=APIKeyType.publishable,
    )


@pytest.fixture(autouse=True)
def price():
    yield PriceFactory(lookup_key=settings.PRICE_LOOKUP_KEY)
