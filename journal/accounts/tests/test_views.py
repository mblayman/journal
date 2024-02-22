from unittest import mock

from django.urls import reverse

from journal.accounts.tests.factories import UserFactory


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
