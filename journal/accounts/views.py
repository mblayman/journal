import json

from django.contrib.auth.decorators import login_required
from django.http import JsonResponse
from django.shortcuts import render
from django.views.decorators.http import require_POST

from journal.payments.gateway import PaymentsGateway


@login_required
@require_POST
def create_checkout_session(request):
    """Create a Stripe checkout session."""
    data = json.loads(request.body)
    price_id = data.get("price_id")
    gateway = PaymentsGateway()
    session_id = gateway.create_checkout_session(price_id, request.user)
    return JsonResponse({"session_id": session_id})


@login_required
def success(request):
    """Show the success after account activation."""
    return render(request, "accounts/success.html", {})
