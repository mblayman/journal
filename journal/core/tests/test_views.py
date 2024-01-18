from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


class TestIndex:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("index"))

        assert response.status_code == 200

    def test_authenticated(self, client):
        """An authenticatd user gets a valid response."""
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
