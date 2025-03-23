from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


class TestAccountSettings:
    def test_authenticated(self, client):
        """An authenticated user gets a valid response."""
        client.force_login(UserFactory())

        response = client.get(reverse("settings"))

        assert response.status_code == 200


class TestSuccess:
    def test_unauthenticated(self, client):
        """Only allow authenticated users."""
        response = client.get(reverse("success"))

        assert response.status_code == 302

    def test_authenticated(self, client):
        """An authenticated user gets a valid response."""
        client.force_login(UserFactory())

        response = client.get(reverse("success"))

        assert response.status_code == 200
