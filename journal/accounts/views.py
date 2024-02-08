import json

from django.http import JsonResponse

from journal.payments.gateway import PaymentsGateway


def create_checkout_session(request):
    """Create a Stripe checkout session."""
    data = json.loads(request.body)
    price_id = data.get("price_id")

    # TODO: should be authenticated.
    # TODO: ensure POST
    gateway = PaymentsGateway()
    session_id = gateway.create_checkout_session(price_id, request.user)
    return JsonResponse({"session_id": session_id})
