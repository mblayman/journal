from django.http import JsonResponse


def create_checkout_session(request):
    """Create a Stripe checkout session."""
    # TODO: should be authenticated.
    # TODO: handle the request body (pull code from homeschool)
    return JsonResponse({"session_id": "fake_session_id"})
