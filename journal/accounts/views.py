import json

from denied.authorizers import any_authorized
from denied.decorators import authorize
from django.http import JsonResponse
from django.shortcuts import render
from django.views.decorators.http import require_POST

from journal.payments.gateway import PaymentsGateway


@authorize(any_authorized)
@require_POST
def create_checkout_session(request):
    """Create a Stripe checkout session."""
    data = json.loads(request.body)
    price_id = data.get("price_id")
    gateway = PaymentsGateway()
    session_id = gateway.create_checkout_session(price_id, request.user)
    return JsonResponse({"session_id": session_id})


@authorize(any_authorized)
@require_POST
def create_billing_portal_session(request):
    """Create a billing portal session for a customer."""
    gateway = PaymentsGateway()
    portal_url = gateway.create_billing_portal_session(request.user)
    return JsonResponse({"url": portal_url})


@authorize(any_authorized)
def account_settings(request):
    """Show the user's settings."""
    return render(request, "accounts/settings.html", {})


@authorize(any_authorized)
def success(request):
    """Show the success after account activation."""
    return render(request, "accounts/success.html", {})
