from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


class TestAbout:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("about"))

        assert response.status_code == 200


class TestFAQ:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("faq"))

        assert response.status_code == 200


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

    def test_signin(self, client):
        """A happy post redirect to the check email page."""
        data = {"email": "test@testing.com"}
        response = client.post(reverse("index"), data=data)

        assert response.status_code == 302


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


class TestUp:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("up"))

        assert response.status_code == 200
