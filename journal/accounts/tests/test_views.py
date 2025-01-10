from unittest import mock

from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


class TestCheckEmail:
    def test_unauthenticated(self, client):
        """An unauthenticated user gets a valid response."""
        response = client.get(reverse("check-email"))

        assert response.status_code == 200


class TestCreateCheckoutSession:
    def test_unauthenticated(self, client):
        """Only allow authenticated users."""
        data = {"price_id": "1234"}

        response = client.post(
            reverse("create_checkout_session"),
            data=data,
            content_type="application/json",
        )

        assert response.status_code == 302

    def test_requires_post(self, client):
        """Only allow POST requests."""
        client.force_login(UserFactory())

        response = client.get(reverse("create_checkout_session"))

        assert response.status_code == 405

    @mock.patch("journal.accounts.views.PaymentsGateway")
    def test_authenticated(self, mock_payments_gateway_cls, client):
        """An authenticated user gets a valid response."""
        gateway = mock.MagicMock()
        gateway.create_checkout_session.return_value = "fake_session_id"
        mock_payments_gateway_cls.return_value = gateway
        user = UserFactory()
        client.force_login(user)
        data = {"price_id": "1234"}

        response = client.post(
            reverse("create_checkout_session"),
            data=data,
            content_type="application/json",
        )

        assert response.status_code == 200
        gateway.create_checkout_session.assert_called_with("1234", user)
        assert response.json()["session_id"] == "fake_session_id"


class TestCreateBillingPortalSession:
    def test_unauthenticated(self, client):
        """Only allow authenticated users."""

        response = client.post(reverse("create_billing_portal_session"))

        assert response.status_code == 302

    def test_requires_post(self, client):
        """Only allow POST requests."""
        client.force_login(UserFactory())

        response = client.get(reverse("create_billing_portal_session"))

        assert response.status_code == 405

    @mock.patch("journal.accounts.views.PaymentsGateway")
    def test_authenticated(self, mock_payments_gateway_cls, client):
        """An authenticated user gets a valid response."""
        gateway = mock.MagicMock()
        gateway.create_billing_portal_session.return_value = "portal_url"
        mock_payments_gateway_cls.return_value = gateway
        user = UserFactory()
        client.force_login(user)

        response = client.post(reverse("create_billing_portal_session"))

        assert response.status_code == 200
        gateway.create_billing_portal_session.assert_called_with(user)
        assert response.json()["url"] == "portal_url"


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
