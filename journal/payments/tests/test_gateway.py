import datetime
from unittest import mock

from django.conf import settings
from django.utils import timezone
from djstripe.enums import APIKeyType
from djstripe.models import APIKey

from journal.accounts.constants import TRIAL_DAYS
from journal.accounts.tests.factories import UserFactory
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


def test_price():
    """The gateway can return a valid price ID."""
    PriceFactory(lookup_key=settings.PRICE_LOOKUP_KEY)
    gateway = PaymentsGateway()

    price = gateway.price

    assert price.lookup_key == settings.PRICE_LOOKUP_KEY


@mock.patch("journal.payments.gateway.stripe")
def test_create_checkout_session(mock_stripe):
    """The gateway gets a valid session ID."""
    APIKey.objects.create(
        type=APIKeyType.secret,
        livemode=settings.STRIPE_LIVE_MODE,
        secret="pk_some_value",  # noqa
    )
    mock_stripe.checkout.Session.create.return_value = {"id": "fake_session_id"}
    user = UserFactory()
    gateway = PaymentsGateway()

    session_id = gateway.create_checkout_session("fake_price_id", user)

    kwargs = mock_stripe.checkout.Session.create.call_args.kwargs
    assert kwargs["customer_email"] == user.email
    assert kwargs["client_reference_id"] == str(user.id)
    assert kwargs["cancel_url"] == "https://example.com/"
    assert kwargs["success_url"] == "https://example.com/success/"
    assert session_id == "fake_session_id"
    assert "subscription_data" in kwargs


@mock.patch("journal.payments.gateway.stripe")
def test_no_trial_in_stripe_limits(mock_stripe):
    """A trial is not added to the session within Stripe's limit.

    Stripe does not permit a trial unless there are at least 48 hours from now.
    """
    APIKey.objects.create(
        type=APIKeyType.secret,
        livemode=settings.STRIPE_LIVE_MODE,
        secret="pk_some_value",  # noqa
    )
    mock_stripe.checkout.Session.create.return_value = {"id": "fake_session_id"}
    gateway = PaymentsGateway()
    too_close_date = timezone.now() - datetime.timedelta(
        days=TRIAL_DAYS - 1  # Set the user 24 hours from trial ending.
    )
    user = UserFactory(date_joined=too_close_date)
    gateway = PaymentsGateway()

    gateway.create_checkout_session("price_fake_id", user)

    kwargs = mock_stripe.checkout.Session.create.call_args.kwargs
    assert "subscription_data" not in kwargs
