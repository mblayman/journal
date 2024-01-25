import pytest
from django.urls import reverse
from djstripe.enums import APIKeyType
from djstripe.models import APIKey

from journal.accounts.tests.factories import UserFactory


@pytest.fixture
def publishable_key():
    yield APIKey.objects.create(
        secret="pk_test_1234",  # noqa: S106 This is test data.
        livemode=False,
        type=APIKeyType.publishable,
    )


@pytest.mark.usefixtures("publishable_key")
class TestIndex:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("index"))

        assert response.status_code == 200

    def test_authenticated(self, client):
        """An authenticated user gets a valid response."""
        client.force_login(UserFactory())

        response = client.get(reverse("index"))

        assert response.status_code == 200


class TestTerms:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("terms"))

        assert response.status_code == 200


class TestPrivacyPolicy:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("privacy"))

        assert response.status_code == 200
