from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


class TestIndex:
    def test_authenticated(self, client):
        """An authenticated user gets a valid response."""
        client.force_login(UserFactory())

        response = client.get(reverse("create_checkout_session"))

        assert response.status_code == 200
