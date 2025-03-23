from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


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


class TestUp:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("up"))

        assert response.status_code == 200
